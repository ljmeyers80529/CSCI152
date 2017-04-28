// package main

// // Trying to run this code will result in an error due to a custom function on my
// // local repo of the Spotify API

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"sort"

// 	"github.com/zmb3/spotify"
// )

// // redirectURI is the OAuth redirect URI for the application.
// // You must register an application at Spotify's developer portal
// // and enter this value.
// const (
// 	redirectURI = "http://localhost:8080/callback"
// 	testID      = "80c614680ee64001a9fe3f5d98880364"
// 	testSecret  = "a3790222803a4f8fbdd5cdd5a2ce64d9"
// 	root        = "https://api.spotify.com/v1/"
// 	authroot    = "https://accounts.spotify.com/authorize"
// )

// var (
// 	auth = spotify.NewAuthenticator(redirectURI, "user-top-read")
// 	//auth  = spotify.NewAuthenticator(redirectURI, "user-read-recently-played")
// 	ch    = make(chan *spotify.Client)
// 	state = "abc123"
// )

// func main() {

// 	// Set SPOTIFY_ID and SPOTIFY_SECRET
// 	auth.SetAuthInfo(testID, testSecret)

// 	// Start HTTP Server
// 	http.HandleFunc("/callback", completeAuth)
// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		log.Println("Got request for:", r.URL.String())
// 	})
// 	go http.ListenAndServe(":8080", nil)

// 	url := auth.AuthURL(state)
// 	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

// 	// wait for auth to complete
// 	client := <-ch

// 	// use the client to make calls that require authorization
// 	user, err := client.CurrentUser()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("You are logged in as:", user.ID)

// 	// data, err := client.GetAudioFeatures("1zHlj4dQ8ZAtrayhuDDmkY")
// 	// data, err := client.CurrentUserRecentTracks(50)
// 	// data, err := client.GetAudioAnalysis("1zHlj4dQ8ZAtrayhuDDmkY")
// 	// data, err := client.CurrentUserTopTracks(50, "short")
// 	// data, err := client.CurrentUserTopArtists(50, "long")

// 	getPersonalizedPlaylist(client)
// 	/*
// 			var trackSeeds []spotify.ID

// 			for _, val := range short.Items {
// 				trackSeeds = append(trackSeeds, val.ID)
// 			}
// 			for _, val := range trackSeeds {
// 				fmt.Println("trackseeds: ", val)
// 			}

// 			// for _, val := range long.Items {
// 			// 	trackSeeds = append(trackSeeds, val.ID)
// 			// }

// 			seeds := spotify.Seeds{Tracks: trackSeeds}
// 			fmt.Println("seeds", seeds)
// 			attr := spotify.NewTrackAttributes()

// 			data, err := client.GetRecommendations(seeds, attr, nil)

// 			// for _, val := range data {
// 			// 	for _, item := range val.Items {
// 			// 		//do something to each track
// 			// 	}
// 			// }

// 		d, err := json.MarshalIndent(data, "", "  ")
// 		fmt.Println("JSON DATA")
// 		fmt.Println(d)
// 		for _, val := range data.Tracks {
// 			fmt.Println(val.Name)
// 		}
// 		err = ioutil.WriteFile("temp.txt", d, 0644)
// 	*/
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// }

// func getPersonalizedPlaylist(client *spotify.Client) (*spotify.SimplePlaylist, error) {
// 	// Gather user's top tracks
// 	limit := 50
// 	timerange := "short_term"
// 	opt := spotify.Options{
// 		Limit:     &limit,
// 		Timerange: &timerange,
// 	}
// 	short, err := client.CurrentUserTopTracks(&opt)
// 	if err != nil {
// 		fmt.Println(err)
// 		return nil, err
// 	}

// 	// Map of the form map[ArtistID][]TrackID
// 	tracksByArtist := make(map[spotify.ID][]spotify.ID)
// 	tracksByName := make(map[string][]string)
// 	var artistIDList []spotify.ID

// 	for _, val := range short.Items {
// 		tracksByArtist[val.Artists[0].ID] = append(tracksByArtist[val.Artists[0].ID], val.ID)
// 		tracksByName[val.Artists[0].Name] = append(tracksByName[val.Artists[0].Name], val.Name)
// 		artistIDList = append(artistIDList, val.Artists[0].ID)
// 	}

// 	for ind, val := range tracksByName {
// 		fmt.Println("Artist:", ind, "\nTrack:", val, "\n")
// 	}

// 	artistObjectList, err := spotify.GetArtists(artistIDList...)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Create a map[genre]count to find distribution of genres for user
// 	countsByGenre := make(map[string]int)
// 	for _, val := range artistObjectList {
// 		for _, genre := range val.Genres {
// 			if _, ok := countsByGenre[genre]; ok { // OK idiom; if key exists in map, execute if statement
// 				countsByGenre[genre] = countsByGenre[genre] + 1
// 			} else {
// 				countsByGenre[genre] = 1 // Key doesnt exist, so initialize
// 			}
// 		}
// 	}

// 	// Sort genres by occurrence
// 	genreBySortedCounts := make(map[int][]string)
// 	var temp []int
// 	for k, v := range countsByGenre {
// 		genreBySortedCounts[v] = append(genreBySortedCounts[v], k)
// 	}

// 	for k := range genreBySortedCounts {
// 		temp = append(temp, k)
// 	}

// 	sort.Sort(sort.Reverse(sort.IntSlice(temp)))

// 	for _, k := range temp {
// 		for _, s := range genreBySortedCounts[k] {
// 			fmt.Printf("%s: %d\n", s, k)
// 		}
// 	}

// 	return nil, nil
// }

// func completeAuth(w http.ResponseWriter, r *http.Request) {
// 	tok, err := auth.Token(state, r)
// 	if err != nil {
// 		http.Error(w, "Couldn't get token", http.StatusForbidden)
// 		log.Fatal(err)
// 	}
// 	if st := r.FormValue("state"); st != state {
// 		http.NotFound(w, r)
// 		log.Fatalf("State mismatch: %s != %s\n", st, state)
// 	}
// 	// use the token to get an authenticated client
// 	client := auth.NewClient(tok)
// 	fmt.Fprintf(w, "Login Completed!")
// 	ch <- &client
// }
