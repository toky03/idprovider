package manager

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"user-service/model"
)

// Service Handler for Config and Services
type ConfigService struct {
	LoginData       model.LoginPageData
	LogoutData      model.LogoutPage
	ConsentData     model.ConsentData
	AcceptLoginData model.AcceptLogin
}

// NewService creates new instance of a Service
func NewConfigService() (manager ConfigService, err error) {
	pwd, err := os.Getwd()
	if pwd == "/" {
		pwd = ""
	}
	var loginPageData model.LoginPageData
	loginFile, err := os.Open(pwd + "/config/login_config.json")
	if err != nil {
		log.Println(err)
	}
	decoder := json.NewDecoder(loginFile)
	err = decoder.Decode(&loginPageData)
	if err != nil {
		log.Println(err)
	}
	var logoutPageData model.LogoutPage
	logoutfile, err := os.Open(pwd + "/config/logout_config.json")
	if err != nil {
		log.Println(err)
	}
	decoder = json.NewDecoder(logoutfile)
	err = decoder.Decode(&logoutPageData)
	if err != nil {
		log.Println(err)
	}

	var consentPageData model.ConsentData
	consentFile, err := os.Open(pwd + "/config/consent_config.json")
	if err != nil {
		log.Println(err)
	}
	decoder = json.NewDecoder(consentFile)
	err = decoder.Decode(&consentPageData)
	if err != nil {
		log.Println(err)
	}

	var acceptLoginData model.AcceptLogin
	acceptLoginFile, err := os.Open(pwd + "/config/accept_login_config.json")
	if err != nil {
		log.Println(err)
	}
	decoder = json.NewDecoder(acceptLoginFile)
	err = decoder.Decode(&acceptLoginData)
	if err != nil {
		log.Println(err)
	}

	manager = ConfigService{
		LoginData:       loginPageData,
		LogoutData:      logoutPageData,
		ConsentData:     consentPageData,
		AcceptLoginData: acceptLoginData,
	}
	return
}

// FetchLoginConfig returns prepared Login Page Data
func (s *ConfigService) FetchLoginConfig(challenge string, withError bool) (loginPageData model.LoginPageData) {

	loginPageData = s.LoginData
	loginPageData.Challenge = challenge
	if withError {
		loginPageData.ErrorMessage = "Benutzername oder Passwort falsch"
	}

	return
}

// FetchLogoutConfig returns prepared Logout Page Data
func (s *ConfigService) FetchLogoutConfig(challenge, subject string) (logoutPageData model.LogoutPage) {
	logoutPageData = s.LogoutData
	logoutPageData.Challenge = challenge
	logoutPageData.Subject = subject
	return
}

// FetchConsentConfig returns prepared Consent Page Data
func (s *ConfigService) FetchConsentConfig(clientID, challenge, userName string, requestedScopes, grantedAccesToken []model.ReqestScope) (consentPageData model.ConsentData) {
	consentPageData = s.ConsentData
	consentPageData.UserName = userName
	consentPageData.ClientName = clientID
	consentPageData.RequestMessage = fmt.Sprintf(s.ConsentData.RequestMessage, clientID)
	consentPageData.Challenge = challenge
	consentPageData.ReqestScopes = requestedScopes
	consentPageData.GrantedAccessToken = grantedAccesToken
	return

}

// FetchAcceptLoginConfig returns prepared Accept Login Data
func (s *ConfigService) FetchAcceptLoginConfig(userName string) (acceptLoginData model.AcceptLogin) {
	acceptLoginData = s.AcceptLoginData
	acceptLoginData.Subject = userName
	return

}
