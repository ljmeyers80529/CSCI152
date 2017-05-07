// +build ignore
package main

// Trying to run this code will result in an error due to a custom function on my
// local repo of the Spotify API

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	spotify "github.com/ljmeyers80529/spot-go-gae"
)

// redirectURI is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.
const (
	redirectURI = "http://localhost:8080/callback"
	testID      = "80c614680ee64001a9fe3f5d98880364"
	testSecret  = "a3790222803a4f8fbdd5cdd5a2ce64d9"
	root        = "https://api.spotify.com/v1/"
	authroot    = "https://accounts.spotify.com/authorize"
)

var (
	auth = spotify.NewAuthenticator(redirectURI, "user-top-read")
	//auth  = spotify.NewAuthenticator(redirectURI, "user-read-recently-played")
	ch    = make(chan *spotify.Client)
	state = "abc123"
)

// Genre contains the title, score, bonus, and a list of artists pertaining to a particular genre.
// Note: none of these fields are exported (private, not public)
type Genre struct {
	title   string   // Title of the genre
	score   int      // Current score of the genre
	bonus   int      // Additional bonus points
	artists []string // List of artists within this genre
}

// construct acts as a constructor for a Genre object, setting the title and first artist
// as the two strings passed in, respectively.
func (g *Genre) construct(name, artist string) {
	g.title = name
	g.score = 1
	g.artists = []string{artist}
}

// setBonus simply sets the bonus points for the Genre object
func (g *Genre) setBonus(bonus int) {
	g.bonus = bonus
}

// addArtist appends a new artist string to the end of Genre object's string slice and increments
// its score value.
func (g *Genre) addArtist(artist string) {
	g.artists = append(g.artists, artist)
	g.score++
}

// removeArtist removes an artist from Genre object's string slice if it exists, then
// decrements its score. If the last artist is removed from the list, the score and bonus
// are reset to 0.
func (g *Genre) removeArtist(artist string) {
	for index, val := range g.artists {
		if val == artist {
			g.artists = append(g.artists[:index], g.artists[index+1:]...)
			g.score--
			if len(g.artists) == 0 {
				g.score = 0
				g.bonus = 0
			}
			break
		}
	}
}

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

	/*
		data, err := client.GetAudioFeatures("1zHlj4dQ8ZAtrayhuDDmkY")
		data, err := client.CurrentUserRecentTracks(50)
		data, err := client.GetAudioAnalysis("1zHlj4dQ8ZAtrayhuDDmkY")
		data, err := client.CurrentUserTopTracks(50, "short")
		data, err := client.CurrentUserTopArtists(50, "long")
	*/

	generateUserGenreStatistics(client, 7, "short_term")

	seeds := generateSeedsFromArtists()
	fmt.Println("seeds", seeds)
	attr := spotify.NewTrackAttributes()

	data, err := client.GetRecommendations(seeds, attr, nil)

	// for _, val := range data {
	// 	for _, item := range val.Items {
	// 		//do something to each track
	// 	}
	// }

	d, err := json.MarshalIndent(data, "", "  ")
	fmt.Println("JSON DATA")
	fmt.Println(d)
	for _, val := range data.Tracks {
		fmt.Println(val.Name)
	}
	err = ioutil.WriteFile("temp.txt", d, 0644)

	if err != nil {
		fmt.Println(err)
	}

}

func getSeeds(genres []string) (seeds []spotify.Seeds) {
	for _, val := range genres {
		var values []string
		values = append(values, val)
		newSeed := spotify.Seeds{Genres: values}
		seeds = append(seeds, newSeed)
	}
	seeds := spotify.Seeds{Tracks: trackSeeds}

	return seeds
}

// generateUserGenreStatistics uses a spotify client, authorized within the "user-top-read" scope, to generate
// a list of the user's top 'numberOfGenres' genres and their respective scores within the given
// time range denoted by timeRange. Note: legal timeRange values are as follows - "short_term", "medium_term",
// and "long_term", strecthing from 6 weeks, to 6 months, and over several years, respectively.
func generateUserGenreStatistics(client *spotify.Client, numberOfGenres int, timeRange string) (topGenreTitles []string, topGenreScores []int, topArtists *spotify.TopArtists, err error) {
	// Gather user's top artists
	topArtists, err = getUserTopArtists(client, timeRange)
	if err != nil {
		return nil, nil, nil, err
	}

	genres := extractGenres(topArtists)
	for _, val := range genres {
		fmt.Println(val)
	}

	topGenreTitles, topGenreScores = calculateTopGenres(numberOfGenres, genres)

	fmt.Println("Top Genre titles: ", topGenreTitles)
	fmt.Println("Top Genre scores: ", topGenreScores)

	return topGenreTitles, topGenreScores, topArtists, nil
}

// getUserTopArtists uses a spotify client, authorized within the "user-top-read" scope, to
// get the users top 50 artists within the given time range using Spotify endpoints.
func getUserTopArtists(client *spotify.Client, timeRange string) (top *spotify.TopArtists, err error) {
	limit := 50

	opt := spotify.Options{
		Limit:     &limit,
		Timerange: &timeRange,
	}
	top, err = client.CurrentUserTopArtists(&opt)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return top, nil
}

// extractGenres parses through a list of artists stored within a TopArtists object in order
// to return a list of Genre objects with their respective titles, artists, scores, and bonuses
// set to their correct values.
func extractGenres(artists *spotify.TopArtists) (genreList []Genre) {
	bonus := 50
	for _, item := range artists.Items {
		for _, val := range item.Genres {
			if index, ok := genreExists(val, genreList); ok {
				genreList[index].addArtist(item.Name)
				genreList[index].bonus += bonus / 10
			} else {
				var temp Genre
				temp.construct(val, item.Name)
				temp.bonus = bonus / 10
				genreList = append(genreList, temp)
			}
		}
		bonus--
	}
	return genreList
}

// genreExists is a helper function that checks for the given genre title within the given
// list of Genres and returns a boolean flag representing its existence along with the element's index.
func genreExists(genre string, list []Genre) (int, bool) {
	for index, val := range list {
		if genre == val.title {
			return index, true
		}
	}
	return 0, false
}

// calculateTopGenres takes a list of Genres and the desired limit of output, and returns an
// ordered list containing the title of each genre and a separate list of ints containing
// their respective final scores.
func calculateTopGenres(limit int, genres []Genre) (titles []string, scores []int) {
	for limit > 0 {
		topIndex := findTopGenreIndex(genres)
		titles = append(titles, genres[topIndex].title)
		scores = append(scores, genres[topIndex].score+genres[topIndex].bonus)
		recalculateGenreScores(topIndex, genres)
		limit--
	}
	return titles, scores
}

// findTopGenreIndex thoroughly parses through the given list of Genres and returns the
// index of the Genre encountered with the highest total score.
func findTopGenreIndex(genres []Genre) int {
	max := 0
	index := 0
	for i, val := range genres {
		if val.score+val.bonus > max {
			max = val.score + val.bonus
			index = i
		}
	}
	return index
}

// recalculateGenreScores parses through the given list of Genres, deleting the the artists encountered
// in the Genre denoted by the given index from every other Genre object in the list, thus recalculating
// the scores for every genre affected.
func recalculateGenreScores(index int, genres []Genre) {
	artists := make([]string, len(genres[index].artists))
	copy(artists, genres[index].artists)
	for i := range genres {
		for _, artist := range artists {
			genres[i].removeArtist(artist)
		}
	}
}

// func getUserTop(client *spotify.Client, timeRange string) (top *spotify.TopTracks, err error) {
// 	limit := 50

// 	opt := spotify.Options{
// 		Limit:     &limit,
// 		Timerange: &timeRange,
// 	}
// 	top, err = client.CurrentUserTopTracks(&opt)
// 	if err != nil {
// 		fmt.Println(err)
// 		return nil, err
// 	}
// 	return top, nil
// }

// func getArtistIDsFromTop(topTracks *spotify.TopTracks) (artistIDs []spotify.ID) {
// 	for _, item := range topTracks.Items {
// 		artistIDs = append(artistIDs, item.Artists[0].ID)
// 	}

// 	return artistIDs
// }

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

////////////////////////////////////////////////////////////

// // +build ignore
// package main

// // Trying to run this code will result in an error due to a custom function on my
// // local repo of the Spotify API

// import (
// 	"fmt"
// 	"log"
// 	"net/http"

// 	spotify "github.com/ljmeyers80529/spot-go-gae"
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

// // Genre contains the title, score and a list of artists pertaining to a particular genre.
// // Note: none of these fields are exported (private, not public)
// type Genre struct {
// 	title   string
// 	score   int
// 	artists []string
// }

// // construct acts as a constructor for a Genre object, setting the title and first artist
// // as the two strings passed in, respectively.
// func (g *Genre) construct(name, artist string) {
// 	g.title = name
// 	g.score = 1
// 	g.artists = []string{artist}
// }

// // addArtist appends a new artist string to the end of Genre object's string slice and increments
// // its score value.
// func (g *Genre) addArtist(artist string) {
// 	g.artists = append(g.artists, artist)
// 	g.score++
// }

// // removeArtist removes an artist from Genre object's string slice, decrementing its score value as well.
// // func (g *Genre) removeArtist(artist string) {
// // 	var s []string
// // 	for _, val := range g.artists {
// // 		if val == artist {
// // 			continue
// // 		}
// // 		s = append(s, val)
// // 	}
// // 	g.artists = s
// // 	g.score--
// // }

// // func (g *Genre) removeArtist(artist string) {
// // 	for index, val := range g.artists {
// // 		if val == artist {
// // 			if len(g.artists) <= 1 {
// // 				g.artists = make([]string, 0)
// // 				g.score--
// // 				return
// // 			}
// // 			g.artists[index] = g.artists[len(g.artists)-1]
// // 			g.artists = g.artists[:len(g.artists)-1]
// // 			g.score--
// // 			break
// // 		}
// // 	}
// // }

// func (g *Genre) removeArtist(artist string) {
// 	for index, val := range g.artists {
// 		if val == artist {
// 			g.artists = append(g.artists[:index], g.artists[index+1:]...)
// 			g.score--
// 			break
// 		}
// 	}
// }

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

// 	/*
// 		data, err := client.GetAudioFeatures("1zHlj4dQ8ZAtrayhuDDmkY")
// 		data, err := client.CurrentUserRecentTracks(50)
// 		data, err := client.GetAudioAnalysis("1zHlj4dQ8ZAtrayhuDDmkY")
// 		data, err := client.CurrentUserTopTracks(50, "short")
// 		data, err := client.CurrentUserTopArtists(50, "long")
// 	*/

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
// 	fmt.Println("Getting personalized playlist...")
// 	// Gather user's top artists
// 	topArtists, err := getUserTopArtists(client, "short_term")
// 	if err != nil {
// 		return nil, err
// 	}

// 	genres := extractGenres(topArtists)
// 	for _, val := range genres {
// 		fmt.Println(val)
// 	}

// 	topGenreTitles, topGenreScores := calculateTopGenres(7, genres)

// 	fmt.Println("Top Genre titles: ", topGenreTitles)
// 	fmt.Println("Top Genre scores: ", topGenreScores)

// 	return nil, nil
// }

// func calculateTopGenres(limit int, genres []Genre) (titles []string, scores []int) {
// 	fmt.Println("Calculating top genres...")
// 	for limit > 0 {
// 		topIndex := findTopGenreIndex(genres)
// 		fmt.Println("Top Genre: ", genres[topIndex])
// 		titles = append(titles, genres[topIndex].title)
// 		scores = append(scores, genres[topIndex].score)
// 		recalculateGenreScores(topIndex, genres)
// 		limit--
// 	}
// 	return titles, scores
// }

// func recalculateGenreScores(index int, genres []Genre) {
// 	fmt.Println("Recalculating scores...")
// 	artists := make([]string, len(genres[index].artists))
// 	copy(artists, genres[index].artists)
// 	for i := range genres {
// 		//for _, artist := range genres[index].artists {
// 		for _, artist := range artists {
// 			genres[i].removeArtist(artist)
// 		}
// 	}
// }

// func findTopGenreIndex(genres []Genre) int {
// 	max := 0
// 	index := 0
// 	for i, val := range genres {
// 		if val.score > max {
// 			max = val.score
// 			index = i
// 		}
// 	}
// 	return index
// }

// /*
// func getPersonalizedPlaylist(client *spotify.Client) (*spotify.SimplePlaylist, error) {
// 	// Gather user's top tracks

// 	topTracks, err := getUserTop(client, "short_term")
// 	if err != nil {
// 		return nil, err
// 	}

// 	topTrackArtistIDs := getArtistIDsFromTop(topTracks)

// 	artistObjects, err := spotify.GetArtists(topTrackArtistIDs...)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// genres := extractGenres(artistObjects)

// 	genreScores := calculateGenreCount(artistObjects)

// 	//artistsPerGenre := sortArtistsByGenre(artistObjects)

// 	// Map of the form map[ArtistID][]TrackID
// 	//tracksByArtist := make(map[spotify.ID][]spotify.ID)
// 	tracksByName := make(map[string][]string)

// 	// for _, val := range shortTop.Items {
// 	// 	tracksByArtist[val.Artists[0].ID] = append(tracksByArtist[val.Artists[0].ID], val.ID)
// 	// 	tracksByName[val.Artists[0].Name] = append(tracksByName[val.Artists[0].Name], val.Name)
// 	// 	artistIDList = append(artistIDList, val.Artists[0].ID)
// 	// }

// 	for ind, val := range tracksByName {
// 		fmt.Println("Artist:", ind, "\nTrack:", val, "\n")
// 	}

// 	// Sort genres by occurrence
// 	genreBySortedCounts := make(map[int][]string)
// 	var temp []int

// 	for k, v := range genreScores {
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
// */

// func getUserTopArtists(client *spotify.Client, timeRange string) (top *spotify.TopArtists, err error) {
// 	fmt.Println("Getting top artists...")
// 	limit := 50

// 	opt := spotify.Options{
// 		Limit:     &limit,
// 		Timerange: &timeRange,
// 	}
// 	top, err = client.CurrentUserTopArtists(&opt)
// 	if err != nil {
// 		fmt.Println(err)
// 		return nil, err
// 	}

// 	return top, nil
// }

// func extractGenres(artists *spotify.TopArtists) (genreList []Genre) {
// 	fmt.Println("Extracting genres...")
// 	for _, item := range artists.Items {
// 		for _, val := range item.Genres {
// 			if index, ok := genreExists(val, genreList); ok {
// 				genreList[index].addArtist(item.Name)
// 			} else {
// 				var temp Genre
// 				temp.construct(val, item.Name)
// 				genreList = append(genreList, temp)
// 			}
// 		}
// 	}
// 	return genreList
// }

// func genreExists(genre string, list []Genre) (int, bool) {
// 	for index, val := range list {
// 		if genre == val.title {
// 			return index, true
// 		}
// 	}
// 	return 0, false
// }

// /*
// func extractGenres(artists *spotify.TopArtists) map[string]Genre {
// 	genreMap := make(map[string]Genre)

// 	for _, item := range artists.Items {
// 		for _, val := range item.Genres {
// 			fmt.Println(val)
// 			if _, ok := genreMap[val]; ok {
// 				genreMap[val].addArtist(item.Name)
// 			} else {
// 				var temp Genre
// 				temp.construct(val, item.Name)
// 				genreMap[val] = temp
// 			}
// 		}
// 	}
// 	return genreMap
// }
// */

// /*
// func calculateGenreCount(artistList []*spotify.FullArtist) map[string]int {
// 	genreScores := make(map[string]int)
// 	for _, artist := range artistList {
// 		for _, genre := range artist.Genres {
// 			if _, ok := genreScores[genre]; ok {
// 				genreScores[genre] = genreScores[genre] + 1
// 			} else {
// 				genreScores[genre] = 1
// 			}
// 		}
// 	}
// 	return genreScores
// }
// */

// func getUserTop(client *spotify.Client, timeRange string) (top *spotify.TopTracks, err error) {
// 	limit := 50

// 	opt := spotify.Options{
// 		Limit:     &limit,
// 		Timerange: &timeRange,
// 	}
// 	top, err = client.CurrentUserTopTracks(&opt)
// 	if err != nil {
// 		fmt.Println(err)
// 		return nil, err
// 	}
// 	return top, nil
// }

// func getArtistIDsFromTop(topTracks *spotify.TopTracks) (artistIDs []spotify.ID) {
// 	for _, item := range topTracks.Items {
// 		artistIDs = append(artistIDs, item.Artists[0].ID)
// 	}

// 	return artistIDs
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
