package main

import (
	"log"
	"strconv"
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

func GetMultiplier(complexityConfig string) int {
	if complexityConfig != "" {
		set := strings.Split(complexityConfig, ",")
		for _, confElement := range set {
			name, complexity := Split(confElement, ":")
			if name == serviceName {
				complexityInt, err := strconv.Atoi(complexity)
				if err != nil || complexityInt < 0 {
					log.Println("WARNING: Invalid complexity numeric value")
				} else {
					//log.Printf("INFO: Complexity is %d", complexityInt)
					return complexityInt
				}
			}
		}
	}
	return 100
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
