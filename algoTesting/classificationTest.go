// +build ignore
package csci152

// package main

// import (
// 	"bufio"
// 	"encoding/gob"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"
// 	"strings"
// 	"time"

// 	"io"

// 	"math"

// 	"github.com/fxsjy/gonn/gonn"
// 	"github.com/mdesenfants/gokmeans"
// 	"github.com/zmb3/spotify"
// )

// const (
// 	redirectURI         = "http://localhost:8080/callback"
// 	testID              = "80c614680ee64001a9fe3f5d98880364"
// 	testSecret          = "a3790222803a4f8fbdd5cdd5a2ce64d9"
// 	root                = "https://api.spotify.com/v1/"
// 	authroot            = "https://accounts.spotify.com/authorize"
// 	dataPerGenreNum     = 300 // 500 breaks something; using 300 for standard
// 	iterationNum        = 1000
// 	confidenceNum       = 0.50
// 	extraNodesNum       = 5
// 	testDataSizeNum     = 10
// 	singletonNum        = 3
// 	setNum              = 5
// 	kmeansIterationsNum = 20
// )

// var (
// 	auth  = spotify.NewAuthenticator(redirectURI, "user-top-read")
// 	ch    = make(chan *spotify.Client)
// 	state = "abc123"
// )

// var logger *log.Logger

// // Command line arguments are as follows:
// // -network : Explicitly use network to generate genre data
// // -test : Run a neural network test using small test cases generated using network
// // -write : Write the analysis data to disk
// // -kmeans : Calculate kmeans instead of reading from file
// // Default: Use local training data and kmeans from files on disk, don't run test, don't write new files
// func main() {
// 	logf, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE, 0640)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	log.SetOutput(logf)
// 	logger = log.New(io.MultiWriter(logf, os.Stdout), "logger:", log.LstdFlags)
// 	startTime := time.Now()

// 	// Get command-line arguments
// 	argList := []string{"-network", "-test", "-write", "-standardize", "-kmeans"}
// 	args := initializeArgMap(argList)
// 	resolveCommandArgs(args)

// 	// Initialize variables assigned by constants
// 	confidence := confidenceNum
// 	dataPerGenre := dataPerGenreNum
// 	extraNodes := extraNodesNum
// 	iterations := iterationNum
// 	testDataSize := testDataSizeNum
// 	singletonCount := singletonNum
// 	setCount := setNum
// 	kmeansIterations := kmeansIterationsNum

// 	// Initialize genre variables
// 	genres := []string{"classical", "pop", "rock", "electronic"}

// 	// Initialize trainingData and error
// 	var trainingData [][]float64
// 	var client *spotify.Client

// 	// Set up local server
// 	if args["-network"] {
// 		client, err = initializeServer()
// 		if err != nil {
// 			logger.Println(err)
// 			return
// 		}

// 		// Generate training data from endpoints
// 		logger.Println("Generating genre data...")
// 		trainingData, err = generateGenreData(client, genres, dataPerGenre, args["-write"], singletonCount, setCount, kmeansIterations)
// 		if err != nil {
// 			logger.Println(err)
// 			return
// 		}
// 	} else {
// 		// Generate training data from file
// 		//fileList := []string{"analyses.txt", "features.txt", "tatumCentroids"}
// 		logger.Println("Reading genre data...")
// 		trainingData, err = generateDataFromFile("analyses.txt", "features.txt", "centroids.txt", args["-write"], args["-kmeans"], singletonCount, setCount, kmeansIterations)
// 		if err != nil {
// 			logger.Println(err)
// 			return
// 		}
// 	}

// 	// Data standardization
// 	if args["-standardize"] {
// 		temp := trainingData
// 		trainingData = standardizeData(temp)
// 	}

// 	// Network operations
// 	dataElementCount := len(trainingData[0])
// 	genreCount := len(genres)
// 	hiddenCount := dataElementCount + genreCount + extraNodes

// 	//network := gonn.DefaultNetwork(dataElementCount, hiddenCount, genreCount, false)
// 	//network := gonn.NewNetwork(dataElementCount, hiddenCount, genreCount, false, 0.01, 0.001) // This is working sort of because of the false regression
// 	//network := gonn.NewNetwork(dataElementCount, hiddenCount, genreCount, true, 0.0001, 0.00001) // Might be overfitting
// 	//network := gonn.NewNetwork(dataElementCount, hiddenCount, genreCount, true, 0.0005, 0.00005) // This seems to only rarely get bad results, best so far
// 	network := gonn.NewNetwork(dataElementCount, hiddenCount, genreCount, true, 0.001, 0.0001)
// 	targetData := generateTargetData(genreCount, dataPerGenre)

// 	// Debugging code
// 	logger.Println("TRAINING DATA")
// 	for index, val := range trainingData {
// 		logger.Println("Index:", index, val)
// 	}
// 	// logger.Println("\nTARGET DATA")
// 	// for index, val := range targetData {
// 	// 	logger.Println("Index:", index, val)
// 	// }

// 	network.Train(trainingData, targetData, iterations)

// 	// Store neural network on disk
// 	gonn.DumpNN("genreclassification.nn", network)

// 	testNetwork(network, trainingData, targetData, genreCount, dataPerGenre, confidence)

// 	if args["-test"] {

// 		if !args["-network"] {
// 			client, err = initializeServer()
// 			if err != nil {
// 				logger.Println(err)
// 				return
// 			}
// 		}

// 		testingData, err := generateGenreData(client, genres, testDataSize, false, singletonCount, setCount, kmeansIterations)
// 		if err != nil {
// 			logger.Println(err)
// 			return
// 		}

// 		testingTargets := generateTargetData(genreCount, testDataSize)
// 		testNetwork(network, testingData, testingTargets, genreCount, 10, confidence)
// 	}
// 	elapsedTime := time.Since(startTime)
// 	logger.Printf("Execution time: %s", elapsedTime)
// }

// func testNetwork(network *gonn.NeuralNetwork, testingData [][]float64, targetData [][]float64, genreCount int, dataPerGenre int, confidence float64) {
// 	logger.Println("\nTesting Neural Network")
// 	passCount := 0
// 	total := genreCount * dataPerGenre
// 	for index, val := range testingData {
// 		key := index / dataPerGenre
// 		target := targetData[index]
// 		result := network.Forward(val)
// 		if result[key] > confidence {
// 			logger.Print("Status: Pass  ")
// 			passCount++
// 		} else {
// 			logger.Print("Status: FAIL  ")
// 		}
// 		logger.Println("Target: ", target, " Result: ", result)
// 	}

// 	logger.Println("Accuracy Report")
// 	successRate := float64(passCount) / float64(total)
// 	logger.Printf("Success Rate: %.2f\n", successRate)
// 	failureRate := 1.0 - successRate
// 	logger.Printf("Failure Rate: %.2f\n", failureRate)
// }

// // Data Generation

// func generateCentroids(analyses []*spotify.AudioAnalysis, singletonClusterCount, setClusterCount, algoIterations int) (rawCentroids [][][]float64, err error) {
// 	logger.Println("Generating centroids...")
// 	tatumNodes := generateTatumNodes(analyses)
// 	tatumCentroids, err := generateGenericCentroids(tatumNodes, singletonClusterCount, algoIterations)
// 	if err != nil {
// 		return nil, err
// 	}
// 	rawCentroids = append(rawCentroids, tatumCentroids)

// 	timbreNodes := generateTimbreNodes(analyses)
// 	timbreCentroids, err := generateGenericCentroids(timbreNodes, setClusterCount, algoIterations)
// 	if err != nil {
// 		return nil, err
// 	}
// 	rawCentroids = append(rawCentroids, timbreCentroids)

// 	pitchNodes := generatePitchNodes(analyses)
// 	pitchCentroids, err := generateGenericCentroids(pitchNodes, setClusterCount, algoIterations)
// 	if err != nil {
// 		return nil, err
// 	}
// 	rawCentroids = append(rawCentroids, pitchCentroids)

// 	return rawCentroids, nil
// }

// func formatCentroids(rawCentroids [][][]float64) [][]float64 {
// 	logger.Println("Formatting centroids...")

// 	height := len(rawCentroids[0])
// 	centroids := initializeSlice(0, height)

// 	for index := range centroids {
// 		for _, val := range rawCentroids {
// 			centroids[index] = append(centroids[index], val[index]...)
// 		}
// 	}
// 	return centroids
// }

// // generateGenericCentroids takes a 2D slice of kmeans nodes and runs the kmeans algorithm on each slice of nodes,
// // for a set amount of round denoted by the iterations variable, yielding the specified amount of centroid clusters
// // denoted by clusterCount.
// func generateGenericCentroids(nodes [][]gokmeans.Node, clusterCount int, iterations int) (centroids [][]float64, err error) {
// 	for index, val := range nodes {
// 		success, rawCentroids := gokmeans.Train(val, clusterCount, iterations)
// 		if !success {
// 			err = errors.New("centroid training has failed")
// 			return nil, err
// 		}
// 		logger.Println("Success! Displaying centroids for track ", index)
// 		for _, c := range rawCentroids {
// 			logger.Println(c)
// 		}
// 		var floats []float64
// 		for _, rawNode := range rawCentroids {
// 			floats = append(floats, rawNode...)
// 		}
// 		centroids = append(centroids, floats)
// 	}

// 	return centroids, err
// }

// // generateTatumNodes parses the provided list of analysis objects and extracts that tatum data
// // for each track as a slice of slices of kmeans nodes.
// func generateTatumNodes(analyses []*spotify.AudioAnalysis) (tatumNodes [][]gokmeans.Node) {
// 	for _, val := range analyses {
// 		var nodes []gokmeans.Node
// 		for _, tatum := range val.Tatums {
// 			tempNode := gokmeans.Node{tatum.Duration}
// 			nodes = append(nodes, tempNode)
// 		}
// 		tatumNodes = append(tatumNodes, nodes)
// 	}
// 	return tatumNodes
// }

// // generateTimbreNodes parses the provided list of analysis objects and extracts a sample of the timbre data
// // for each track as a slice of slices of kmeans nodes.
// func generateTimbreNodes(analyses []*spotify.AudioAnalysis) (timbreNodes [][]gokmeans.Node) {
// 	for _, val := range analyses {
// 		var nodes []gokmeans.Node
// 		for index, segment := range val.Segments {
// 			mod := index % 10
// 			if mod == 0 || mod == 5 {
// 				var tempNode gokmeans.Node
// 				tempNode = segment.Timbre
// 				nodes = append(nodes, tempNode)
// 			}
// 		}
// 		timbreNodes = append(timbreNodes, nodes)
// 	}
// 	return timbreNodes
// }

// // generatePitchNodes parses the provided list of analysis objects and extracts a sample of the timbre data
// // for each track as a slice of slices of kmeans nodes.
// func generatePitchNodes(analyses []*spotify.AudioAnalysis) (pitchNodes [][]gokmeans.Node) {
// 	for _, val := range analyses {
// 		var nodes []gokmeans.Node
// 		for index, segment := range val.Segments {
// 			mod := index % 10
// 			if mod == 0 || mod == 5 {
// 				var tempNode gokmeans.Node
// 				tempNode = segment.Pitches
// 				nodes = append(nodes, tempNode)
// 			}
// 		}
// 		pitchNodes = append(pitchNodes, nodes)
// 	}
// 	return pitchNodes
// }

// // generateTargetData takes the amount of genres and the data needed per genre
// // then returns a slice of slices of float64s that denote the target data for the given
// // parameters
// func generateTargetData(genreCount int, dataPerGenre int) (targets [][]float64) {
// 	total := genreCount * dataPerGenre

// 	for i := 0; i < total; i++ {
// 		key := i / dataPerGenre
// 		target := make([]float64, genreCount)
// 		// for index := range target {
// 		// 	target[index] = -1
// 		// }
// 		target[key] = 1
// 		targets = append(targets, target)
// 	}

// 	return targets
// }
// func generateGenreData(client *spotify.Client, genres []string, dataPerGenre int, shouldWrite bool, singleCount, setCount, iterations int) (genreData [][]float64, err error) {
// 	seeds := getSeeds(genres)

// 	recs, err := generateRecommendations(client, seeds, dataPerGenre)
// 	if err != nil {
// 		return nil, err
// 	}

// 	ids := getIDs(recs)

// 	rawAnalyses, err := getAnalyses(client, ids)
// 	if err != nil {
// 		return nil, err
// 	}
// 	rawCentroids, err := generateCentroids(rawAnalyses, singleCount, setCount, iterations)
// 	if err != nil {
// 		return nil, err
// 	}
// 	rawFeatures, err := getFeatures(client, ids)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if shouldWrite {
// 		err = writeFiles("analyses.txt", "features.txt", "centroids.txt", rawAnalyses, rawFeatures, rawCentroids)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	//
// 	for index := range rawAnalyses {
// 		logger.Println("Index:", index, "analysis")
// 	} //
// 	centroids := formatCentroids(rawCentroids)

// 	analyses := formatAnalyses(rawAnalyses, centroids)

// 	features := formatFeatures(rawFeatures)
// 	for index, val := range features {
// 		logger.Println("Index:", index, val)
// 	}

// 	genreData = formatData(analyses, features)
// 	return genreData, nil
// }

// func generateDataFromFile(analysisFileName, featureFileName, centroidFileName string, shouldWrite, shouldCalc bool, singleCount, setCount, iterations int) (genreData [][]float64, err error) {
// 	rawAnalyses, rawFeatures, rawCentroids, err := readFiles(analysisFileName, featureFileName, centroidFileName)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if shouldCalc {
// 		rawCentroids, err = generateCentroids(rawAnalyses, singleCount, setCount, iterations)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	centroids := formatCentroids(rawCentroids)

// 	analyses := formatAnalyses(rawAnalyses, centroids)
// 	features := formatFeatures(rawFeatures)

// 	if shouldWrite {
// 		err = writeCentroidsFile("centroids.txt", rawCentroids)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	genreData = formatData(analyses, features)
// 	return genreData, nil
// }

// func formatData(analyses [][]float64, features [][]float64) (data [][]float64) {
// 	logger.Println("Formatting genre data...")

// 	for index, val := range analyses {
// 		var datum []float64
// 		datum = append(datum, val...)
// 		datum = append(datum, features[index]...)

// 		data = append(data, datum)
// 	}

// 	return data
// }

// func getFeatures(client *spotify.Client, ids []spotify.ID) (features []*spotify.AudioFeatures, err error) {
// 	logger.Println("Getting features...")

// 	idCopies := ids // Prevent modifying ID slice
// 	idsLeftToProcess := len(idCopies)

// 	for idsLeftToProcess > 100 {
// 		currentIDs := idCopies[:100]
// 		currentFeatures, err := client.GetAudioFeatures(currentIDs...)
// 		if err != nil {
// 			return nil, err
// 		}
// 		features = append(features, currentFeatures...)

// 		idCopies = idCopies[100:]
// 		idsLeftToProcess = len(idCopies)
// 	}

// 	if idsLeftToProcess > 0 {
// 		currentFeatures, err := client.GetAudioFeatures(idCopies...)
// 		if err != nil {
// 			return nil, err
// 		}
// 		features = append(features, currentFeatures...)
// 	}

// 	return features, nil
// }

// // standardizeData will transform the raw data provided into a standardized set of float64 values
// // using gaussian normalization
// func standardizeData(rawData [][]float64) [][]float64 {
// 	rowWidth := len(rawData[0])
// 	colHeight := len(rawData)
// 	normData := initializeSlice(rowWidth, colHeight)

// 	for i := 0; i < rowWidth; i++ {
// 		total := 0.0
// 		for j := 0; j < colHeight; j++ {
// 			total += rawData[j][i]
// 		}
// 		mean := total / float64(colHeight)
// 		deviation := 0.0
// 		var differences []float64
// 		for j := 0; j < colHeight; j++ {
// 			difference := rawData[j][i] - mean
// 			differences = append(differences, difference)
// 			deviation += math.Pow(difference, 2)
// 		}
// 		stdDeviation := math.Sqrt(deviation)
// 		for j := 0; j < colHeight; j++ {
// 			normData[j][i] = (differences[j]) / stdDeviation
// 		}
// 	}

// 	return normData
// }

// // normalizeValue performs gaussian normalization on the specified value with the provided
// // mean and standard deviation
// func normalizeValue(val, mean, dev float64) float64 {
// 	return (val - mean) / dev
// }

// // initializeSlice simply creates a two-dimensional slice of float64s with the specified
// // width and height
// func initializeSlice(width int, height int) [][]float64 {
// 	slice := make([][]float64, height)
// 	for index := range slice {
// 		slice[index] = make([]float64, width)
// 	}
// 	return slice
// }

// // formatAnalyses will take the raw analysis objects and centroid data, parse through each and
// // return a 2D slice of float64s with correctly formatted data
// func formatAnalyses(rawAnalyses []*spotify.AudioAnalysis, centroids [][]float64) (analyses [][]float64) {
// 	logger.Println("Formatting analyses...")

// 	for index, val := range rawAnalyses {
// 		var analysis []float64
// 		analysis = append(analysis, val.TrackInfo.Duration)
// 		analysis = append(analysis, val.TrackInfo.Tempo)
// 		//analysis = append(analysis, float64(val.TrackInfo.TimeSignature)) // Might remove
// 		//analysis = append(analysis, float64(val.TrackInfo.Key))           //
// 		analysis = append(analysis, val.TrackInfo.Loudness)
// 		//analysis = append(analysis, float64(val.TrackInfo.Mode)) //
// 		analysis = append(analysis, centroids[index]...)

// 		analyses = append(analyses, analysis)
// 	}
// 	logger.Println("Formatted analyses: ", analyses)
// 	return analyses
// }

// // formatFeatures will take the raw audio features objects, parse though each, and return
// // a 2D slice of float64s with correctly formatted data
// func formatFeatures(tracks []*spotify.AudioFeatures) (features [][]float64) {
// 	logger.Println("Formatting features...")

// 	for _, val := range tracks {
// 		var track []float64
// 		track = append(track, float64(val.Acousticness))
// 		track = append(track, float64(val.Danceability))
// 		track = append(track, float64(val.Energy))
// 		//track = append(track, float64(val.Instrumentalness)) // Might remove
// 		track = append(track, float64(val.Liveness))
// 		track = append(track, float64(val.Speechiness))
// 		track = append(track, float64(val.Valence))
// 		features = append(features, track)
// 	}
// 	return features
// }

// // getAnalyses will use the specified authorized client to make seperate API calls that will return
// // the audio analysis of every song specified by the slice of IDs.
// func getAnalyses(client *spotify.Client, ids []spotify.ID) (analyses []*spotify.AudioAnalysis, err error) {
// 	logger.Println("Getting analyses...")

// 	for index, val := range ids {
// 		logger.Println("Getting analysis", index, "for", val)
// 		analysis, err := client.GetAudioAnalysis(val)
// 		if err != nil {
// 			return nil, err
// 		}
// 		analyses = append(analyses, analysis)
// 		logger.Println("Analysis successful.")
// 	}
// 	logger.Println("Returning analyses")

// 	return analyses, nil
// }

// // getIDs will take a list of raw recommendation objects, parse them, then return a slice
// // of IDs that pertain to each recommended track.
// func getIDs(recs []*spotify.Recommendations) (ids []spotify.ID) {
// 	logger.Println("Getting ID's...")

// 	for _, val := range recs {
// 		tracks := val.Tracks
// 		for _, track := range tracks {
// 			ids = append(ids, track.ID)
// 		}
// 	}
// 	return ids
// }

// // generateRecommendations uses the specified authorized client and the provided slice of seeds
// // to make seperate API calls that return recommendations for the provided seeds.
// func generateRecommendations(client *spotify.Client, seeds []spotify.Seeds, limit int) (recs []*spotify.Recommendations, err error) {
// 	logger.Println("Generating recommendations...")

// 	iterationsPerSeed := limit / 100
// 	remainder := limit % 100
// 	hundred := 100
// 	attr := spotify.NewTrackAttributes()
// 	maxOpts := spotify.Options{Limit: &hundred}
// 	remainderOpts := spotify.Options{Limit: &remainder}

// 	logger.Println("Iterating through seeds")
// 	for _, val := range seeds {
// 		for i := 0; i < iterationsPerSeed; i++ {
// 			newRec, err := client.GetRecommendations(val, attr, &maxOpts)
// 			if err != nil {
// 				return nil, err
// 			}
// 			recs = append(recs, newRec)
// 		}
// 		if remainder > 0 {
// 			newRec, err := client.GetRecommendations(val, attr, &remainderOpts)
// 			if err != nil {
// 				return nil, err
// 			}
// 			recs = append(recs, newRec)
// 		}
// 	}

// 	return recs, nil
// }

// // getSeeds takes a list of genres and returns a slice of seed objects for the specified genres.
// func getSeeds(genres []string) (seeds []spotify.Seeds) {
// 	for _, val := range genres {
// 		var values []string
// 		values = append(values, val)
// 		newSeed := spotify.Seeds{Genres: values}
// 		seeds = append(seeds, newSeed)
// 	}

// 	return seeds
// }

// // I/O Functions

// // initializeArgMap takes a list of expected command-line arguments as strings and
// // returns a map using the arguments as keys and initialized false bools as values.
// func initializeArgMap(argList []string) map[string]bool {
// 	argMap := make(map[string]bool, 3)
// 	for _, val := range argList {
// 		argMap[val] = false
// 	}
// 	return argMap
// }

// // resolveCommandArgs takes an
// func resolveCommandArgs(argMap map[string]bool) {
// 	for _, val := range os.Args {
// 		for index := range argMap {
// 			if strings.EqualFold(index, val) {
// 				argMap[index] = true
// 			}
// 		}
// 	}
// }

// func writeFiles(analysisFileName, featureFileName, centroidFileName string, analyses []*spotify.AudioAnalysis, features []*spotify.AudioFeatures, centroids [][][]float64) (err error) {
// 	logger.Println("Writing files...")

// 	err = writeAnalysesFile(analysisFileName, analyses)
// 	if err != nil {
// 		return err
// 	}

// 	err = writeFeaturesFile(featureFileName, features)
// 	if err != nil {
// 		return err
// 	}

// 	err = writeCentroidsFile(centroidFileName, centroids)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func readFiles(analysisFileName, featureFileName, centroidFileName string) (analyses []*spotify.AudioAnalysis, features []*spotify.AudioFeatures, centroids [][][]float64, err error) {
// 	analyses, err = readAnalysisFile(analysisFileName)
// 	if err != nil {
// 		return nil, nil, nil, err
// 	}
// 	features, err = readFeaturesFile(featureFileName)
// 	if err != nil {
// 		return nil, nil, nil, err
// 	}
// 	centroids, err = readCentroidsFile(centroidFileName)
// 	if err != nil {
// 		return nil, nil, nil, err
// 	}

// 	return analyses, features, centroids, nil
// }

// func readAnalysisFile(fileName string) (analyses []*spotify.AudioAnalysis, err error) {
// 	file, err := os.Open(fileName)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()

// 	reader := bufio.NewReader(file)
// 	decoder := json.NewDecoder(reader)

// 	for decoder.More() {
// 		var analysis spotify.AudioAnalysis
// 		err = decoder.Decode(&analysis)
// 		if err != nil {
// 			return nil, err
// 		}
// 		analyses = append(analyses, &analysis)
// 	}

// 	return analyses, err
// }

// func readFeaturesFile(fileName string) (features []*spotify.AudioFeatures, err error) {
// 	file, err := os.Open(fileName)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()

// 	reader := bufio.NewReader(file)
// 	decoder := json.NewDecoder(reader)

// 	for decoder.More() {
// 		var feature spotify.AudioFeatures
// 		err = decoder.Decode(&feature)
// 		if err != nil {
// 			return nil, err
// 		}
// 		features = append(features, &feature)
// 	}

// 	return features, err
// }

// // readCentroidsFile takes the name of file and reads a slice of [][]float64's that represent
// // precalculated centroids and returns that 3D slice.
// func readCentroidsFile(fileName string) (centroids [][][]float64, err error) {
// 	file, err := os.Open(fileName)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()

// 	decoder := gob.NewDecoder(file)

// 	err = decoder.Decode(&centroids)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return centroids, err
// }

// func writeAnalysesFile(fileName string, analyses []*spotify.AudioAnalysis) error {
// 	file, err := os.Create(fileName)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	writer := bufio.NewWriter(file)
// 	encoder := json.NewEncoder(writer)
// 	for _, val := range analyses {
// 		err := encoder.Encode(val)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	writer.Flush()
// 	return nil
// }

// func writeFeaturesFile(fileName string, features []*spotify.AudioFeatures) error {
// 	file, err := os.Create(fileName)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	writer := bufio.NewWriter(file)
// 	encoder := json.NewEncoder(writer)
// 	for _, val := range features {
// 		err := encoder.Encode(val)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	writer.Flush()
// 	return nil
// }

// func writeCentroidsFile(fileName string, centroids [][][]float64) error {
// 	file, err := os.Create(fileName)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	encoder := gob.NewEncoder(file)
// 	err = encoder.Encode(centroids)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// // Network Functions

// func initializeServer() (client *spotify.Client, err error) {
// 	// Set SPOTIFY_ID and SPOTIFY_SECRET
// 	auth.SetAuthInfo(testID, testSecret)

// 	// Start local HTTP Server
// 	http.HandleFunc("/callback", completeAuth)
// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		logger.Println("Got request for:", r.URL.String())
// 	})
// 	go http.ListenAndServe(":8080", nil)

// 	url := auth.AuthURL(state)
// 	fmt.Println("Please logger in to Spotify by visiting the following page in your browser:", url)

// 	// Block until authorization complete
// 	client = <-ch
// 	return client, nil
// }

// func completeAuth(w http.ResponseWriter, r *http.Request) {
// 	tok, err := auth.Token(state, r)
// 	if err != nil {
// 		http.Error(w, "Couldn't get token", http.StatusForbidden)
// 		logger.Fatal(err)
// 	}
// 	if st := r.FormValue("state"); st != state {
// 		http.NotFound(w, r)
// 		logger.Fatalf("State mismatch: %s != %s\n", st, state)
// 	}
// 	// Retrieve authenticated client using token
// 	client := auth.NewClient(tok)
// 	fmt.Fprintf(w, "Login Completed!")
// 	ch <- &client
// }
