package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	oidc "github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

const (
	ClientID     = "myclient"
	ClientSecret = "43354724-c159-4c16-a888-28e193179c01"
)

func main() {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, "http://localhost:8080/auth/realms/myrealm")
	checkError(err)

	config := oauth2.Config{
		ClientID:     ClientID,
		ClientSecret: ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  "http://localhost:8081/auth/callback",
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "roles"},
	}

	state := strconv.Itoa(rand.Intn(5000) + 1000)

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		http.Redirect(writer, request, config.AuthCodeURL(state), http.StatusFound)
	})

	http.HandleFunc("/auth/callback", func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Query().Get("state") != state {
			http.Error(writer, "State inválido", http.StatusBadRequest)
			return
		}

		token, err := config.Exchange(ctx, request.URL.Query().Get("code"))
		if err != nil {
			http.Error(writer, "Falha ao trocar o token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		resp := struct {
			AccessToken *oauth2.Token
		}{
			AccessToken: token,
		}

		data, err := json.Marshal(resp)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Add("Content-Type", "application/json")
		writer.Write(data)

	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
