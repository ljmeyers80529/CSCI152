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
	dataPerGenreNum = 500
	iterationNum    = 1000
	confidenceNum   = 0.75
	extraNodesNum   = 20
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
	// Get command-line arguments
	choice, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Println(err)
		return
	}

	// Initialize variables assigned by constants
	confidence := confidenceNum
	dataPerGenre := dataPerGenreNum
	extraNodes := extraNodesNum
	iterations := iterationNum

	// Initialize genre variables
	var genres []string
	genres = append(genres, "classical")
	genres = append(genres, "pop")

	// Initialize trainingData
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

		// Block until authorization complete
		client := <-ch

		// Assign current user according to client
		user, err := client.CurrentUser()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("You are logged in as:", user.ID)

		// Generate training data from endpoints
		fmt.Println("Generating genre data...")
		trainingData, err = generateGenreData(client, genres, dataPerGenre)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		// Generate training data from file
		log.Println("Reading genre data...")
		trainingData, err = generateDataFromFile("analyses.txt", "features.txt")
		if err != nil {
			log.Println(err)
			return
		}
	}

	// Network operations
	dataElementCount := len(trainingData[0])
	genreCount := 2
	hiddenCount := dataElementCount + genreCount + extraNodes

	//network := gonn.DefaultNetwork(dataElementCount, hiddenCount, genreCount, false)
	//network := gonn.NewNetwork(dataElementCount, hiddenCount, genreCount, false, 0.01, 0.001) // This is working sort of because of the false regression
	//network := gonn.NewNetwork(dataElementCount, hiddenCount, genreCount, true, 0.0001, 0.00001) // Might be overfitting
	//network := gonn.NewNetwork(dataElementCount, hiddenCount, genreCount, true, 0.0005, 0.00005) // This seems to only rarely get bad results, best so far
	network := gonn.NewNetwork(dataElementCount, hiddenCount, genreCount, true, 0.0005, 0.00005)
	targetData := generateTargetData(genreCount, dataPerGenre)

	// Debugging code
	log.Println("TRAINING DATA")
	for index, val := range trainingData {
		log.Println("Index:", index, val)
	}
	log.Println("\nTARGET DATA")
	for index, val := range targetData {
		log.Println("Index:", index, val)
	}

	network.Train(trainingData, targetData, iterations)

	testNetwork(network, trainingData, targetData, genreCount, dataPerGenre, confidence)

	// Store neural network on disk
	gonn.DumpNN("genreclassification.nn", network)
}

func testNetwork(network *gonn.NeuralNetwork, testingData [][]float64, targetData [][]float64, genreCount int, dataPerGenre int, confidence float64) {
	log.Println("\nTesting Neural Network")
	passCount := 0
	total := genreCount * dataPerGenre
	for index, val := range testingData {
		key := index / dataPerGenre
		target := targetData[index]
		result := network.Forward(val)
		if result[key] > confidence {
			fmt.Print("Status: Pass  ")
			passCount++
		} else {
			fmt.Print("Status: FAIL  ")
		}
		fmt.Println("Target: ", target, " Result: ", result)
	}

	log.Println("Accuracy Report")
	successRate := float64(passCount) / float64(total)
	log.Printf("Success Rate: %.2f\n", successRate)
	failureRate := 1.0 - successRate
	log.Printf("Failure Rate: %.2f\n", failureRate)
}

func generateTargetData(genreCount int, dataPerGenre int) (targets [][]float64) {
	total := genreCount * dataPerGenre

	for i := 0; i < total; i++ {
		key := i / dataPerGenre
		target := make([]float64, genreCount)
		target[key] = 1
		targets = append(targets, target)
	}

	return targets
}

func generateGenreData(client *spotify.Client, genres []string, dataPerGenre int) (data [][]float64, err error) {
	seeds := formatSeeds(genres)

	log.Println("Generating recommendations...")
	recs, err := generateRecommendations(client, seeds, dataPerGenre)
	if err != nil {
		return nil, err
	}

	log.Println("Getting ID's...")
	ids := getIDs(recs)
	log.Println(ids)
	log.Println(len(ids))

	log.Println("Getting analyses...")
	analyses, err := getAnalyses(client, ids)
	if err != nil {
		return nil, err
	}
	for index := range analyses {
		log.Println("Index:", index, "analysis")
	}

	log.Println("Getting Features...")
	rawFeatures, err := getFeatures(client, ids)
	if err != nil {
		return nil, err
	}
	features := formatFeatures(rawFeatures)

	log.Println("Writing Files...")
	err = writeFiles("analyses.txt", "features.txt", analyses, rawFeatures)
	if err != nil {
		return nil, err
	}

	for index, val := range features {
		log.Println("Index:", index, val)
	}

	log.Println("Formatting Data...")
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

func getFeatures(client *spotify.Client, ids []spotify.ID) (features []*spotify.AudioFeatures, err error) {
	idCopies := ids // Prevent modifying ID slice
	idsLeftToProcess := len(idCopies)

	for idsLeftToProcess > 100 {
		currentIDs := idCopies[:100]
		currentFeatures, err := client.GetAudioFeatures(currentIDs...)
		if err != nil {
			return nil, err
		}
		features = append(features, currentFeatures...)

		idCopies = idCopies[100:]
		idsLeftToProcess = len(idCopies)
	}

	if idsLeftToProcess > 0 {
		currentFeatures, err := client.GetAudioFeatures(idCopies...)
		if err != nil {
			return nil, err
		}
		features = append(features, currentFeatures...)
	}

	return features, nil
}

func formatFeatures(tracks []*spotify.AudioFeatures) (features [][]float64) {
	for _, val := range tracks {
		var track []float64
		track = append(track, float64(val.Acousticness))
		track = append(track, float64(val.Danceability))
		track = append(track, float64(val.Energy))
		//track = append(track, float64(val.Instrumentalness))
		track = append(track, float64(val.Liveness))
		track = append(track, float64(val.Speechiness))
		track = append(track, float64(val.Valence))
		features = append(features, track)
	}
	return features
}

func getAnalyses(client *spotify.Client, ids []spotify.ID) (analyses []*spotify.AudioAnalysis, err error) {
	for index, val := range ids {
		log.Println("\tGetting analysis", index, "for", val)
		analysis, err := client.GetAudioAnalysis(val)
		if err != nil {
			return nil, err
		}
		analyses = append(analyses, analysis)
		log.Println("\tAnalysis successful.")
	}
	log.Println("\tReturning analyses")

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

// I/O Functions

func writeFiles(analysisFileName string, featureFileName string, analyses []*spotify.AudioAnalysis, features []*spotify.AudioFeatures) (err error) {
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

// Network Functions

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
