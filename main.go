package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	OutDir         string
	IntegrationId  string
	InstallationId int64
	PrivateKeyFile string
)

type GHResponse struct {
	Token     string    `json:"token,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

func init() {
	flag.StringVar(&OutDir, "git-config-out", "", "path where we will write the .gitconfig file")
	flag.StringVar(&IntegrationId, "gh-integration-id", "1234", "Github Integration's id")
	flag.StringVar(&PrivateKeyFile, "gh-private-key", "private_key.pem", "Full path to the Github Integration's private key")
	flag.Int64Var(&InstallationId, "gh-installation-id", 5678, "Github Integtation's installation id")
}

func fatalErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	flag.Parse()
	data, err := ioutil.ReadFile(PrivateKeyFile)
	fatalErr(err)
	key, err := jwt.ParseRSAPrivateKeyFromPEM(data)
	fatalErr(err)
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
		Issuer:    IntegrationId,
	}).SignedString(key)
	fatalErr(err)

	var result GHResponse
	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.github.com/installations/%d/access_tokens", InstallationId), nil)
	fatalErr(err)

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Accept", "application/vnd.github.machine-man-preview+json")
	fatalErr(JsonResponse(req, &result))

	if OutDir == "" {
		fmt.Print(result.Token)
	} else {
		path, err := filepath.Abs(OutDir)
		fatalErr(err)
		f, err := os.Create(path)
		fatalErr(err)
		defer f.Close()
		t := template.Must(template.New("t1").
			Parse(`[url "https://x-access-token:{{.}}@github.com/"]
	insteadOf = https://github.com/`))
		fatalErr(t.Execute(f, result.Token))
	}
}

func JsonResponse(req *http.Request, target interface{}) error {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}
		log.Fatalln(string(data))
	}
	return json.NewDecoder(resp.Body).Decode(target)
}
