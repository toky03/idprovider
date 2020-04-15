package main

import (
	"log"
	"net/http"
	"user-service/manager"
)

func main() {

	handler := manager.NewServieHandler()

	http.HandleFunc("/login", handler.LoginHandler)
	http.HandleFunc("/consent", handler.ConsentHandler)
	http.HandleFunc("/acceptConsent", handler.AcceptConsentHandler)
	http.HandleFunc("/logout", handler.LogoutHandler)

	log.Println("Server is running at 3000 port.")
	http.ListenAndServe(":3000", nil)

}
