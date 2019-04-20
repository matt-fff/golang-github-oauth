package main

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		Scopes:       []string{"public_repo"},
		Endpoint:     github.Endpoint,
	}
	url := conf.AuthCodeURL("", oauth2.AccessTypeOffline)
	fmt.Printf("Login with GitHub: %v", url)

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Fatal(err)
		}

		token, err := conf.Exchange(ctx, r.Form["code"][0])
		if err != nil {
			log.Fatal(err)
		}

		client := conf.Client(ctx, token)
		response, err := client.Get("https://api.github.com/user/repos?page=0&per_page=100")
		if err != nil {
			log.Fatal(err)
		}

		defer response.Body.Close()
		repos, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprintf(w, "Repos: %s", repos)
	})

	log.Fatal(http.ListenAndServe(":4567", nil))
}
