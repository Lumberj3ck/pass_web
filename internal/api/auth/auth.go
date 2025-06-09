package auth

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	templ "pass_web/internal/api/template"
	"path/filepath"
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
			w.WriteHeader(http.StatusNetworkAuthenticationRequired)
			w.Write([]byte("Please provide auth cookie"))
			return
		}

		_, err = DecodeJWT(auth_cookie.Value)
		if err != nil {
			log.Println(err)
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

func Handler(t *templ.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sign := r.FormValue("signature")
		challengeId := r.FormValue("challengeId")
		if sign != "" && challengeId != "" {
			original, ok := uc.chalenges[challengeId]
			// log.Println(uc.chalenges)
			if !ok {
				t.Render(w, "oob-auth-id-fail", struct{}{})
				return
			}

			pgp := crypto.PGP()

			file, err := os.Open("pubKeyTest.gpg")
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
				w.Header().Set("HX-Redirect", "/show")
				w.WriteHeader(http.StatusOK)
				return
			}
		}
		dir := templ.GetTemplateDir()

		id := uuid.New()
		randomId := id.String()

		randomChallenge := GenerateChallenge(challengeAmountOfChars)
		uc.chalenges[randomId] = randomChallenge
		t, err := template.ParseFiles(filepath.Join(dir, "base.tmpl"), filepath.Join(dir, "auth.tmpl"))
		if err != nil {
			panic(err)
		}

		t.Execute(w, Page{randomChallenge, randomId})
	}
}
