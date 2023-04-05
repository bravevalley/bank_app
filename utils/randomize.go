package utils

import (
	"math/rand"
	"strings"
	"time"
)

var rdm *rand.Rand
var alphabets = "abcdefghijklmnopqstuvxyz"

func init() {
	// Seed the random number generator
	rdm = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// rdmNumbBtwRange return a random number within the range
func RdmNumbBtwRange(min, max int64) int64 {
	return min + rdm.Int63n(max-min+1)
}

// rdmstring returns a random collection of strings
func rdmString(length int) string {
	var sb strings.Builder

	lenAlp := len(alphabets)

	for i := 0; i < length; i++ {
		rndAlp := alphabets[rdm.Intn(lenAlp)]

		sb.WriteByte(rndAlp)
	}

	return sb.String()
}

// RandowmName returns a name from a random collection of strings
func RandomName() string {
	return rdmString(7)
}

// RandomAmount returns a random amount of money
func RandomAmount() int64 {
	return RdmNumbBtwRange(1, 1000)
}

// RdmCurr returns a ramdom currency
func RdnCurr() string {
	xCurr := []string{"NGN", "USD", "YEN"}

	return xCurr[rdm.Intn(len(xCurr))]
}
