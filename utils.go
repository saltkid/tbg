package main

import (
	"math/rand/v2"
)

// shuffles the slice from the current index up to the end
//
// it does not affect the elements before the current index
func ShuffleFrom[T any](currentIndex int, slice []T) {
	for range slice[currentIndex:] {
		i := rand.IntN(len(slice)-currentIndex) + currentIndex
		slice[i], slice[currentIndex] = slice[currentIndex], slice[i]
	}
}
