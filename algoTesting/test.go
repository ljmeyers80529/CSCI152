package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"

	spotify "github.com/ljmeyers80529/spot-go-gae"
)

// redirectURI is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.
const (
	redirectURI       = "http://localhost:8080/callback"
	testID            = "80c614680ee64001a9fe3f5d98880364"
	testSecret        = "a3790222803a4f8fbdd5cdd5a2ce64d9"
	root              = "https://api.spotify.com/v1/"
	authroot          = "https://accounts.spotify.com/authorize"
	playlistNameConst = "Taste Test - Personal Playlist"
)

var (
	auth = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserTopRead, spotify.ScopePlaylistModifyPublic)
	//auth  = spotify.NewAuthenticator(redirectURI, "user-read-recently-played")
	ch    = make(chan *spotify.Client)
	ch2   = make(chan *spotify.Client)
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

	// topGenreTitle, topGenreScores, topArtists, err := generateUserGenreStatistics(client, 7, "short_term")
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	topGenres, topScores, topArtists, err := generateUserGenreStatistics(client, 3, "long_term")
	if err != nil {
		log.Println(err)
		return
	}

	_, err = generateUserPlaylist(client, 30, topGenres, topScores, topArtists)
	if err != nil {
		log.Println(err)
		return
	}

	client2 := <-ch2
	fmt.Println(client2)

	// d, err := json.MarshalIndent(data, "", "  ")
	// fmt.Println("JSON DATA")
	// fmt.Println(d)

	// err = ioutil.WriteFile("temp.txt", d, 0644)
	// if err != nil {
	// 	fmt.Println(err)
	// }

}

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

// generateUserGenreStatistics uses a spotify client, authorized within the "user-top-read" scope, to generate
// a list of the user's top 'numberOfGenres' (limited to 10) genres and their respective scores within the given
// time range denoted by timeRange. Additionally, the user's top artists are returned as a TopArtists object
// Note: legal timeRange values are as follows - "short_term", "medium_term", and "long_term", stretching from
// 6 weeks, to 6 months, and over several years, respectively.
func generateUserGenreStatistics(client *spotify.Client, numberOfGenres int, timeRange string) (topGenreTitles []string, topGenreScores []int, topArtists *spotify.TopArtists, err error) {
	// Gather user's top artists
	if numberOfGenres > 10 {
		err = errors.New("number of genre's requested exceeds 10")
		return nil, nil, nil, err
	}

	if !(timeRange == "short_term" || timeRange == "medium_term" || timeRange == "long_term") {
		err = errors.New("invalid time range input")
		return nil, nil, nil, err
	}

	topArtists, err = getUserTopArtists(client, timeRange)
	if err != nil {
		return nil, nil, nil, err
	}

	genres, err := extractGenres(topArtists, numberOfGenres)
	if err != nil {
		return nil, nil, nil, err
	}
	for _, val := range genres {
		fmt.Println(val)
	}

	topGenreTitles, topGenreScores, err = calculateTopGenres(numberOfGenres, genres)
	if err != nil {
		return nil, nil, nil, err
	}

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
func extractGenres(artists *spotify.TopArtists, genreFloorLimit int) (genreList []Genre, err error) {
	if len(artists.Items) < 15 {
		err = errors.New("extract genres: insufficient 'top artist' information in user data")
		return nil, err
	}
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
	if len(genreList) < genreFloorLimit {
		err = errors.New("extract genres: insufficient amount of genre information extracted")
		return nil, err
	}
	return genreList, nil
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

// calculateTopGenres takes a list of Genres and the desired floor limit of output, and returns an
// ordered list containing the title of each genre and a separate list of ints containing
// their respective final scores.
func calculateTopGenres(floorLimit int, genres []Genre) (titles []string, scores []int, err error) {
	if floorLimit <= 0 {
		err = errors.New("calculateTopGenres: invalid 0 or negative floor limit")
		return nil, nil, err
	}
	for floorLimit > 0 {
		if len(genres) == 0 {
			return titles, scores, nil
		}
		topIndex, err := findTopGenreIndex(genres)
		if err != nil {
			return nil, nil, err
		}
		titles = append(titles, genres[topIndex].title)
		scores = append(scores, genres[topIndex].score+genres[topIndex].bonus)
		recalculateGenreScores(topIndex, genres)
		floorLimit--
	}
	return titles, scores, nil
}

// findTopGenreIndex thoroughly parses through the given list of Genres and returns the
// index of the Genre encountered with the highest total score.
func findTopGenreIndex(genres []Genre) (int, error) {
	if len(genres) <= 0 {
		err := errors.New("findTopGenreIndex: genre list is empty")
		return 0, err
	}
	max := 0
	index := 0
	for i, val := range genres {
		if val.score+val.bonus > max {
			max = val.score + val.bonus
			index = i
		}
	}
	return index, nil
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

// generatePersonalizedPlaylist uses uses an authenticated spotify client and the user's genre
// statistics as input to generate a playlist of a size denoted by the playlistSize parameter.
// The required statistics input are two lists and an object that can be retrieved using the
// generateUserGenreStatistics function and are as follows: a list of the users top genres as strings,
// a list of the user's genre's top scores as ints stored respective to the previous top genres list,
// and a topArtists object for the user. The output of the function is a complete FullPlaylist object
// containing all of the identifying information needed such as ID, URI, Name, Owner, etc.
func generateUserPlaylist(client *spotify.Client, playlistSize int, topGenres []string, topScores []int, topArtists *spotify.TopArtists) (playlist *spotify.FullPlaylist, err error) {
	seeds, err := generateSeedsByGenre(topGenres, topArtists.Items)
	if err != nil {
		return nil, err
	}

	genreWeights := calculateGenreWeights(topScores)
	tracksPerGenre := calculateTracksPerGenre(genreWeights, playlistSize)
	recommendations, err := getSeededRecommendations(client, seeds, tracksPerGenre)
	if err != nil {
		return nil, err
	}
	playlistTrackIDs := extractIDsFromRecommendations(recommendations)

	user, err := client.CurrentUser()
	if err != nil {
		return nil, err
	}

	playlist, err = client.CreatePlaylistForUser(user.ID, playlistNameConst, true)
	if err != nil {
		return nil, err
	}
	_, err = client.AddTracksToPlaylist(user.ID, playlist.ID, playlistTrackIDs...)
	if err != nil {
		return nil, err
	}
	return playlist, nil
}

// generateSeedsByGenre takes a list of artists and their corresponding genres and returns a list
// of seeds generated using the aforementioned lists.
func generateSeedsByGenre(topGenres []string, artists []spotify.ArtistItem) (seeds []spotify.Seeds, err error) {
	for _, targetGenre := range topGenres {
		newSeed := getArtistSeedForGenre(artists, targetGenre)
		seeds = append(seeds, newSeed)
	}

	if len(seeds) <= 0 {
		err = errors.New("no seeds generated from artists")
		return nil, err
	}

	return seeds, nil
}

// getArtistSeedForGenre returns a list of seeds for the provided genre using the available
// artists in the provided artist list as its seed contents.
func getArtistSeedForGenre(artists []spotify.ArtistItem, targetGenre string) spotify.Seeds {
	maxSeedInput := 5
	newSeed := spotify.Seeds{}

	for _, artist := range artists {
		if len(newSeed.Artists) >= maxSeedInput {
			break
		}
		for _, currentGenre := range artist.Genres {
			if currentGenre == targetGenre {
				newSeed.Artists = append(newSeed.Artists, artist.ID)
			}
		}
	}
	return newSeed
}

// calculateGenreWeights takes a list of scores and returns a list of corresponding weights
// representing each scores proportion relative to the others.
func calculateGenreWeights(scores []int) (weights []float64) {
	sum := calculateSum(scores)
	for _, score := range scores {
		weight := 0.0
		weight = float64(score) / float64(sum)
		weights = append(weights, weight)
	}
	return weights
}

// calculateSum takes a list of ints and returns the sum its contents.
func calculateSum(scores []int) (sum int) {
	for _, score := range scores {
		sum += score
	}
	return sum
}

// calculateTracksPerGenre calculates the amount of tracks to be fetched for each genre
// according to the provided weights and total number of max tracks to be fetched and
// returns it as a slice with indecies corresponding to the aforementioned genres.
func calculateTracksPerGenre(weights []float64, total int) (tracksPerGenre []int) {
	for _, weight := range weights {
		decimal := weight * float64(total)
		rounded := math.Ceil(decimal)
		tracksPerGenre = append(tracksPerGenre, int(rounded))
	}
	return tracksPerGenre
}

// getSeededRecommendations uses an authorized spotify client to fetch a list of recommendation
// objects from a Spotify endpoint using the provided seeds and the corresponding track limit.
func getSeededRecommendations(client *spotify.Client, seeds []spotify.Seeds, tracksPerGenre []int) (recommendations []*spotify.Recommendations, err error) {
	attr := spotify.NewTrackAttributes()

	for index, seed := range seeds {
		opts := spotify.Options{Limit: &tracksPerGenre[index]}
		newRecommendations, err := client.GetRecommendations(seed, attr, &opts)
		if err != nil {
			return nil, err
		}
		recommendations = append(recommendations, newRecommendations)
	}
	return recommendations, nil
}

// extractIDsFromRecommendations parses through a list of recommendation objects and returns
// a list of ID's corresponding to the provided recommended tracks.
func extractIDsFromRecommendations(recommendations []*spotify.Recommendations) (trackIDs []spotify.ID) {
	for _, rec := range recommendations {
		for _, track := range rec.Tracks {
			trackIDs = append(trackIDs, track.ID)
		}
	}
	return trackIDs
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
	// Retrieve authenticated client using token
	client := auth.NewClient(tok)
	fmt.Fprintf(w, "Login Completed!")
	ch <- &client
}
