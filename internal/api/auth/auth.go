package auth

import (
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	templ "pass_web/internal/api/template"
	"path/filepath"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/google/uuid"
)

const challengeAmountOfChars = 20

type Page struct{
    Challenge string
    ChallengeID string
}

type UserChalenges struct{
    chalenges map[string]string   
}

var uc = UserChalenges{
    chalenges: make(map[string]string),
}

func generateChallenge(m int) string{
    var resp string
    asciiLetters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
    for i := 0; i < m; i++{
        idx := rand.Intn(52)
        resp += string(asciiLetters[idx])
    }
    return resp
}


func Handler(t *templ.Template) http.HandlerFunc{
    return func (w http.ResponseWriter, r *http.Request){
        sign := r.FormValue("signature")
        challengeId:= r.FormValue("challengeId")
        if sign != "" && challengeId != ""{
            original, ok := uc.chalenges[challengeId]
            // log.Println(uc.chalenges)
            if !ok{
                t.Render(w, "oob-auth-id-fail", struct{}{})
                return
            }

            pgp := crypto.PGP()

            file, err := os.Open("pubKeyTest.gpg")
            defer file.Close()         
            if err != nil{
                panic(err)
            }

            buffer, err := io.ReadAll(file)
            if err != nil{
                panic(err)
            }
            pubkeyArmored := string(buffer)
            


            publicKey, err := crypto.NewKeyFromArmored(pubkeyArmored)

            if err != nil{
                panic(err)
            }

            verifier, err := pgp.Verify().VerificationKey(publicKey).New()

            verifyResult, err := verifier.VerifyCleartext([]byte(sign))
            if err != nil{
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
            if decrypted == original{
                t.Render(w, "oob-auth-success", struct{}{})
                return
            }
        }
        dir := templ.GetTemplateDir()

        id := uuid.New()
        randomId := id.String() 

        randomChallenge :=  generateChallenge(challengeAmountOfChars)
        uc.chalenges[randomId] = randomChallenge
        t, err := template.ParseFiles(filepath.Join(dir, "base.tmpl"), filepath.Join(dir,"auth.tmpl"))
        if err != nil{
            panic(err)
        }

        t.Execute(w, Page{randomChallenge, randomId})
    }
}
