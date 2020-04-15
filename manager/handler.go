package manager

import (
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net/http"
	"user-service/adapter"
	"user-service/model"
)

type Handler struct {
	Service Service
	Adapter adapter.Adapter
}

func NewServieHandler() Handler {
	service, err := NewService()
	if err != nil {
		log.Fatal(err)
	}
	adapter := adapter.NewAdapter()

	return Handler{
		Service: service,
		Adapter: adapter,
	}

}

type ServiceHandler interface {
	LoginHandler(w http.ResponseWriter, r *http.Request)
	LogoutHandler(w http.ResponseWriter, r *http.Request)
	ConsentHandler(w http.ResponseWriter, r *http.Request)
	AceptconsentHandler(w http.ResponseWriter, r *http.Request)
}

// LoginHandler handles login requests
func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {

	challenge, err := readURLChallangeParams(r, "login")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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

			acceptLoginBody := h.Service.FetchAcceptLoginConfig(userName)
			rawJson, err := json.Marshal(acceptLoginBody)

			redirectURL, err := h.Adapter.SendAcceptBody("login", loginChallenge, rawJson)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Fatal(err)
			}
			http.Redirect(w, r, redirectURL, http.StatusFound)
		}

		w.WriteHeader(http.StatusForbidden)
		templLogin := template.Must(template.ParseFiles("templates/login.html"))
		loginData := h.Service.FetchLoginConfig(challenge, true)
		templLogin.Execute(w, loginData)
	} else {
		challengeBody, err := h.Adapter.ReadChallenge(challenge, "login")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Print(err)
		}

		if !challengeBody.Skip {
			templLogin := template.Must(template.ParseFiles("templates/login.html"))
			loginData := h.Service.FetchLoginConfig(challenge, false)
			templLogin.Execute(w, loginData)
		} else {

			acceptLoginBody := h.Service.FetchAcceptLoginConfig(challengeBody.Subject)
			rawJson, err := json.Marshal(acceptLoginBody)

			redirectURL, err := h.Adapter.SendAcceptBody("login", challenge, rawJson)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Fatal(err)
			}
			http.Redirect(w, r, redirectURL, http.StatusFound)

		}
	}
}

// LogoutHandler handles logout requests
func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {

	challenge, err := readURLChallangeParams(r, "logout")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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
			redirectURL, err = h.Adapter.SendAcceptBody("logout", logoutChallenge, nil)

		} else {
			redirectURL, err = h.Adapter.SendRejectBody("logout", logoutChallenge, nil)
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
		}

		http.Redirect(w, r, redirectURL, http.StatusFound)
	} else {

		challengeBody, err := h.Adapter.ReadChallenge(challenge, "logout")
		if err != nil {
			log.Println(err)
		}

		if challengeBody.RpInitiated {
			templLogout := template.Must(template.ParseFiles("templates/logout.html"))
			logoutData := h.Service.FetchLogoutConfig(challenge, challengeBody.Subject)
			templLogout.Execute(w, logoutData)
		} else {
			redirectURL, err := h.Adapter.SendAcceptBody("logout", challenge, nil)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Fatal(err)
			}

			http.Redirect(w, r, redirectURL, http.StatusFound)
		}
	}

}

func (h *Handler) ConsentHandler(w http.ResponseWriter, r *http.Request) {

	challenge, err := readURLChallangeParams(r, "consent")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	challengeBody, err := h.Adapter.ReadChallenge(challenge, "consent")

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
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
		consentData := h.Service.FetchConsentConfig(challengeBody.Client.ClientID, challenge, requestedScopes, grantedAccesToken)
		templConsent := template.Must(template.ParseFiles("templates/consent.html"))

		templConsent.Execute(w, consentData)
	} else {

		redirectURL, err := h.Adapter.RedirectFromConsent(challengeBody.RequestedScope, challengeBody.RequestedAccessToken, challenge)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
		}
		http.Redirect(w, r, redirectURL, http.StatusFound)

	}

}

func (h *Handler) AcceptConsentHandler(w http.ResponseWriter, r *http.Request) {
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
		redirectURL, err := h.Adapter.RedirectFromConsent(allowedScopes, allowedAccessToken, consentChallenge)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
		}
		http.Redirect(w, r, redirectURL, http.StatusFound)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func readURLChallangeParams(r *http.Request, challengeMethod string) (string, error) {
	urlChallengeParams := r.URL.Query()[challengeMethod+"_challenge"]
	if len(urlChallengeParams) != 1 {
		return "", errors.New("Challange Parameters not existing or to many parameters!")
	}
	return urlChallengeParams[0], nil
}
