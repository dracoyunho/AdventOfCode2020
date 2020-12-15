package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

// TermMemoryShift takes in a term number (cardinal, e.g. 1st, 2nd) and a 2-long int slice representing the two previous terms
//   that a term has appeared in (e.g. [7th, 3rd], or [7, 3])
// It will assign the first value to the second and the term input to the first slot
// In the example above, if the 10th term is submitted, then [7, 3] becomes [10, 7]
func TermMemoryShift(term int, recall [2]int) [2]int {
	recall[1] = recall[0]
	recall[0] = term
	return recall
}

// PatternSolve returns the value of the term number specified given some slice of ints that start the pattern.
func PatternSolve(initial []int, finalTerm int) int {
	// Spawn a map of ints to a slice of two ints - the keys are numbers that appear in the pattern
	// The values are the two most recent terms the key appears in the pattern
	// If the key does not exist in the map, then the second-most recent term is 0
	termMem := make(map[int][2]int)
	var pattern []int
	for term := 0; term < finalTerm; term++ {
		if term < len(initial) {
			pattern = append(pattern, initial[term])
			termMem[pattern[term]] = [2]int{term + 1, 0}
		} else {
			if termMem[pattern[term-1]][1] == 0 {
				pattern = append(pattern, 0)
				termMem[0] = TermMemoryShift(term+1, termMem[0])
			} else {
				pattern = append(pattern, termMem[pattern[term-1]][0]-termMem[pattern[term-1]][1])
				if _, exists := termMem[pattern[term]]; exists {
					termMem[pattern[term]] = TermMemoryShift(term+1, termMem[pattern[term]])
				} else {
					termMem[pattern[term]] = [2]int{term + 1, 0}
				}
			}
		}
	}
	return pattern[finalTerm-1]
}

func main() {
	// Reader
	path := "./input.txt"
	buf, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(buf)

	// Interpret input
	var input []string
	for scanner.Scan() {
		input = append(input, scanner.Text())
	}
	for index := range input {
		initialStrings := strings.Split(input[index], ",")
		var initials []int
		for _, str := range initialStrings {
			num, err := strconv.Atoi(str)
			if err != nil {
				log.Fatal(err)
			}
			initials = append(initials, num)
		}
		log.Println("P1 | INITIALS:", initials, "| 2020th term:", PatternSolve(initials, 2020))
		log.Println("P2 | INITIALS:", initials, "| 30000000th term:", PatternSolve(initials, 30000000))
	}
}
