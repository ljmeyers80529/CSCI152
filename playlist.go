package csci152

import (
	"errors"
	"math"

	spotify "github.com/ljmeyers80529/spot-go-gae"
)

// generatePersonalizedPlaylist uses uses an authenticated spotify client and the user's genre
// statistics as input to generate a playlist of a size denoted by the playlistSize parameter.
// The required statistics input are two lists and an object that can be retrieved using the
// generateUserGenreStatistics function and are as follows: a list of the users top genres as strings,
// a list of the user's genre's top scores as ints stored respective to the previous top genres list,
// and a topArtists object for the user. The output of the function is a complete FullPlaylist object
// containing all of the identifying information needed such as ID, URI, Name, Owner, etc.
func generateUserPlaylist(client *spotify.Client, playlistSize int, topGenres []string, topScores []int, topArtists *spotify.TopArtists) (playlist *spotify.FullPlaylist, err error) {
	if playlistSize > 50 {
		err = errors.New("generateUserPlaylist: queried playlist creation size exceeds 50")
		return nil, err
	}
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
