package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
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

	// P1: For each group, build up a slice of unique letters
	groupAnswers := make(map[string]struct{})
	sumYes := 0
	for i := 0; i < len(input); i++ {
		if len(strings.Split(input[i], "")) == 0 {
			sumYes += len(groupAnswers)
			// log.Println("P1 | DEBUG | Group Answers:", groupAnswers, "| Count:", len(groupAnswers))
			groupAnswers = make(map[string]struct{})
			continue
		}
		for _, char := range strings.Split(input[i], "") {
			groupAnswers[char] = struct{}{}
		}
	}
	// If there's still content in groupAnswers at this point then it's because the input doesn't contain a trailing empty string
	if len(groupAnswers) > 0 {
		sumYes += len(groupAnswers)
	}
	log.Print(fmt.Sprintf("P1 | Count of yes answers: %d", sumYes))

	// P2: For each group, associate each letter with occurrences; also count up each line in the group
	// If the group line count equals the occurrences then tick up the sum
	groupAnswerCounts := make(map[string]int)
	groupSize := 0
	sumAllYes := 0
	for i := 0; i < len(input); i++ {
		if len(strings.Split(input[i], "")) == 0 {
			for _, count := range groupAnswerCounts {
				if count == groupSize {
					sumAllYes++
				}
			}
			groupAnswerCounts = make(map[string]int)
			groupSize = 0
			continue
		}
		groupSize++
		for _, char := range strings.Split(input[i], "") {
			groupAnswerCounts[char]++
		}
	}
	// If there's still content in groupAnswerCounts at this point then it's because the input doesn't contain a trailing empty string
	if len(groupAnswerCounts) > 0 {
		for _, count := range groupAnswerCounts {
			if count == groupSize {
				sumAllYes++
			}
		}
	}
	log.Print(fmt.Sprintf("P2 | Count of group yes answers: %d", sumAllYes))
}
