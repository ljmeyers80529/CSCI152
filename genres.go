package csci152

import (
	"fmt"

	spotify "github.com/ljmeyers80529/spot-go-gae"
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

// generateUserGenreStatistics uses a spotify client, authorized within the "user-top-read" scope, to generate
// a list of the user's top 'numberOfGenres' genres and their respective scores within the given
// time range denoted by timeRange. Note: legal timeRange values are as follows - "short_term", "medium_term",
// and "long_term", strecthing from 6 weeks, to 6 months, and over several years, respectively.
func generateUserGenreStatistics(client *spotify.Client, numberOfGenres int, timeRange string) (topGenreTitles []string, topGenreScores []int, err error) {
	// Gather user's top artists
	topArtists, err := getUserTopArtists(client, timeRange)
	if err != nil {
		return nil, nil, err
	}

	genres := extractGenres(topArtists)
	for _, val := range genres {
		fmt.Println(val)
	}

	topGenreTitles, topGenreScores = calculateTopGenres(numberOfGenres, genres)

	fmt.Println("Top Genre titles: ", topGenreTitles)
	fmt.Println("Top Genre scores: ", topGenreScores)

	return topGenreTitles, topGenreScores, nil
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
