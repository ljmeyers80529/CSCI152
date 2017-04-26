// package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"

// 	"github.com/fxsjy/gonn/gonn"
// 	"github.com/zmb3/spotify"
// )

// const (
// 	redirectURI = "http://localhost:8080/callback"
// 	testID      = "80c614680ee64001a9fe3f5d98880364"
// 	testSecret  = "a3790222803a4f8fbdd5cdd5a2ce64d9"
// 	root        = "https://api.spotify.com/v1/"
// 	authroot    = "https://accounts.spotify.com/authorize"
// )

// var (
// 	auth  = spotify.NewAuthenticator(redirectURI, "user-top-read")
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

// 	var genres []string
// 	genres = append(genres, "classical")
// 	genres = append(genres, "pop")

// 	data, err := generateTrainingData(client, genres)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	dataElementCount := len(data[0])
// 	genreCount := 2
// 	hiddenCount := dataElementCount + genreCount + 10

// 	network := gonn.NewNetwork(dataElementCount, hiddenCount, genreCount, false, 0.01, 0.001) // This is working sort of

// 	targets := make([][]float64, 40)
// 	for i := 0; i < 20; i++ {
// 		targets[i] = []float64{1.0, 0.0}
// 	}
// 	for j := 20; j < 40; j++ {
// 		targets[j] = []float64{0.0, 1.0}
// 	}
// 	// targets := make([][]float64, 40)

// 	// for i := 0; i < 20; i++ {
// 	// 	targets[i] = []float64{0.0}
// 	// }

// 	// for j := 20; j < 40; j++ {
// 	// 	targets[j] = []float64{1.0}
// 	// }

// 	fmt.Println("DATA")
// 	for index, val := range data {
// 		fmt.Println("Index:", index, val)
// 	}
// 	fmt.Println("\nTARGETS")
// 	for index, val := range targets {
// 		fmt.Println("Index:", index, val)
// 	}

// 	network.Train(data, targets, 1000)

// 	for _, val := range data {
// 		fmt.Println(network.Forward(val))
// 	}
// }

// func generateTrainingData(client *spotify.Client, genres []string) (data [][]float64, err error) {
// 	seeds := formatSeeds(genres)

// 	recs, err := generateRecommendations(client, seeds)
// 	if err != nil {
// 		return nil, err
// 	}

// 	ids := getIDs(recs)

// 	analyses, err := getAnalyses(client, ids)
// 	if err != nil {
// 		return nil, err
// 	}

// 	features, err := getFeatures(client, ids)
// 	if err != nil {
// 		return nil, err
// 	}

// 	data = formatData(analyses, features)
// 	return data, nil
// }

// func formatData(analyses []*spotify.AudioAnalysis, features [][]float64) (data [][]float64) {
// 	for index, val := range analyses {
// 		var datum []float64
// 		datum = append(datum, val.TrackInfo.Duration)
// 		datum = append(datum, val.TrackInfo.Tempo)
// 		datum = append(datum, float64(val.TrackInfo.TimeSignature))
// 		datum = append(datum, float64(val.TrackInfo.Key))
// 		datum = append(datum, val.TrackInfo.Loudness)
// 		datum = append(datum, float64(val.TrackInfo.Mode))

// 		datum = append(datum, features[index]...) // Concatenate datum and features
// 		data = append(data, datum)
// 	}
// 	return data
// }

// func getFeatures(client *spotify.Client, ids []spotify.ID) (features [][]float64, err error) {
// 	tracks, err := client.GetAudioFeatures(ids...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	for _, val := range tracks {
// 		var track []float64
// 		track = append(track, float64(val.Acousticness))
// 		track = append(track, float64(val.Danceability))
// 		track = append(track, float64(val.Energy))
// 		track = append(track, float64(val.Instrumentalness))
// 		track = append(track, float64(val.Liveness))
// 		track = append(track, float64(val.Speechiness))
// 		track = append(track, float64(val.Valence))
// 		features = append(features, track)
// 	}
// 	return features, nil
// }

// func getAnalyses(client *spotify.Client, ids []spotify.ID) (analyses []*spotify.AudioAnalysis, err error) {
// 	for _, val := range ids {
// 		analysis, err := client.GetAudioAnalysis(val)
// 		if err != nil {
// 			return nil, err
// 		}
// 		analyses = append(analyses, analysis)
// 	}

// 	return analyses, nil
// }

// func getIDs(recs []*spotify.Recommendations) (ids []spotify.ID) {
// 	for _, val := range recs {
// 		tracks := val.Tracks
// 		for _, track := range tracks {
// 			ids = append(ids, track.ID)
// 		}
// 	}
// 	return ids
// }

// func generateRecommendations(client *spotify.Client, seeds []spotify.Seeds) (recs []*spotify.Recommendations, err error) {
// 	attr := spotify.NewTrackAttributes()
// 	opts := spotify.Options{}

// 	for _, val := range seeds {
// 		newRec, err := client.GetRecommendations(val, attr, &opts)
// 		if err != nil {
// 			return nil, err
// 		}
// 		recs = append(recs, newRec)
// 	}

// 	return recs, nil
// }

// func formatSeeds(genres []string) (seeds []spotify.Seeds) {
// 	for _, val := range genres {
// 		var values []string
// 		values = append(values, val)
// 		newSeed := spotify.Seeds{Genres: values}
// 		seeds = append(seeds, newSeed)
// 	}

// 	return seeds
// }

// func trainNetwork() {
// 	nn := gonn.DefaultNetwork(2, 3, 1, true)
// 	inputs := [][]float64{
// 		[]float64{0, 0},
// 		[]float64{0, 1},
// 		[]float64{1, 0},
// 		[]float64{1, 1},
// 	}

// 	targets := [][]float64{
// 		[]float64{0}, //0+0=0
// 		[]float64{1}, //0+1=1
// 		[]float64{1}, //1+0=1
// 		[]float64{2}, //1+1=2
// 	}

// 	nn.Train(inputs, targets, 1000)

// 	for _, p := range inputs {
// 		fmt.Println(nn.Forward(p))
// 	}
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
