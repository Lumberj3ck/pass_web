package auth

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	templ "pass_web/internal/api/template"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

const challengeAmountOfChars = 20
const tokenExpirationDuration = time.Hour * 24
const challengeTimeToLive = time.Minute * 2

type Page struct {
	Challenge   string
	ChallengeID string
}

type Challenge struct {
	Challenge string
	Created   time.Time
}

type UserChalenges struct {
	chalenges map[string]Challenge
	mu        sync.RWMutex
}

func NewUserChalenges() *UserChalenges {
	uc := &UserChalenges{
		chalenges: make(map[string]Challenge),
	}

	go uc.cleanupExpired()
	return uc
}

func (uc *UserChalenges) cleanupExpired() {
	t := time.Tick(time.Minute * 1)

	for range t {
		uc.mu.Lock()
		for k, v := range uc.chalenges {
			if time.Since(v.Created) >= challengeTimeToLive {
				slog.Info("Deleted expired challenge", "challenge", v.Challenge)
				delete(uc.chalenges, k)
			}
		}
		uc.mu.Unlock()
	}
}

func (uc *UserChalenges) Add() (string, string) {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	id := uuid.New()
	randomId := id.String()

	randomChallenge := GenerateChallenge(challengeAmountOfChars)

	slog.Info("New challenge created", "challenge", randomChallenge)
	uc.chalenges[randomId] = Challenge{randomChallenge, time.Now()}
	return randomChallenge, randomId
}

func (uc *UserChalenges) Get(id string) (string, bool) {
	uc.mu.RLock()
	defer uc.mu.RUnlock()
	challenge, ok := uc.chalenges[id]
	return challenge.Challenge, ok
}

type JWTClaims struct {
	jwt.RegisteredClaims
}

var jwtSecret []byte

func init() {
	var err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	secretString := os.Getenv("jwt_secret")
	if secretString == "" {
		log.Fatal("jwt_secret environment variable not set")
	}
	jwtSecret = []byte(secretString)

}

func DecodeJWT(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse or validate token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims type")
	}

	return claims, nil
}

func GenerateJWT() (string, error) {
	claims := JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExpirationDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func AuthMiddlerware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth_cookie, err := r.Cookie("auth-token")
		t := templ.NewTemplate()

		if err != nil {
			unescaped_url := r.URL.String()
			next := url.QueryEscape(unescaped_url)
			http.Redirect(w, r, fmt.Sprintf("/auth?next=%v", next), http.StatusSeeOther)
			return
		}

		_, err = DecodeJWT(auth_cookie.Value)
		if err != nil {
			log.Println(err)

			cookie := http.Cookie{
				Name:     "auth-token",
				Value:    "",
				HttpOnly: true,
				Path:     "/",
				Expires:  time.Now().Add(tokenExpirationDuration),
				MaxAge:   -1,
				SameSite: http.SameSiteLaxMode,
			}

			http.SetCookie(w, &cookie)
			w.WriteHeader(http.StatusNetworkAuthenticationRequired)
			t.Render(w, "auth_provided_cookie_invalid", struct{}{})
			return
		}

		next(w, r)
	}
}

func GenerateChallenge(m int) string {
	var resp string
	asciiLetters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	for i := 0; i < m; i++ {
		idx := rand.Intn(52)
		resp += string(asciiLetters[idx])
	}
	return resp
}

type Notification struct {
	Title string
	Message string
	DisplayDuration int
}

func Handler(uc *UserChalenges) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			challenge, randomId := uc.Add()
			t := templ.NewTemplate("templates/base.tmpl", "templates/auth.tmpl", "templates/notifications.tmpl")
			t.Render(w, "", Page{challenge, randomId})
		} else {
			cookie, _ := r.Cookie("auth-token")

			next := r.URL.Query().Get("next")

			if cookie != nil && next == "" {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			t := templ.NewTemplate()
			sign := r.FormValue("signature")
			challengeId := r.FormValue("challengeId")
			if sign != "" && challengeId != "" {
				original, ok := uc.Get(challengeId)
				if !ok {
					t.Render(w, "oob-warning-notification", Notification{"Warning", "Challenge id either expired or incorect", 200000})
					return
				}

				pgp := crypto.PGP()

				pub_key_path := os.Getenv("PASSWEB_PUB_KEY_PATH")
				if len(pub_key_path) == 0 {
					pub_key_path = "pubKeyTest.gpg"
				}
				log.Println(pub_key_path)
				file, err := os.Open(pub_key_path)
				defer file.Close()
				if err != nil {
					t.Render(w, "oob-warning-notification", Notification{"Warning", "Public key path not found", 2000})
					return
				}

				buffer, err := io.ReadAll(file)
				if err != nil {
					t.Render(w, "oob-warning-notification", Notification{"Warning", "Public key could not be read", 2000})
					return
				}
				pubkeyArmored := string(buffer)

				publicKey, err := crypto.NewKeyFromArmored(pubkeyArmored)

				if err != nil {
					t.Render(w, "oob-warning-notification", Notification{"Warning", " Invalid public key", 2000})
					return
				}

				verifier, err := pgp.Verify().VerificationKey(publicKey).New()

				verifyResult, err := verifier.VerifyCleartext([]byte(sign))
				if err != nil {
					log.Println(err)
					t.Render(w, "oob-warning-notification", Notification{"Warning", "Provided signature is not valid", 2000})
					return
				}

				if sigErr := verifyResult.SignatureError(); sigErr != nil {
					log.Println("Check sign sig err")
					log.Println(sigErr.Error())
					t.Render(w, "oob-warning-notification", Notification{"Warning", "Provided signature is not valid", 2000})
					return
				}

				log.Println(string(verifyResult.Cleartext()))

				decrypted := string(verifyResult.Cleartext())
				if decrypted == original {

					jwt_auth_token, err := GenerateJWT()
					if err != nil {
						slog.Error("Failed to generate JWT token", "error", err)
						t.Render(w, "oob-warning-notification", Notification{"Warning", "Failed to generate JWT token", 2000})
						return
					}

					cookie := http.Cookie{
						Name:     "auth-token",
						Value:    jwt_auth_token,
						HttpOnly: true,
						Path:     "/",
						Expires:  time.Now().Add(tokenExpirationDuration),
						SameSite: http.SameSiteLaxMode,
					}

					http.SetCookie(w, &cookie)
					w.Header().Set("HX-Redirect", next)
					w.WriteHeader(http.StatusOK)
					return
				}
			}
		}

	}
}
