package main

import (
	"context"
	"log"
	"net/http"

	oidc "github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

const (
	ClientID     = "myclient"
	ClientSecret = "8c67e844-bb1f-4f8d-a048-b1171fc2f882"
)

func main() {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, "http://localhost:8080/auth/realms/my-realm")
	checkError(err)

	config := oauth2.Config{
		ClientID:     ClientID,
		ClientSecret: ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  "http://localhost:8081/auth/callback",
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "roles"},
	}

	state := "123"

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		http.Redirect(writer, request, config.AuthCodeURL(state), http.StatusFound)
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
