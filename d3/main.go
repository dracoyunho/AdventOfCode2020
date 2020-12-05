package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func ouchieCounter(field []string, startRight, startDown, cadenceRight, cadenceDown int) int {
	ouchies := 0
	positionRight := startRight
	positionDown := startDown
	for lineN, line := range field {
		if lineN%cadenceDown == 0 {
			rowSplit := strings.Split(line, "")
			if positionRight >= len(rowSplit) {
				positionRight %= len(rowSplit)
			}
			if rowSplit[positionRight] == "#" {
				ouchies++
			}
			positionRight += cadenceRight
			positionDown += cadenceDown
		}
	}
	log.Print(fmt.Sprintf("P2 | Cadence [R,D]: [%d,%d] | Ouchies: %d", cadenceRight, cadenceDown, ouchies))
	return ouchies
}

func main() {
	// P1: Traversal through the field can be simply performed by iteration and mod math.
	// This is because the tree pattern repeats infinitely to the right, which is the direction of traversal anyway.
	// As the iteration over lines proceeds, if the next index for pulling up the field value is beyond the bounds of the upcoming line,
	//   then simply take the next index mod the length of the line
	// For example, presuming a period of 15 and a starting index of 6:
	//   6 mod 15 = 6
	//   9 mod 15 = 9
	//   12 mod 15 = 12
	//   15 mod 15 = 0
	// In order to get the individual chars, just split each line by empty string.

	// Reader
	path := "./input.txt"
	buf, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(buf)

	// Retrieve input
	var input []string
	for scanner.Scan() {
		input = append(input, scanner.Text())
	}

	// Time to slam into trees
	ouchies := 0
	positionRight := 0
	positionDown := 0
	cadenceRight := 3
	cadenceDown := 1
	for lineN, line := range input {
		if lineN%cadenceDown == 0 {
			rowSplit := strings.Split(line, "")
			if positionRight >= len(rowSplit) {
				positionRight %= len(rowSplit)
			}
			if rowSplit[positionRight] == "#" {
				ouchies++
			}
			positionRight += cadenceRight
			positionDown += cadenceDown
		}
	}
	log.Print(fmt.Sprintf("P1 | Ouchies: %d", ouchies))

	// P2: This is probably a good time to turn the content above into a function.
	// The original stuff is preserved up top for posterity, though.
	ouchieProduct := 1
	ouchieProduct *= ouchieCounter(input, 0, 0, 1, 1)
	ouchieProduct *= ouchieCounter(input, 0, 0, 3, 1)
	ouchieProduct *= ouchieCounter(input, 0, 0, 5, 1)
	ouchieProduct *= ouchieCounter(input, 0, 0, 7, 1)
	ouchieProduct *= ouchieCounter(input, 0, 0, 1, 2)
	log.Print(fmt.Sprintf("P2 | Ouchie Product: %d", ouchieProduct))
}
