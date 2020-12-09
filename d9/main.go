package main

import (
	"bufio"
	"log"
	"os"
	"sort"
	"strconv"
)

// FindRogueValue will attempt to find a value in a slice of ints where, given some preamble size, the following int
//   is not a sum of any two elements in the preamble
// If such a value is found, the return values are the discovered value plus a success flag
// If not, the return value is 0 and the success flag is not active
func FindRogueValue(data []int, preambleSize int) (int, bool) {
	for i := preambleSize; i < len(data); i++ {
		var preamble = make([]int, preambleSize)
		copy(preamble, data[i-preambleSize:i])
		sort.Ints(preamble)
		lowIndex := 0
		highIndex := preambleSize - 1
		for preamble[lowIndex]+preamble[highIndex] != data[i] || preamble[lowIndex] == preamble[highIndex] {
			// Now it is known that the sum is not the target
			// If the low and high search indexes are 1 apart, then the target value is a rogue value - return success
			if highIndex-lowIndex == 1 {
				return data[i], true
			}
			if preamble[lowIndex]+preamble[highIndex] > data[i] {
				highIndex--
			}
			if preamble[lowIndex]+preamble[highIndex] < data[i] {
				lowIndex++
			}
		}
		// If the loop exits naturally, then a match was found, so it's time to move on to the next value
	}
	// If the loop exits naturally, then every number (that would have a properly-sized preamble) is not rogue
	return 0, false
}

// FindVulnerability finds the contiguous block of numbers in a slice and returns the sum of the lowest and highest values
// If it successfully did the above, it sets a success flag - if not, it returns 0 and no success (false)
func FindVulnerability(data []int, target int) (int, bool) {
	if len(data) < 2 {
		return 0, false
	}
	for size := 2; size < len(data); size++ {
		for start := 0; start+size-1 < len(data); start++ {
			subset := data[start : start+size]
			sum := 0
			for _, v := range subset {
				sum += v
			}
			if sum == target {
				var clone = make([]int, size)
				copy(clone, subset)
				sort.Ints(clone)
				return clone[0] + clone[size-1], true
			}
		}
	}
	return 0, false
}

func main() {
	// Reader
	path := "./input.txt"
	buf, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(buf)

	// Golang does this real nice thing with arrays and slices
	// Let's make an array of all the contents; slices can be made out of it as needed

	// Retrieve input
	var input []int
	for scanner.Scan() {
		val, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		input = append(input, val)
	}

	// P1: Start by looking at index 25 - it must be a sum of any two numbers in the previous 25 indexes
	// Like in Day 1's puzzles, sort this 25-long preamble, and then sum the lowest and highest values
	// If the resulting sum is higher than the desired sum, then the top end is too high - drop the index and try again
	// If the resulting sum is lower than the desired sum, then the bottom end is too low - up the index and try again
	// If the resulting sum is discovered, then return success - unless the two values used in the sum are of equal value
	// However, if it gets to the point where the two indexes collide (i.e. are equal), then the rogue value has been found
	rogue, ok := FindRogueValue(input, 25)
	log.Println("P1 | Rogue value:", rogue, "| Success:", ok)

	// P2: This puzzle could potentially be done without iterating through slice sizes, but this seems like a very non-trivial task
	// Instead, just loop through all of the potential slice sizes, which shouldn't be too big of a deal since the elements must be contiguous
	// There might be some kind of assurance that the numbers behind the rogue value can never sum to the rogue value,
	//   but this doesn't seem very likely
	vul, ok := FindVulnerability(input, rogue)
	log.Println("P2 | Vulnerability:", vul, "| Success:", ok)
}
