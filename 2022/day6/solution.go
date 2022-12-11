package main

import (
	"container/ring"
	"errors"
	"os"
)

func solveTuningTrouble(markerLen int) (int, error) {
	// Just read all file, no bufferring.
	data, err := os.ReadFile("input.txt")
	if err != nil {
		return 0, err
	}

	circleBuf := ring.New(markerLen)
	for _, ch := range data[:markerLen] {
		circleBuf.Value = ch
		circleBuf = circleBuf.Next()
	}

	uniqueChars := make(map[byte]bool, markerLen)
	for i := markerLen; ; i++ {
		circleBuf.Do(func(a any) {
			uniqueChars[a.(byte)] = true
		})
		if len(uniqueChars) == markerLen {
			return i, nil
		}
		for k := range uniqueChars {
			delete(uniqueChars, k)
		}

		if i > len(data)-1 {
			return 0, errors.New("marker not found")
		}

		circleBuf.Value = data[i]
		circleBuf = circleBuf.Next()
	}
}
