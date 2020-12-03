package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	// P1 requires regex to split the lines in data into two pieces:
	//   1) The requirement
	//   2) The password in storage
	// This may be split on the string ": "
	// The requirement string may then be further split into a count range and required char on the string " "
	// The count range may then be split into a min and max on the string "-"

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

	// Split pw from requirement
	var passwords []string
	var requirements []string
	rePassword := regexp.MustCompile(": ")
	for i := 0; i < len(input); i++ {
		split := rePassword.Split(input[i], -1)
		requirements = append(requirements, split[0])
		passwords = append(passwords, split[1])
	}

	// Split count from char
	var reqCount []string
	var reqChar []string
	reReq := regexp.MustCompile(" ")
	for i := 0; i < len(requirements); i++ {
		split := reReq.Split(requirements[i], -1)
		reqCount = append(reqCount, split[0])
		reqChar = append(reqChar, split[1])
	}

	// Split min from max
	var countMin []int
	var countMax []int
	reCount := regexp.MustCompile("-")
	for i := 0; i < len(reqCount); i++ {
		split := reCount.Split(reqCount[i], -1)
		min, errMin := strconv.ParseInt(split[0], 10, 0)
		if errMin != nil {
			log.Fatal(errMin)
		}
		max, errMax := strconv.ParseInt(split[1], 10, 0)
		if errMax != nil {
			log.Fatal(errMax)
		}
		countMin = append(countMin, int(min))
		countMax = append(countMax, int(max))
	}

	// Now count the number of OK pws
	valid := 0
	for i := 0; i < len(passwords); i++ {
		count := strings.Count(passwords[i], reqChar[i])
		if count >= countMin[i] && count <= countMax[i] {
			valid++
		}
	}
	log.Print(fmt.Sprintf("P1 | Valid passwords: %d", valid))

	// With Puzzle 2, the countMin and countMax slices now indicate the positions of where characters should be searched in the password (not indexes)
	// Before checking that the char exists at the given locations, the index (countMin[i] - 1 or countMax[i] - 1) should be checked that it's within the length of the pw string
	// Otherwise skip checking for that index
	valid = 0
	for i := 0; i < len(passwords); i++ {
		posMin := countMin[i] - 1
		posMax := countMax[i] - 1
		hits := 0
		chars := strings.Split(passwords[i], "")
		if posMin < len(chars) {
			if chars[posMin] == reqChar[i] {
				hits++
			}
		}
		if posMax < len(chars) {
			if chars[posMax] == reqChar[i] {
				hits++
			}
		}
		if hits == 1 {
			valid++
		}
	}
	log.Print(fmt.Sprintf("P2 | Valid passwords: %d", valid))
}
