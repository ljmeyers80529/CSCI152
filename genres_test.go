package csci152

import (
	"log"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShuffleListsInParallel(t *testing.T) {

	strings := []string{"one", "two", "three", "four"}
	ints := []int{1, 2, 3, 4}

	shuffleListsInParallel(strings, ints)

	for index, val := range ints {
		s := strings[index]
		stringAsInt, err := strconv.Atoi(s)
		if err != nil {
			log.Println(err)
			return
		}
		assert.Equal(t, val, stringAsInt, "elements out of respective order")
	}
}
