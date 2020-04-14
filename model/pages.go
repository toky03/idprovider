package model

type LoginPageData struct {
	PageTitle        string
	LoginLabel       string
	UserNameLabel    string
	PasswordLabel    string
	LoginButtonLabel string
	Challenge        string
	ErrorMessage     string
}

type ConsentData struct {
	PageTitle            string
	AuthorizeTitle       string
	RequestMessage       string
	ReqestScopes         []ReqestScope
	AuthorizeButtonLabel string
	Challenge            string
	GrantedAccessLabel   string
	GrantedAccessToken   []ReqestScope
}

type ConsentForm struct {
	Challenge     string   `json:"challenge"`
	RequestScopes []string `json:"scope"`
}

type LogoutPage struct {
	PageTitle         string
	LogoutTitle       string
	LogoutButtonLabel string
	Challenge         string
	LogoutDenyLabel   string
	Subject           string `json:"subject"`
}

type ReqestScope struct {
	ScopeName  string
	ScopeValue string
}
