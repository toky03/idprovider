package adapter

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"user-service/model"
)

func TestHydraAdapter_ReadChallenge(t *testing.T) {

	srv := httptest.NewServer(mockChallenge())
	os.Setenv("HYDRA_URL", srv.URL)
	defer srv.Close()

	adapter := NewHydraAdapter()
	body, err := adapter.ReadChallenge("loginChallenge", "login")
	if err != nil {
		t.Error("error should be nil", err)
		t.FailNow()
	}
	if body.RequestURL != "url" {
		t.Error("false url")
	}

}

func TestHydraAdapter_ReadChallenge_Failure(t *testing.T) {
	srv := httptest.NewServer(mockError())
	os.Setenv("HYDRA_URL", srv.URL)
	defer srv.Close()

	adapter := NewHydraAdapter()
	_, err := adapter.ReadChallenge("loginChallenge", "login")
	if err == nil {
		t.Error("error should not be nil")
		t.FailNow()
	}

}

func mockError() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func mockChallenge() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		loginResponse := model.LoginChallenge{
			Skip:                 false,
			Subject:              "",
			Client:               model.Client{},
			RequestURL:           "url",
			RequestedScope:       nil,
			RequestedAccessToken: nil,
			RpInitiated:          false,
			X:                    nil,
		}

		json.NewEncoder(w).Encode(loginResponse)
	}
}
