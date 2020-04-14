package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"user-service/model"
)

const (
	hydraEndpoint = "http://127.0.0.1:4445"
)

func main() {

	http.HandleFunc("/login", loginHanlder)
	http.HandleFunc("/consent", consentHandler)
	http.HandleFunc("/acceptConsent", acceptConsentHandler)
	http.HandleFunc("/logout", logoutHandler)
	log.Println("Server is running at 3000 port.")
	http.ListenAndServe(":3000", nil)

}

func loginHanlder(w http.ResponseWriter, r *http.Request) {

	urlChallengeParams := r.URL.Query()["login_challenge"]
	if len(urlChallengeParams) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Challange Parameters not existing or to many parameters!"))
		return
	}
	challenge := urlChallengeParams[0]

	if r.Method == "POST" {
		if r.Form == nil {
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		userName := r.Form.Get("username")
		password := r.Form.Get("password")
		loginChallenge := r.Form.Get("challenge")

		if userName == "toky" && password == "pwd" {

			// könnte noch gekürzt werden mit dem obigen acceptLogin body
			acceptLoginBody := model.AcceptLogin{
				Subject:     userName,
				Remember:    true,
				RememberFor: 20,
			}
			rawJson, err := json.Marshal(acceptLoginBody)

			redirectURL, err := sendAcceptBody("login", loginChallenge, rawJson)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Fatal(err)
			}
			http.Redirect(w, r, redirectURL, http.StatusFound)
		}
		w.WriteHeader(http.StatusForbidden)
		templLogin := template.Must(template.ParseFiles("templates/login.html"))
		loginData := model.LoginPageData{
			PageTitle:        "Test Login",
			LoginButtonLabel: "Einloggen",
			UserNameLabel:    "Benutzername",
			PasswordLabel:    "Passwort",
			LoginLabel:       "Login",
			Challenge:        challenge,
			ErrorMessage:     "Benutzername oder Passwort falsch",
		}
		templLogin.Execute(w, loginData)
	} else {

		challengeBody, err := readChallenge(challenge, "login")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Print(err)
		}

		if !challengeBody.Skip {
			templLogin := template.Must(template.ParseFiles("templates/login.html"))
			loginData := model.LoginPageData{
				PageTitle:        "Test Login",
				LoginButtonLabel: "Einloggen",
				UserNameLabel:    "Benutzername",
				PasswordLabel:    "Passwort",
				LoginLabel:       "Login",
				Challenge:        challenge,
			}
			templLogin.Execute(w, loginData)
		} else {

			acceptLoginBody := model.AcceptLogin{
				Subject:     challengeBody.Subject,
				Remember:    true,
				RememberFor: 3600,
			}
			rawJson, err := json.Marshal(acceptLoginBody)

			redirectURL, err := sendAcceptBody("login", challenge, rawJson)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Fatal(err)
			}
			http.Redirect(w, r, redirectURL, http.StatusFound)

		}
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	urlChallengeParams := r.URL.Query()["logout_challenge"]
	if len(urlChallengeParams) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Challange Parameters not existing or to many parameters!"))
		return
	}
	challenge := urlChallengeParams[0]

	if r.Method == "POST" {
		var err error
		if r.Form == nil {
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		accept := r.Form.Get("accept")
		logoutChallenge := r.Form.Get("challenge")
		var redirectURL string

		if accept == "true" {
			redirectURL, err = sendAcceptBody("logout", logoutChallenge, nil)

		} else {
			redirectURL, err = sendRejectBody("logout", logoutChallenge, nil)
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
		}

		http.Redirect(w, r, redirectURL, http.StatusFound)
	} else {

		challengeBody, err := readChallenge(challenge, "logout")
		if err != nil {
			log.Println(err)
		}

		if challengeBody.RpInitiated {
			templLogout := template.Must(template.ParseFiles("templates/logout.html"))
			logoutData := model.LogoutPage{
				Subject:           challengeBody.Subject,
				LogoutButtonLabel: "Ausloggen",
				LogoutDenyLabel:   "Nicht ausloggen",
				Challenge:         challenge,
				LogoutTitle:       "Logout Anforderung",
				PageTitle:         "Logout",
			}
			templLogout.Execute(w, logoutData)
		} else {
			redirectURL, err := sendAcceptBody("logout", challenge, nil)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Fatal(err)
			}

			http.Redirect(w, r, redirectURL, http.StatusFound)
		}
	}

}

func consentHandler(w http.ResponseWriter, r *http.Request) {

	urlChallengeParams := r.URL.Query()["consent_challenge"]
	if len(urlChallengeParams) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Challange Parameters not existing or to many parameters!"))
		return
	}
	challenge := urlChallengeParams[0]

	challengeBody, err := readChallenge(challenge, "consent")

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)

	}

	requestedScopes := make([]model.ReqestScope, 0, len(challengeBody.RequestedScope))
	for _, scope := range challengeBody.RequestedScope {
		requestedScopes = append(requestedScopes, model.ReqestScope{ScopeValue: scope, ScopeName: scope})
	}

	grantedAccesToken := make([]model.ReqestScope, 0, len(challengeBody.RequestedAccessToken))
	for _, accessToken := range challengeBody.RequestedAccessToken {
		grantedAccesToken = append(grantedAccesToken, model.ReqestScope{ScopeName: accessToken, ScopeValue: "true"}) // hier könnten die tokens gefiltert werden
	}
	if !challengeBody.Skip {
		consentData := model.ConsentData{
			PageTitle:            "Auth",
			RequestMessage:       fmt.Sprintf("Die Seite %s fordert die Folgenden Berechtigungen an", challengeBody.Client.ClientID),
			AuthorizeButtonLabel: "Authorisieren",
			AuthorizeTitle:       "Authorisierung",
			GrantedAccessLabel:   "Die Folgenden Authorisierungen sind für diesen Client Erlaubt",
			ReqestScopes:         requestedScopes,
			Challenge:            challenge,
			GrantedAccessToken:   grantedAccesToken,
		}

		templConsent := template.Must(template.ParseFiles("templates/consent.html"))

		templConsent.Execute(w, consentData)
	} else {

		redirectURL, err := redirectFromConsent(challengeBody.RequestedScope, challengeBody.RequestedAccessToken, challenge)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
		}
		http.Redirect(w, r, redirectURL, http.StatusFound)

	}

}

func acceptConsentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if r.Form == nil {
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		r.ParseForm()
		allowedScopes := r.PostForm["scope"]
		allowedAccessToken := r.PostForm["accesToken"]
		consentChallenge := r.Form.Get("challenge")
		redirectURL, err := redirectFromConsent(allowedScopes, allowedAccessToken, consentChallenge)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
		}
		http.Redirect(w, r, redirectURL, http.StatusFound)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func redirectFromConsent(allowedScopes, allowedAccessToken []string, consentChallenge string) (redirectUrl string, err error) {
	scope := make([]string, 0, len(allowedScopes))
	for _, allowedScope := range allowedScopes {
		scope = append(scope, string(allowedScope))
	}
	accesToken := make([]string, 0, len(allowedAccessToken))
	for _, allowedToken := range allowedAccessToken {
		accesToken = append(accesToken, string(allowedToken))
	}

	acceptConsentBody := model.AcceptConsent{
		GrantScope:               scope,
		GrantAccessTokenAudience: accesToken, // TODO sollte noch nachgeführt werden
		Remember:                 true,
		RememberFor:              0,
		Session: model.SessionInfo{
			//TODO Beispieldaten
			AccessToken: map[string]string{
				"email":    "marco.jakob3@gmail.com",
				"userName": "Toky",
			},
			IDToken: map[string]string{
				"username": "Marco",
				"lastname": "Jakob",
				"email":    "marco.jakob3@gmail.com",
			},
		},
	}

	rawJson, err := json.Marshal(acceptConsentBody)

	return sendAcceptBody("consent", consentChallenge, rawJson)
}

func readChallenge(loginChallenge, challengeMethod string) (challengeBody model.LoginChallenge, err error) {
	headers := map[string][]string{
		"Accept": []string{"application/json"},
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/oauth2/auth/requests/%s?%s_challenge=%s", hydraEndpoint, challengeMethod, challengeMethod, loginChallenge), nil)
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

func sendRejectBody(method, challenge string, rawJson []byte) (redirectUrl string, err error) {
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/oauth2/auth/requests/%s/reject?%s_challenge=%s", hydraEndpoint, method, method, challenge), bytes.NewBuffer(rawJson))
	if err != nil {
		log.Println(err)
	}
	return sendRequest(req)
}

func sendAcceptBody(method, challenge string, rawJson []byte) (redirectUrl string, err error) {
	headers := map[string][]string{
		"Accept":       []string{"application/json"},
		"Content-Type": []string{"application/json"},
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/oauth2/auth/requests/%s/accept?%s_challenge=%s", hydraEndpoint, method, method, challenge), bytes.NewBuffer(rawJson))
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
