package utils

import (
	"fmt"
	"strings"

	"math/rand"
)

// Functions

// GenerateString returns a random string from the
// alphabet [a-z,0-9] of length "strlen".
func GenerateString(strlen int) string {

	// Define alphabet.
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"

	result := ""
	for i := 0; i < strlen; i++ {
		index := rand.Intn(len(chars))
		result += chars[index:(index + 1)]
	}

	return result
}

// GenerateFlag returns a random choice of message flags.
func GenerateFlags() (string, []string) {

	// Define alphabet.
	flags := []string{"\\Seen", "\\Answered", "\\Flagged", "\\Deleted", "\\Draft"}

	numFlags := rand.Intn(len(flags)) + 1

	// Generate an array of random but different indices.
	var genIndex []int
	for len(genIndex) < numFlags {

		index := rand.Intn(len(flags))

		for i := 0; i < len(genIndex); i++ {

			if index == genIndex[i] {
				index = rand.Intn(len(flags))
				i = -1
			}
		}

		genIndex = append(genIndex, index)
	}

	// Add the corresponding flag of the previously generated
	// index to the string array "genFlags".
	var genFlags []string
	for i := 0; i < len(genIndex); i++ {
		genFlags = append(genFlags, flags[genIndex[i]])
	}

	// Generate final flag string.
	flagString := fmt.Sprintf("(%s)", strings.Join(genFlags, " "))

	return flagString, genFlags
}
