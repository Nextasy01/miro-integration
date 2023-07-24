package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	handler "github.com/Nextasy01/miro-integration/api"
	"github.com/gorilla/mux"
)

var (
	l = log.New(os.Stdout, "miro-api ", log.LstdFlags)
)

func main() {
	r := mux.NewRouter()

	postRouter := r.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/api/v1/miro/login", func(w http.ResponseWriter, r *http.Request) {
		l.Println("Handling POST Login Request to Miro")
		reqLogin := &handler.LoginRequest{}
		reqBody := json.NewDecoder(r.Body)
		reqBody.Decode(&reqLogin)

		email := reqLogin.Email
		password := reqLogin.Password

		if email == "" || password == "" {
			http.Error(w, "error: please provide email and password", http.StatusInternalServerError)
			l.Println("no email and password")
			return
		}

		token, err := handler.StartSelenium(email, password)
		if err != nil {
			http.Error(w, "error: unable to execute selenium", http.StatusInternalServerError)
			l.Printf("error occured: %v", err)
			return
		}
		l.Println("Successfully retrieve the token")
		fmt.Fprintf(w, "token: %s", token)
	})

	getRouter := r.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/api/v1/miro/teams", func(w http.ResponseWriter, r *http.Request) {
		l.Println("Handling GET Teams request")
		reqTeam := &handler.TeamRequest{}
		reqBody := json.NewDecoder(r.Body)
		reqBody.Decode(&reqTeam)

		token := reqTeam.Token

		if token == "" {
			http.Error(w, "error: please provide the session token", http.StatusInternalServerError)
			l.Println("no token provided")
			return
		}

		data, err := handler.GetTeamID(token)
		if err != nil {
			http.Error(w, "error: couldn't get response from miro", http.StatusInternalServerError)
			l.Printf("error occured: %v", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		l.Println("Successfully return the list of teams in Miro")
		fmt.Fprint(w, string(data))
	})

	getRouter.HandleFunc("/api/v1/miro/team/users", func(w http.ResponseWriter, r *http.Request) {
		l.Println("Handling GET users in team request")
		reqTeamMembers := &handler.TeamMembersRequest{}
		reqBody := json.NewDecoder(r.Body)
		reqBody.Decode(&reqTeamMembers)

		token := reqTeamMembers.Token
		teamId := reqTeamMembers.TeamID

		data, err := handler.GetTeamMembers(token, teamId)
		if err != nil {
			http.Error(w, "error: couldn't get response from miro", http.StatusInternalServerError)
			l.Printf("error occured: %v", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		l.Println("Successfully returned the list of users")
		fmt.Fprint(w, string(data))
	})

	l.Println("Starting server at 8080 port")

	if os.Getenv("PORT") == "" {
		log.Fatal(http.ListenAndServe(":8080", r))
	}

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), r))

}
