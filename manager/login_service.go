package manager

import (
	"encoding/json"
	"user-service/adapter"
	"user-service/model"
)

type LoginAdapter interface {
	ReadChallenge(string, string) (model.LoginChallenge, error)
	SendRejectBody(string, string, []byte) (string, error)
	SendAcceptBody(string, string, []byte) (string, error)
}

type LoginService struct {
	UserService  UserService
	HydraAdapter LoginAdapter
}

func NewLoginService() LoginService {

	loginAdapter := adapter.NewHydraAdapter()

	return LoginService{
		UserService:  NewUserService(),
		HydraAdapter: &loginAdapter,
	}
}

func (s *LoginService) CheckPasswords(userName, password string) (bool, error) {
	return s.UserService.CheckPassword(userName, password)
}

// ReadChallenge fetch data from Challgnge
func (s *LoginService) ReadChallenge(loginChallenge, challengeMethod string) (challengeBody model.LoginChallenge, err error) {
	return s.HydraAdapter.ReadChallenge(loginChallenge, challengeMethod)
}

// SendRejectBody used to reqject requests for login, logout or consent
func (s *LoginService) SendRejectBody(method, challenge string, rawJson []byte) (redirectUrl string, err error) {
	return s.HydraAdapter.SendRejectBody(method, challenge, rawJson)
}

// SendAcceptBody used to accept requests
func (s *LoginService) SendAcceptBody(method, challenge string, rawJson []byte) (redirectUrl string, err error) {
	return s.HydraAdapter.SendAcceptBody(method, challenge, rawJson)
}

func (s *LoginService) RedirectFromConsent(allowedScopes, allowedAccessToken []string, consentChallenge, userName, clientName string) (redirectUrl string, err error) {
	scope := make([]string, 0, len(allowedScopes))
	for _, allowedScope := range allowedScopes {
		scope = append(scope, string(allowedScope))
	}
	accesToken := make([]string, 0, len(allowedAccessToken))
	for _, allowedToken := range allowedAccessToken {
		accesToken = append(accesToken, string(allowedToken))
	}

	user, err := s.UserService.FindUserByEmailOrUserName(userName)
	if err != nil {
		return
	}
	roles := make([]string, 0, len(user.Applications))

	for _, application := range user.Applications {
		if application.ApplicationName == clientName {
			roles = append(roles, application.Roles...)
		}
	}

	userInfoToken := model.UserInfoToken{
		EMail:    user.Email,
		LastName: user.LastName,
		UserName: user.UserName,
		Roles:    roles,
	}

	acceptConsentBody := model.AcceptConsent{
		GrantScope:               scope,
		GrantAccessTokenAudience: accesToken,
		Remember:                 true,
		RememberFor:              0,
		Session: model.SessionInfo{
			AccessToken: userInfoToken,
			IDToken:     userInfoToken,
		},
	}

	rawJson, err := json.Marshal(acceptConsentBody)

	return s.HydraAdapter.SendAcceptBody("consent", consentChallenge, rawJson)
}
