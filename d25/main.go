package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

const (
	// InputFilePath is the path to the input for this puzzle
	InputFilePath string = "./input.txt"
	// DefaultSubject is the default subject for this puzzle, which is 7
	DefaultSubject int = 7
	// ModulusKey is the default modulus for subject transformation, which is 20201227
	ModulusKey int = 20201227
)

// Transform performs n iterations of transformation, given some initial value and a subject
func Transform(i, s, n int) int {
	for iter := 0; iter < n; iter++ {
		i = IterateTransform(i, s)
	}
	return i
}

// IterateTransform performs one transformation iteration given some initial value and a subject
func IterateTransform(i, s int) int {
	i *= s
	i %= ModulusKey
	return i
}

func main() {
	// Reader
	buf, err := os.Open(InputFilePath)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(buf)

	// Retrieve input
	var input []string
	for scanner.Scan() {
		input = append(input, scanner.Text())
	}

	// The two strings are the card and door public key (kpc & kpd), both associated with some secret key (ks)
	kpc, kpcErr := strconv.Atoi(input[0])
	if kpcErr != nil {
		log.Fatal(kpcErr)
	}
	kpd, kpdErr := strconv.Atoi(input[1])
	if kpdErr != nil {
		log.Fatal(kpdErr)
	}

	// Transforming the default Subject some lc times should yield kpc; transforming the default Subject some ld times should yield kpd
	var loops, lc, ld int = 0, 0, 0
	var transform int = 1
	for lc == 0 || ld == 0 {
		loops++
		transform = IterateTransform(transform, DefaultSubject)
		// log.Println("Iteration", loops, "| Transform is now", transform) // Don't turn this on unless you enjoy wasting roughly 17 minutes
		if transform == kpc {
			lc = loops
			log.Println("Set lc to", loops)
		}
		if transform == kpd {
			ld = loops
			log.Println("Set ld to", loops)
		}
	}

	// With lc and ld in hand, validate that the ksc and ksd from transforming kpd lc times and transforming kpc ld times are equal
	var ks, ksc, ksd int = 0, Transform(1, kpd, lc), Transform(1, kpc, ld)
	if ksc == ksd {
		ks = ksc
	} else {
		log.Fatalln("Encryption keys from card and door don't match! Card:", ksc, "| Door:", ksd)
	}
	log.Println("P1 | Encryption key:", ks)

	// P2:
	// I thought this was going to end with "Your vacation was a COVID-19 fever dream, haha, get wrecked kid" but instead it ends with a broken soft serve machine.
	// I don't know which one's worse tbh
}
