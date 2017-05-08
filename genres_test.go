package csci152

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenreConstruct(t *testing.T) {
	var g Genre
	g.construct("Rock", "Queen")

	assert.Equal(t, "Rock", g.title)
	assert.Contains(t, g.artists, "Queen")
	assert.Equal(t, 1, g.score)
	assert.Equal(t, 0, g.bonus, "bonus not initialized to 0")
}

func TestSetBonus(t *testing.T) {
	var g Genre
	assert.Equal(t, 0, g.bonus, "bonus not initialized to 0")

	g.setBonus(10)
	assert.Equal(t, 10, g.bonus)

	g.setBonus(-99)
	assert.Equal(t, -99, g.bonus)

	g.setBonus(0)
	assert.Equal(t, 0, g.bonus)

	assert.Empty(t, g.artists, "setBonus interfering with artist list")
	assert.Equal(t, 0, g.score, "setBonus interfering with score")
	assert.Equal(t, "", g.title, "setBonus interfering with title")
}

func TestAddArtist(t *testing.T) {
	var g Genre
	testCases := []string{"Queen", "Zion.T", "The Smiths", "DREAMCAR", "King Gizzard and the Lizard Wizard", "t e l e p a t h テレパシー能力者", "2 8 1 4"}
	for _, testCase := range testCases {
		g.addArtist(testCase)
		t.Log("Current state: ", g.artists)
		assert.Equal(t, len(g.artists), g.score, "score not keeping up with artist list manipulation")
	}
	for _, testCase := range testCases {
		assert.Contains(t, g.artists, testCase)
	}

	assert.Equal(t, 0, g.bonus, "addArtist interfering with bonus")
	assert.Equal(t, "", g.title, "addArtist interfering with title")
}

func TestRemoveArtist(t *testing.T) {
	var g Genre
	testCases := []string{"Queen", "Zion.T", "The Smiths", "DREAMCAR", "King Gizzard and the Lizard Wizard", "t e l e p a t h テレパシー能力者", "2 8 1 4"}
	for _, testCase := range testCases {
		g.addArtist(testCase)
		t.Log("Current state: ", g.artists)
	}

	newCases := []string{"King Gizzard and the Lizard Wizard", "t e l e p a t h テレパシー能力者", "2 8 1 4", "DREAMCAR", "Queen", "Zion.T", "The Smiths"}
	for _, newCase := range newCases {
		g.removeArtist(newCase)
		t.Log("Current state: ", g.artists)
		assert.NotContains(t, g.artists, newCase)
		assert.Equal(t, len(g.artists), g.score, "score not keeping up with artist list manipulation")
	}

	assert.Empty(t, g.artists)
	assert.Equal(t, 0, g.bonus, "removeArtist interfering with bonus")
	assert.Equal(t, "", g.title, "removeArtist interfering with title")
}
