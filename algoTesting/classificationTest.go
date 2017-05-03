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
// 	redirectURI     = "http://localhost:8080/callback"
// 	testID          = "80c614680ee64001a9fe3f5d98880364"
// 	testSecret      = "a3790222803a4f8fbdd5cdd5a2ce64d9"
// 	root            = "https://api.spotify.com/v1/"
// 	authroot        = "https://accounts.spotify.com/authorize"
// 	dataPerGenreNum = 300 // 500 breaks something; using 300 for standard
// 	iterationNum    = 2500
// 	confidenceNum   = 0.35
// 	extraNodesNum   = 50
// 	testDataSizeNum = 10
// )

// var (
// 	auth  = spotify.NewAuthenticator(redirectURI, "user-top-read")
// 	ch    = make(chan *spotify.Client)
// 	state = "abc123"
// )

// // Centroid holds the slice of slices of node data for different input data
// type Centroid struct {
// 	Tatums [][]gokmeans.Node
// }

// var logger *log.Logger

// // Command line arguments are as follows:
// // -network : Explicitly use network to generate genre data
// // -test : Run a neural network test using small test cases generated using network
// // -write : Write the analysis data to disk
// // Default: Use local training data from files on disk, don't run test, don't write new files

// func main() {
// 	logf, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE, 0640)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	log.SetOutput(logf)
// 	logger = log.New(io.MultiWriter(logf, os.Stdout), "logger:", log.LstdFlags)
// 	startTime := time.Now()
// 	// Get command-line arguments
// 	argList := []string{"-network", "-test", "-write", "-standardize"}
// 	args := initializeArgMap(argList)
// 	resolveCommandArgs(args)

// 	// Initialize variables assigned by constants
// 	confidence := confidenceNum
// 	dataPerGenre := dataPerGenreNum
// 	extraNodes := extraNodesNum
// 	iterations := iterationNum
// 	testDataSize := testDataSizeNum

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
// 		trainingData, err = generateGenreData(client, genres, dataPerGenre, args["-write"])
// 		if err != nil {
// 			logger.Println(err)
// 			return
// 		}
// 	} else {
// 		// Generate training data from file
// 		logger.Println("Reading genre data...")
// 		trainingData, err = generateDataFromFile("analyses.txt", "features.txt", "tatumCentroids.txt", args["-write"])
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
// 	network := gonn.NewNetwork(dataElementCount, hiddenCount, genreCount, true, 0.0005, 0.00005)
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

// 		testingData, err := generateGenreData(client, genres, testDataSize, false) // Never write test data to file
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

// func formatRawCentroids(rawCentroids [][]gokmeans.Node) (centroids [][]float64) {
// 	for _, val := range rawCentroids {
// 		singleCentroid := val[0]
// 		centroids = append(centroids, singleCentroid)
// 	}

// 	return centroids
// }

// func generateCentroidData(analyses []*spotify.AudioAnalysis) (tatumCentroids [][]gokmeans.Node, err error) {
// 	tatumCentroids, err = generateRawTatumCentroids(analyses)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return tatumCentroids, nil
// }

// func generateRawTatumCentroids(analyses []*spotify.AudioAnalysis) (rawTatumCentroids [][]gokmeans.Node, err error) {
// 	tatumNodes := generateTatumNodes(analyses)
// 	for _, val := range tatumNodes {
// 		success, rawCentroids := gokmeans.Train(val, 2, 25)
// 		if !success {
// 			err = errors.New("centroid training has failed")
// 			return nil, err
// 		}
// 		logger.Println("Success!\nDisplaying centroids")
// 		for _, centroid := range rawCentroids {
// 			logger.Println(centroid)
// 		}
// 		rawTatumCentroids = append(rawTatumCentroids, rawCentroids)
// 	}

// 	return rawTatumCentroids, nil
// }

// // just input all the timbres or maybe a set of the timbres into the kmeans algorithm
// func generatePrimitiveTimbreCentroids(analyses []*spotify.AudioAnalysis) (timbreCentroids [][]float64, err error) {

// 	return timbreCentroids, nil
// }

// func generatePrimitiveTatumCentroids(analyses []*spotify.AudioAnalysis) (tatumCentroids [][]float64, err error) {
// 	tatumNodes := generateTatumNodes(analyses)
// 	for index, val := range tatumNodes {
// 		success, rawCentroids := gokmeans.Train(val, 3, 15)
// 		if !success {
// 			err = errors.New("centroid training has failed")
// 			return nil, err
// 		}
// 		logger.Println("Success! Displaying centroids for track ", index)
// 		for _, centroid := range rawCentroids {
// 			logger.Println(centroid)
// 		}
// 		var centroids []float64
// 		// Loop through the current slice of nodes and extract the underlying floats from each
// 		for _, rawNode := range rawCentroids {
// 			underlyingFloat := []float64(rawNode)[0]
// 			centroids = append(centroids, underlyingFloat)
// 		}

// 		tatumCentroids = append(tatumCentroids, centroids)
// 	}

// 	return tatumCentroids, nil
// }

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

// func generateTargetData(genreCount int, dataPerGenre int) (targets [][]float64) {
// 	total := genreCount * dataPerGenre

// 	for i := 0; i < total; i++ {
// 		key := i / dataPerGenre
// 		target := make([]float64, genreCount)
// 		target[key] = 1
// 		targets = append(targets, target)
// 	}

// 	return targets
// }

// func generateGenreData(client *spotify.Client, genres []string, dataPerGenre int, shouldWrite bool) (genreData [][]float64, err error) {
// 	seeds := formatSeeds(genres)

// 	logger.Println("Generating recommendations...")
// 	recs, err := generateRecommendations(client, seeds, dataPerGenre)
// 	if err != nil {
// 		return nil, err
// 	}

// 	logger.Println("Getting ID's...")
// 	ids := getIDs(recs)
// 	logger.Println(ids)
// 	logger.Println(len(ids))

// 	logger.Println("Getting analyses...")
// 	rawAnalyses, err := getAnalyses(client, ids)
// 	if err != nil {
// 		return nil, err
// 	}
// 	for index := range rawAnalyses {
// 		logger.Println("Index:", index, "analysis")
// 	}

// 	logger.Println("Generating centroids...")
// 	// rawTatumCentroids := generateCentroidData(rawAnalyses)
// 	// tatumCentroids := formatTatumCentroids(rawTatumCentroids)
// 	tatumCentroids, err := generatePrimitiveTatumCentroids(rawAnalyses)

// 	logger.Println("Formatting analyses...")
// 	analyses := formatAnalyses(rawAnalyses, tatumCentroids)

// 	logger.Println("Getting features...")
// 	rawFeatures, err := getFeatures(client, ids)
// 	if err != nil {
// 		return nil, err
// 	}

// 	logger.Println("Formatting features...")
// 	features := formatFeatures(rawFeatures)
// 	for index, val := range features {
// 		logger.Println("Index:", index, val)
// 	}

// 	if shouldWrite {
// 		logger.Println("Writing files...")
// 		err = writeFiles("analyses.txt", "features.txt", "tatumCentroids.txt", rawAnalyses, rawFeatures, tatumCentroids)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	logger.Println("Formatting genre data...")
// 	genreData = formatData(analyses, features)
// 	return genreData, nil
// }

// func generateDataFromFile(analysisFileName string, featureFileName string, tatumCentroidFileName string, shouldWrite bool) (genreData [][]float64, err error) {
// 	rawAnalyses, rawFeatures, tatumCentroids, err := readFiles(analysisFileName, featureFileName, tatumCentroidFileName)
// 	if err != nil {
// 		return nil, err
// 	}

// 	analyses := formatAnalyses(rawAnalyses, tatumCentroids)
// 	features := formatFeatures(rawFeatures)

// 	if shouldWrite {
// 		err = writeCentroidsFile("tatumCentroids.txt", tatumCentroids)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	genreData = formatData(analyses, features)
// 	return genreData, nil
// }

// func formatData(analyses [][]float64, features [][]float64) (data [][]float64) {
// 	for index, val := range analyses {
// 		var datum []float64
// 		datum = append(datum, val...)
// 		datum = append(datum, features[index]...)

// 		data = append(data, datum)
// 	}

// 	return data
// }

// func getFeatures(client *spotify.Client, ids []spotify.ID) (features []*spotify.AudioFeatures, err error) {
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

// // original
// // func standardizeData(rawData [][]float64) [][]float64 {
// // 	colNum := len(rawData[0])
// // 	rowNum := len(rawData)
// // 	normData := initializeSlice(colNum, rowNum)

// // 	for i := 0; i < rowNum; i++ {
// // 		total := 0.0
// // 		for j := 0; j < colNum; j++ {
// // 			total += rawData[j][i]
// // 		}
// // 		mean := total / float64(rowNum)
// // 		deviation := 0.0
// // 		var differences []float64
// // 		for j := 0; j < colNum; j++ {
// // 			difference := rawData[j][i] - mean
// // 			differences = append(differences, difference)
// // 			deviation += math.Pow(difference, 2)
// // 		}
// // 		stdDeviation := math.Sqrt(deviation)
// // 		for j := 0; j < colNum; j++ {
// // 			normData[j][i] = (differences[j]) / stdDeviation
// // 			//normData[j][i] = (rawData[j][i] - mean) / stdDeviation // gaussian normalization equation
// // 		}
// // 	}

// // 	return normData
// // }

// // standardizeData will transform the raw data provided into a standardized set of float64 values
// // using gaussian normalization
// // row/cols in loops are wrong; see goplayground code
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
// 			//normData[j][i] = (rawData[j][i] - mean) / stdDeviation // gaussian normalization equation
// 		}
// 	}

// 	return normData
// }

// /*
// i := [][]int{{0,1,2},{3,4,5}}
// 	fmt.Println(i)
// 	fmt.Println()
// 	for _, val := range i {
// 		fmt.Println(val)
// 		}
// 	//fmt.Println()
// 	//fmt.Println(i[0])
// 	//colNum := len(i[0])
// 	//fmt.Println(colNum)
// 	fmt.Println()
// 	fmt.Println(i[0][1])
// */

// // normalizeValue performs gaussian normalization on the specified value with the provided
// // mean and standard deviation
// func normalizeValue(val, mean, dev float64) float64 {
// 	return (val - mean) / dev
// }

// func initializeSlice(width int, height int) [][]float64 {
// 	slice := make([][]float64, height)
// 	for index := range slice {
// 		slice[index] = make([]float64, width)
// 	}
// 	return slice
// }

// func formatAnalyses(rawAnalyses []*spotify.AudioAnalysis, tatumCentroids [][]float64) (analyses [][]float64) {
// 	for index, val := range rawAnalyses {
// 		var analysis []float64
// 		analysis = append(analysis, val.TrackInfo.Duration)
// 		analysis = append(analysis, val.TrackInfo.Tempo)
// 		//analysis = append(analysis, float64(val.TrackInfo.TimeSignature)) // Might remove
// 		//analysis = append(analysis, float64(val.TrackInfo.Key))           //
// 		analysis = append(analysis, val.TrackInfo.Loudness)
// 		//analysis = append(analysis, float64(val.TrackInfo.Mode)) //
// 		analysis = append(analysis, tatumCentroids[index]...)

// 		analyses = append(analyses, analysis)
// 	}
// 	return analyses
// }

// func formatFeatures(tracks []*spotify.AudioFeatures) (features [][]float64) {
// 	for _, val := range tracks {
// 		var track []float64
// 		track = append(track, float64(val.Acousticness))
// 		track = append(track, float64(val.Danceability))
// 		track = append(track, float64(val.Energy))
// 		track = append(track, float64(val.Instrumentalness)) // Might remove
// 		track = append(track, float64(val.Liveness))
// 		track = append(track, float64(val.Speechiness))
// 		track = append(track, float64(val.Valence))
// 		features = append(features, track)
// 	}
// 	return features
// }

// func getAnalyses(client *spotify.Client, ids []spotify.ID) (analyses []*spotify.AudioAnalysis, err error) {
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

// func getIDs(recs []*spotify.Recommendations) (ids []spotify.ID) {
// 	for _, val := range recs {
// 		tracks := val.Tracks
// 		for _, track := range tracks {
// 			ids = append(ids, track.ID)
// 		}
// 	}
// 	return ids
// }

// func generateRecommendations(client *spotify.Client, seeds []spotify.Seeds, limit int) (recs []*spotify.Recommendations, err error) {
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

// // I/O Functions

// func initializeArgMap(argList []string) map[string]bool {
// 	argMap := make(map[string]bool, 3)
// 	for _, val := range argList {
// 		argMap[val] = false
// 	}
// 	return argMap
// }

// func resolveCommandArgs(argMap map[string]bool) {
// 	for _, val := range os.Args {
// 		for index := range argMap {
// 			if strings.EqualFold(index, val) {
// 				argMap[index] = true
// 			}
// 		}
// 	}
// }

// func writeFiles(analysisFileName, featureFileName, tatumCentroidFileName string, analyses []*spotify.AudioAnalysis, features []*spotify.AudioFeatures, tatumCentroids [][]float64) (err error) {
// 	err = writeAnalysesFile(analysisFileName, analyses)
// 	if err != nil {
// 		return err
// 	}

// 	err = writeFeaturesFile(featureFileName, features)
// 	if err != nil {
// 		return err
// 	}

// 	err = writeCentroidsFile(tatumCentroidFileName, tatumCentroids)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func readFiles(analysisFileName string, featureFileName string, tatumCentroidFileName string) (analyses []*spotify.AudioAnalysis, features []*spotify.AudioFeatures, tatumCentroids [][]float64, err error) {
// 	analyses, err = readAnalysisFile(analysisFileName)
// 	if err != nil {
// 		return nil, nil, nil, err
// 	}
// 	features, err = readFeaturesFile(featureFileName)
// 	if err != nil {
// 		return nil, nil, nil, err
// 	}
// 	tatumCentroids, err = readCentroidsFile(tatumCentroidFileName)
// 	if err != nil {
// 		return nil, nil, nil, err
// 	}

// 	return analyses, features, tatumCentroids, nil
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

// func readCentroidsFile(fileName string) (centroids [][]float64, err error) {
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

// func writeCentroidsFile(fileName string, centroids [][]float64) error {
// 	file, err := os.Create(fileName)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	encoder := gob.NewEncoder(file)
// 	err = encoder.Encode(centroids)
// 	if err != nil {
// 		return nil
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
// 	fmt.Fprintf(w, "log Completed!")
// 	ch <- &client
// }
