package csci152

import (
	"errors"
	"math"

	spotify "github.com/ljmeyers80529/spot-go-gae"
)

// generateUserPlaylist uses an authenticated spotify client to generate a personalized playlist
// to the currently logged in user, leveraging their top played artists and their most listened to
// genres as seeds for the playlist creation. The goal of the playlist is to demonstrate the user's
// taste in a music in a compact manner. The returned value a FullPlaylist object containing all of
// the identifying information needed such as ID, URI, Name, Owner, etc.
func generateUserPlaylist(client *spotify.Client) (playlist *spotify.FullPlaylist, err error) {
	playlistSize := 30

	topGenreTitles, topGenreScores, topArtists, err := generateUserGenreStatistics(client, 7, "short_term")
	if err != nil {
		return nil, err
	}

	seeds, err := generateSeedsByGenre(topGenreTitles, topArtists.Items)
	if err != nil {
		return nil, err
	}

	genreWeights := calculateGenreWeights(topGenreScores)
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
	attr := spotify.NewTrackAttributes().TargetPopularity(80)

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
