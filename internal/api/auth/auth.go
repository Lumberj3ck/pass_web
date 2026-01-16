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
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

const challengeAmountOfChars = 20
const tokenExpirationDuration = time.Hour * 24

type Page struct {
	Challenge   string
	ChallengeID string
}

type UserChalenges struct {
	chalenges map[string]string
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

var uc = UserChalenges{
	chalenges: make(map[string]string),
}

func AuthMiddlerware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth_cookie, err := r.Cookie("auth-token")
		t := templ.NewTemplate()

		if err != nil {
			unescaped_url := r.URL.String()
			slog.Info("Not authentication ", "unescaped_url", unescaped_url)
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

func Handler(w http.ResponseWriter, r *http.Request) {
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
		original, ok := uc.chalenges[challengeId]
		if !ok {
			t.Render(w, "oob-auth-id-fail", struct{}{})
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
			panic(err)
		}

		buffer, err := io.ReadAll(file)
		if err != nil {
			panic(err)
		}
		pubkeyArmored := string(buffer)

		publicKey, err := crypto.NewKeyFromArmored(pubkeyArmored)

		if err != nil {
			panic(err)
		}

		verifier, err := pgp.Verify().VerificationKey(publicKey).New()

		verifyResult, err := verifier.VerifyCleartext([]byte(sign))
		if err != nil {
			log.Println(err)
			t.Render(w, "oob-auth-signature-fail", struct{}{})
			return
		}

		if sigErr := verifyResult.SignatureError(); sigErr != nil {
			log.Println("Check sign sig err")
			log.Println(sigErr.Error())
			t.Render(w, "oob-auth-signature-fail", struct{}{})

			return
		}

		log.Println(string(verifyResult.Cleartext()))

		decrypted := string(verifyResult.Cleartext())
		if decrypted == original {

			jwt_auth_token, err := GenerateJWT()
			if err != nil {
				panic(err)
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

	id := uuid.New()
	randomId := id.String()

	randomChallenge := GenerateChallenge(challengeAmountOfChars)
	uc.chalenges[randomId] = randomChallenge
	// t, err := template.ParseFiles(filepath.Join(dir, "base.tmpl"), filepath.Join(dir, "auth.tmpl"))
	t = templ.NewTemplate("templates/base.tmpl", "templates/auth.tmpl")
	// if err != nil {
	// 	panic(err)
	// }

	t.Render(w, "", Page{randomChallenge, randomId})
}
