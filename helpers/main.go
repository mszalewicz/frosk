// Implementation of helper functions responsible for common tasks, not included in standard library
// e.x. parsing time into ISO 8601 string
package helpers

import (
	"math/rand"
	"time"
)

// Cast date into string comforming to ISO 8601.
func TimeTo8601String(timeToParse time.Time) string {
	return timeToParse.Format("2006-01-02 15:04:05")
}

// Returns random string of a given length.
// You can choose whether to use special characters - disable option if given website for some reason does not allow them.
func RandString(length int, allowSpecialChars bool) string {
	const chars = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"0123456789"

	const specialChars = "!#$%&'()*+,-./:;<=>?@"

	charsPool := []byte(chars)

	if allowSpecialChars {
		charsPool = append(charsPool, []byte(specialChars)...)
	}

	randString := make([]byte, length)

	for i := range length {
		randString[i] = charsPool[rand.Int63()%int64(len(charsPool))]
	}

	return string(randString)
}
