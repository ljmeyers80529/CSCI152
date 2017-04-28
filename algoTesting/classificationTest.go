package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/fxsjy/gonn/gonn"
	"github.com/zmb3/spotify"
)

const (
	redirectURI     = "http://localhost:8080/callback"
	testID          = "80c614680ee64001a9fe3f5d98880364"
	testSecret      = "a3790222803a4f8fbdd5cdd5a2ce64d9"
	root            = "https://api.spotify.com/v1/"
	authroot        = "https://accounts.spotify.com/authorize"
	dataAmount      = 100
	iterationAmount = 2500
)

var (
	auth  = spotify.NewAuthenticator(redirectURI, "user-top-read")
	ch    = make(chan *spotify.Client)
	state = "abc123"
)

// Command line arguments are as follows:
// 1 - Generate training data using endpoint and save to local file
// 2 - Generate training data using existing local file
func main() {
	var choice int
	choice, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	// Specify genres
	var genres []string
	genres = append(genres, "classical")
	genres = append(genres, "pop")

	var trainingData [][]float64

	// Set up local server
	if choice == 1 {
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

		// Generate Data
		trainingData, err = generateGenreData(client, genres)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		trainingData, err = generateDataFromFile("analyses.txt", "features.txt")
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	dataElementCount := len(trainingData[0])
	genreCount := 2
	// hiddenCount := (dataElementCount + genreCount) / 2
	hiddenCount := dataElementCount + genreCount + 100

	//network := gonn.NewNetwork(dataElementCount, hiddenCount, genreCount, false, 0.01, 0.001) // This is working sort of because of the false regression
	network := gonn.NewNetwork(dataElementCount, hiddenCount, genreCount, false, 0.0001, 0.00001) // Might be overfitting
	//network := gonn.NewNetwork(dataElementCount, hiddenCount, genreCount, true, 0.0001, 0.00001) // Might be overfitting

	//network := gonn.DefaultNetwork(dataElementCount, hiddenCount, genreCount, false)
	targets := generateTargets()

	fmt.Println("DATA")
	for index, val := range trainingData {
		fmt.Println("Index:", index, val)
	}
	fmt.Println("\nTARGETS")
	for index, val := range targets {
		fmt.Println("Index:", index, val)
	}

	network.Train(trainingData, targets, iterationAmount)

	// for _, val := range trainingData {
	// 	fmt.Println(network.Forward(val))
	// }
	testNetwork(trainingData, network)

	// Store neural network on disk
	gonn.DumpNN("genreclassification.nn", network)

	// Test neural network
	// testingData, err := generateGenreData(client, genres)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// testNetwork(testingData, network)
}

func testNetwork(testingData [][]float64, network *gonn.NeuralNetwork) {
	fmt.Println("\nTesting Neural Network")
	for index, val := range testingData {
		result := network.Forward(val)
		if index < dataAmount/2 {
			fmt.Print("Key: {1, 0} ", result)
			if result[0] > 0.5 {
				fmt.Println("\tPass")
			} else {
				fmt.Println("\tFail")
			}
		} else {
			fmt.Print("Key: {0, 1} ", result)
			if result[1] > 0.5 {
				fmt.Println("\tPass")
			} else {
				fmt.Println("\tFail")
			}
		}
	}
}

func generateTargets() [][]float64 {
	targets := make([][]float64, dataAmount)
	for i := 0; i < dataAmount/2; i++ {
		targets[i] = []float64{1.0, 0.0}
	}
	for j := dataAmount / 2; j < dataAmount; j++ {
		targets[j] = []float64{0.0, 1.0}
	}

	// targets := make([][]float64, 40)

	// for i := 0; i < 20; i++ {
	// 	targets[i] = []float64{0.0}
	// }

	// for j := 20; j < 40; j++ {
	// 	targets[j] = []float64{1.0}
	// }

	return targets
}

func generateGenreData(client *spotify.Client, genres []string) (data [][]float64, err error) {
	seeds := formatSeeds(genres)

	recs, err := generateRecommendations(client, seeds, dataAmount/2)
	if err != nil {
		return nil, err
	}

	ids := getIDs(recs)

	err = writeFiles(client, "analyses.txt", "features.txt", ids)
	if err != nil {
		return nil, err
	}

	analyses, err := getAnalyses(client, ids)
	if err != nil {
		return nil, err
	}

	features, err := getFeatures(client, ids)
	if err != nil {
		return nil, err
	}

	data = formatData(analyses, features)
	return data, nil
}

func generateDataFromFile(analysisFileName string, featureFileName string) (data [][]float64, err error) {
	analyses, rawFeatures, err := readFiles(analysisFileName, featureFileName)
	if err != nil {
		return nil, err
	}

	features := formatFeatures(rawFeatures)

	data = formatData(analyses, features)
	return data, nil
}

func formatData(analyses []*spotify.AudioAnalysis, features [][]float64) (data [][]float64) {
	for index, val := range analyses {
		var datum []float64
		datum = append(datum, val.TrackInfo.Duration)
		datum = append(datum, val.TrackInfo.Tempo)
		//datum = append(datum, float64(val.TrackInfo.TimeSignature))
		//datum = append(datum, float64(val.TrackInfo.Key))
		datum = append(datum, val.TrackInfo.Loudness)
		//datum = append(datum, float64(val.TrackInfo.Mode))

		datum = append(datum, features[index]...) // Concatenate datum and features
		data = append(data, datum)
	}
	return data
}

func getFeatures(client *spotify.Client, ids []spotify.ID) (features [][]float64, err error) {
	tracks, err := client.GetAudioFeatures(ids...)
	if err != nil {
		return nil, err
	}
	features = formatFeatures(tracks)
	// for _, val := range tracks {
	// 	var track []float64
	// 	track = append(track, float64(val.Acousticness))
	// 	track = append(track, float64(val.Danceability))
	// 	track = append(track, float64(val.Energy))
	// 	track = append(track, float64(val.Instrumentalness))
	// 	track = append(track, float64(val.Liveness))
	// 	track = append(track, float64(val.Speechiness))
	// 	track = append(track, float64(val.Valence))
	// 	features = append(features, track)
	// }
	return features, nil
}

func formatFeatures(tracks []*spotify.AudioFeatures) (features [][]float64) {
	for _, val := range tracks {
		var track []float64
		track = append(track, float64(val.Acousticness))
		track = append(track, float64(val.Danceability))
		track = append(track, float64(val.Energy))
		track = append(track, float64(val.Instrumentalness))
		track = append(track, float64(val.Liveness))
		track = append(track, float64(val.Speechiness))
		track = append(track, float64(val.Valence))
		features = append(features, track)
	}
	return features
}

func getAnalyses(client *spotify.Client, ids []spotify.ID) (analyses []*spotify.AudioAnalysis, err error) {
	for _, val := range ids {
		analysis, err := client.GetAudioAnalysis(val)
		if err != nil {
			return nil, err
		}
		analyses = append(analyses, analysis)
	}

	return analyses, nil
}

func getIDs(recs []*spotify.Recommendations) (ids []spotify.ID) {
	for _, val := range recs {
		tracks := val.Tracks
		for _, track := range tracks {
			ids = append(ids, track.ID)
		}
	}
	return ids
}

func generateRecommendations(client *spotify.Client, seeds []spotify.Seeds, limit int) (recs []*spotify.Recommendations, err error) {
	//limit := dataAmount / 2
	attr := spotify.NewTrackAttributes()
	opts := spotify.Options{Limit: &limit}

	for _, val := range seeds {
		newRec, err := client.GetRecommendations(val, attr, &opts)
		if err != nil {
			return nil, err
		}
		recs = append(recs, newRec)
	}

	return recs, nil
}

func formatSeeds(genres []string) (seeds []spotify.Seeds) {
	for _, val := range genres {
		var values []string
		values = append(values, val)
		newSeed := spotify.Seeds{Genres: values}
		seeds = append(seeds, newSeed)
	}

	return seeds
}

func trainNetwork() {
	nn := gonn.DefaultNetwork(2, 3, 1, true)
	inputs := [][]float64{
		[]float64{0, 0},
		[]float64{0, 1},
		[]float64{1, 0},
		[]float64{1, 1},
	}

	targets := [][]float64{
		[]float64{0}, //0+0=0
		[]float64{1}, //0+1=1
		[]float64{1}, //1+0=1
		[]float64{2}, //1+1=2
	}

	nn.Train(inputs, targets, 1000)

	for _, p := range inputs {
		fmt.Println(nn.Forward(p))
	}
}

func writeFiles(client *spotify.Client, analysisFileName string, featureFileName string, ids []spotify.ID) error {
	analyses, err := getAnalyses(client, ids)
	if err != nil {
		return err
	}

	features, err := client.GetAudioFeatures(ids...)
	if err != nil {
		return err
	}

	err = writeAnalysesFile(analysisFileName, analyses)
	if err != nil {
		return err
	}

	err = writeFeaturesFile(featureFileName, features)
	if err != nil {
		return err
	}

	return nil
}

func readFiles(analysisFileName string, featureFileName string) (analyses []*spotify.AudioAnalysis, features []*spotify.AudioFeatures, err error) {
	analyses, err = readAnalysisFile(analysisFileName)
	if err != nil {
		return nil, nil, err
	}
	features, err = readFeaturesFile(featureFileName)
	if err != nil {
		return nil, nil, err
	}
	return analyses, features, nil
}

func readAnalysisFile(fileName string) (analyses []*spotify.AudioAnalysis, err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)

	for decoder.More() {
		var analysis spotify.AudioAnalysis
		err = decoder.Decode(&analysis)
		if err != nil {
			return nil, err
		}
		analyses = append(analyses, &analysis)
	}

	return analyses, err
}

func readFeaturesFile(fileName string) (features []*spotify.AudioFeatures, err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)

	for decoder.More() {
		var feature spotify.AudioFeatures
		err = decoder.Decode(&feature)
		if err != nil {
			return nil, err
		}
		features = append(features, &feature)
	}

	return features, err
}

func writeAnalysesFile(fileName string, analyses []*spotify.AudioAnalysis) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	encoder := json.NewEncoder(writer)
	for _, val := range analyses {
		err := encoder.Encode(val)
		if err != nil {
			return err
		}
	}
	writer.Flush()
	return nil
}

func writeFeaturesFile(fileName string, features []*spotify.AudioFeatures) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	encoder := json.NewEncoder(writer)
	for _, val := range features {
		err := encoder.Encode(val)
		if err != nil {
			return err
		}
	}
	writer.Flush()
	return nil
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
