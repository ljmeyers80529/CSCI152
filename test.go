package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/zmb3/spotify"
)

// redirectURI is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.
const (
	redirectURI = "http://localhost:8080/callback"
	testID = "80c614680ee64001a9fe3f5d98880364"
	testSecret = "a3790222803a4f8fbdd5cdd5a2ce64d9"
	root = "https://api.spotify.com/v1/"
	authroot = "https://accounts.spotify.com/authorize"
)

var (
	auth  = spotify.NewAuthenticator(redirectURI, "user-read-recently-played")
	ch    = make(chan *spotify.Client)
	state = "abc123"
)

func main() {

	// Set SPOTIFY_ID and SPOTIFY_SECRET
	auth.SetAuthInfo(testID, testSecret)

	// Start HTTP Server
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go http.ListenAndServe(":8080", nil)
	
	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	// wait for auth to complete
	client := <-ch

	// use the client to make calls that require authorization
	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", user.ID)
	
	//URL := root + "me/player/recently-played" 
	data, _ := client.GetRecent() // Function only available on local copy of spotify api

	/*
	data := make([]interface{}, 20, 100)

	if err = json.Unmarshal(body, &data); err != nil {
		return nil, err 
	}
	*/
	fmt.Println("JSON DATA")
	fmt.Println(data)
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	// use the token to get an authenticated client
	client := auth.NewClient(tok)
	fmt.Fprintf(w, "Login Completed!")
	ch <- &client
}
