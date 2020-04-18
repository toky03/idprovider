package main

import (
	"log"
	"net/http"
	"user-service/manager"
)

func main() {

	loginHandler := manager.NewLoginHandler()
	userHandler := manager.NewUserHandler()

	http.HandleFunc("/login", loginHandler.LoginHandler)
	http.HandleFunc("/consent", loginHandler.ConsentHandler)
	http.HandleFunc("/acceptConsent", loginHandler.AcceptConsentHandler)
	http.HandleFunc("/logout", loginHandler.LogoutHandler)
	http.HandleFunc("/user", userHandler.ManageUser)
	http.HandleFunc("/user/", userHandler.ManageUser)
	http.HandleFunc("/user/application/", userHandler.ManageApplications)

	log.Println("Server is running at 3000 port.")
	http.ListenAndServe(":3000", nil)

}
