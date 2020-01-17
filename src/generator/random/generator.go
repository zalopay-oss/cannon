package random

import (
	"errors"
	"fmt"
	"math/rand"
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
