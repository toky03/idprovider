package model

type LoginChallenge struct {
	Skip                 bool                   `json:"skip"`
	Subject              string                 `json:"subject"`
	Client               Client                 `json:"client"`
	RequestURL           string                 `json:"request_url"`
	RequestedScope       []string               `json:"requested_scope"`
	RequestedAccessToken []string               `json:"requested_access_token_audience"`
	RpInitiated          bool                   `json:"rp_initiated"`
	X                    map[string]interface{} `json:"-"`
}

type Client struct {
	ClientID   string `json:"client_id"`
	ClientName string `json:"client_name"`
}

type AcceptLogin struct {
	Subject     string `json:"subject"`
	Remember    bool   `json:"remember"`
	RememberFor int16  `json:"remember_for"`
}

type AcceptConsent struct {
	GrantScope               []string    `json:"grant_scope"`
	GrantAccessTokenAudience []string    `json:"grant_access_token_audience"`
	Remember                 bool        `json:"remember"`
	RememberFor              int16       `json:"remember_for"`
	Session                  SessionInfo `json:"session"`
}

type SessionInfo struct {
	AccessToken map[string]string `json:"access_token"`
	IDToken     map[string]string `json:"id_token"`
}

type Redirect struct {
	RedirectURL string `json:"redirect_to"`
}
