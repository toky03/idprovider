package adapter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"user-service/model"
)

type HydraAdapter struct {
	hydraEndpoint string
}

// NewHydraAdapter creates new Instance of Adapter Service
func NewHydraAdapter() HydraAdapter {
	hydraEndpoint := os.Getenv("HYDRA_BASE_URL")
	if hydraEndpoint == "" {
		hydraEndpoint = "http://127.0.0.1:4445"
	}
	return HydraAdapter{
		hydraEndpoint: hydraEndpoint,
	}

}

// ReadChallenge fetch data from Challgnge
func (a *HydraAdapter) ReadChallenge(loginChallenge, challengeMethod string) (challengeBody model.LoginChallenge, err error) {
	headers := map[string][]string{
		"Accept": []string{"application/json"},
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/oauth2/auth/requests/%s?%s_challenge=%s", a.hydraEndpoint, challengeMethod, challengeMethod, loginChallenge), nil)
	if err != nil {
		log.Print(err)
		return
	}
	req.Header = headers

	client := &http.Client{}
	res, err := client.Do(req)

	body, err := ioutil.ReadAll(res.Body)

	var challangeBody model.LoginChallenge

	err = json.Unmarshal(body, &challangeBody)

	if err != nil {
		log.Print(err)
		return
	}

	res.Body.Close()
	return challangeBody, err
}

// SendRejectBody used to reqject requests for login, logout or consent
func (a *HydraAdapter) SendRejectBody(method, challenge string, rawJson []byte) (redirectUrl string, err error) {
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/oauth2/auth/requests/%s/reject?%s_challenge=%s", a.hydraEndpoint, method, method, challenge), bytes.NewBuffer(rawJson))
	if err != nil {
		log.Println(err)
	}
	return sendRequest(req)
}

// SendAcceptBody used to accept requests
func (a *HydraAdapter) SendAcceptBody(method, challenge string, rawJson []byte) (redirectUrl string, err error) {
	headers := map[string][]string{
		"Accept":       []string{"application/json"},
		"Content-Type": []string{"application/json"},
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/oauth2/auth/requests/%s/accept?%s_challenge=%s", a.hydraEndpoint, method, method, challenge), bytes.NewBuffer(rawJson))
	if err != nil {
		log.Println(err)
	}
	req.Header = headers
	return sendRequest(req)
}

func sendRequest(req *http.Request) (redirectUrl string, err error) {

	client := &http.Client{}

	res, err := client.Do(req)

	var redirect model.Redirect

	body, err := ioutil.ReadAll(res.Body)

	err = json.Unmarshal(body, &redirect)

	if err != nil {
		log.Print(err)
		return
	}

	redirectUrl = redirect.RedirectURL

	res.Body.Close()
	return

}
