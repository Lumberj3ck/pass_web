package auth

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"net/http/httptest"
	"os"
	templ "pass_web/internal/api/template"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/net/html"

	"github.com/stretchr/testify/assert"
)

func filterChallenges(s string) string{
    var buff bytes.Buffer
    scanner := bufio.NewScanner(strings.NewReader(s))

    for scanner.Scan(){
        line := scanner.Text()

        if strings.Contains(line, "name=\"challenge"){
            continue
        }

        buff.WriteString(line + "\n")
    }
    return buff.String()
}
func getAuthPage() httptest.ResponseRecorder{
    r, err := http.NewRequest(http.MethodGet,  "auth", nil)
    if err != nil{
        panic(err)
    }

    recorder := httptest.NewRecorder()


    te := templ.NewTemplate()

    Handler(&te)(recorder, r)
    return *recorder
}


func traverse(n *html.Node){
    if n.Type == html.ElementNode{
        for _, attr := range n.Attr{
            if attr.Key == "id" && attr.Val == "challenge"{
            }
            fmt.Println(attr)
        }
    }
    for c := n.FirstChild; c != nil; c = c.NextSibling{
        traverse(c)
    }
}

func getChallengeAndChallengeId() (string, string){
    recorder := getAuthPage()

    doc, err := html.Parse(recorder.Result().Body)
    if err != nil{
        panic(err)
    }

    var challenge, challengeId string
    var findElement func(*html.Node)
    findElement = func(n *html.Node){
        if n.Type == html.ElementNode{
            elFound := ""
            for _, attr := range n.Attr{
                if attr.Key == "id"{
                    if attr.Val == "challenge" || attr.Val == "challengeId"{
                        elFound = attr.Val
                    }                 
                }
            }
            switch elFound{
            case "challenge":
                for _, attr := range n.Attr{
                    if attr.Key == "value"{
                        challenge = attr.Val
                    }
                }
            case "challengeId":
                for _, attr := range n.Attr{
                    if attr.Key == "value"{
                        challengeId = attr.Val
                    }
                }
            }
        }
        if challenge == "" || challengeId == ""{
            for c := n.FirstChild; c != nil; c = c.NextSibling{
                findElement(c)
            }
        }
    }

    findElement(doc)
    return challenge, challengeId
}

func TestAuthHappyPathWithoutChallenge(t *testing.T){
    recorder := getAuthPage()
    assert.Equal(t, "text/html; charset=utf-8", recorder.Header().Get("Content-Type"))
                                 
    dir := filepath.Join(os.Getenv(templ.PROJECT_ROOT_ENV), "templates")
    test_templ, err := template.ParseFiles(filepath.Join(dir, "base.tmpl"), filepath.Join(dir,"auth.tmpl"))

    if err != nil{
        panic(err)
    }


    var buf bytes.Buffer

    test_templ.Execute(&buf, Page{})
    
    want := filterChallenges(buf.String())
    got := filterChallenges(recorder.Body.String())
    assert.Equal(t, want, got)
} 


func TestAuthChallengeValid(t *testing.T){
    challenge, challengeId := getChallengeAndChallengeId()
    assert.NotEmpty(t, challenge)
    assert.NotEmpty(t, challengeId)
    assert.Regexp(t, `\w+`, challenge)
    assert.Regexp(t, "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$", challengeId)
    
    t.Log(challenge, challengeId)
}

func TestSignFail(t *testing.T){
    challenge, challengeId := getChallengeAndChallengeId()
    formData := url.Values{}
    formData.Set("signature", challenge)
    formData.Set("challengeId", challengeId)
    payload := formData.Encode()

    req, err := http.NewRequest(http.MethodPost, "auth", bytes.NewBufferString(payload))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    if err != nil{
        panic(err)
    }

    recorder := httptest.NewRecorder()
    
    te := templ.NewTemplate()
    Handler(&te)(recorder, req)
    t.Log(recorder.Body.String())
}

