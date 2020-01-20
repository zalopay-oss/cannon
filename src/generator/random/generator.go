package random

import (
	"errors"
	"fmt"
	"math/rand"
	"encoding/base64"
)

type NumberBoundary struct {
	Start int
	End int
}

var (
	nBoundary = NumberBoundary{Start: 0, End: 100}
	testRandZero = false
	randomSize = 100
	randomStringLen = 25
)

const (
	letterIdxBits         = 6                    // 6 bits to represent a letter index
	letterIdxMask         = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax          = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	letterBytes           = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberBytes			  = "0123456789"
)

func SetRandomNumberBoundaries(start, end int) error {
	if start > end {
		return errors.New("Start value is bigger than End value")
	}
	nBoundary = NumberBoundary{Start: start, End: end}
	return nil
}

// SetRandomStringLength sets a length for random string generation
func SetRandomStringLength(size int) error {
	if size < 0 {
		return fmt.Errorf("Err Smaller Than Zero", size)
	}
	randomStringLen = size
	return nil
}

// SetRandomMapAndSliceSize sets the size for maps and slices for random generation.
func SetRandomMapAndSliceSize(size int) error {
	if size < 0 {
		return fmt.Errorf("Err Smaller Than Zero", size)
	}
	randomSize = size
	return nil
}

func RandomString() string {
	b := make([]byte, randomStringLen)
	for i, cache, remain := randomStringLen-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// randomIntegerWithBoundary returns a random integer between input start and end boundary. [start, end)
func RandomIntegerWithBoundary(boundary NumberBoundary) int {
	return rand.Intn(boundary.End-boundary.Start) + boundary.Start
}

// randomInteger returns a random integer between start and end boundary. [start, end)
func RandomInteger() int {
	return rand.Intn(nBoundary.End-nBoundary.Start) + nBoundary.Start
}

// randomSliceAndMapSize returns a random integer between [0,randomSliceAndMapSize). If the testRandZero is set, returns 0
// Written for test purposes for shouldSetNil
func RandomSliceAndMapSize() int {
	if testRandZero {
		return 0
	}
	return rand.Intn(randomSize)
}

func RandomElementFromSliceString(s []string) string {
	return s[rand.Int()%len(s)]
}
func RandomStringNumber(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(numberBytes) {
			b[i] = numberBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// RandomInt Get three parameters , only first mandatory and the rest are optional
// 		If only set one parameter :  This means the minimum number of digits and the total number
// 		If only set two parameters : First this is min digit and second max digit and the total number the difference between them
// 		If only three parameters: the third argument set Max count Digit
func RandomInt(parameters ...int) (p []int, err error) {
	switch len(parameters) {
	case 1:
		minCount := parameters[0]
		p = rand.Perm(minCount)
		for i := range p {
			p[i] += minCount
		}
	case 2:
		minDigit, maxDigit := parameters[0], parameters[1]
		p = rand.Perm(maxDigit - minDigit + 1)

		for i := range p {
			p[i] += minDigit
		}
	default:
		err = fmt.Errorf("Error more arguments", len(parameters))
	}
	return p, err
}

func RandomBytesData() string {
	data := RandomString()
	b64Data := base64.StdEncoding.EncodeToString([]byte(data))
	return b64Data
}
