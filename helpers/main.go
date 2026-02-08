// Implementation of helper functions responsible for common tasks, not included in standard library
// e.x. parsing time into ISO 8601 string
package helpers

import (
	"fmt"
	"log/slog"
	"math/rand/v2"
	"runtime/debug"
	"time"

	"golang.org/x/exp/constraints"
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

	const specialChars = "^!#$%&'()*+,-./:;<=>?@_|\"\\"

	charsPool := []byte(chars)

	if allowSpecialChars {
		charsPool = append(charsPool, []byte(specialChars)...)
	}

	randString := make([]byte, length)

	for i := range length {
		randString[i] = charsPool[rand.Int64()%int64(len(charsPool))]
	}

	return string(randString)
}

// Asserting if two values are equall. If not, it stops execution of program and logs error to the slog logger and stdout.
func Assert[T comparable](x T, y T) {
	if x != y {
		err := fmt.Errorf("Assert error. Value x does not equal value y. Where x = %v, y = %v", x, y)
		slog.Error(err.Error())
		panic(err)
	}
}

// Asserting first value is bigger then the second one. If not, it stops execution of program and logs error to the slog logger and stdout.
func AssertBigger[T constraints.Ordered](x T, y T) {
	if x < y {
		err := fmt.Errorf("Assert error. Value x is not bigger then y. Where x = %v, y = %v", x, y)
		slog.Error(err.Error())
		panic(err)
	}
}

// Gracefuly handle panic - useful for external packages
func RunAndHandlePanic(f func() []byte) (encryptedValue []byte, err error) {
	if r := recover(); r != nil {
		encryptedValue = nil
		err = fmt.Errorf("Panic in %s: %v", string(debug.Stack()), r)
	}

	encryptedValue = f()

	return encryptedValue, nil
}
