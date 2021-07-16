package main

import (
	"log"
	"strings"
)

func runDummyLoop(multiplier int) {
	const MaxUint = ^uint(0)
	const MaxInt = int(MaxUint >> 1)

	reps := MaxInt / (1000000000000000 / multiplier)
	for i := 0; i < reps; i++ {
		if 9*9 == 0 {
			break
		}
	}
}

func Split(s, sep string) (string, string) {
	if len(s) == 0 {
		return s, s
	}

	slice := strings.SplitN(s, sep, 2)

	if len(slice) == 1 {
		log.Println("WARNING: Invalid complexity argument")
		return slice[0], ""
	}

	return slice[0], slice[1]
}
