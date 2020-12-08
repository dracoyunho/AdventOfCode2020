package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Execute executes the program - if it is an infinite loop, it performs one complete period and quits
func Execute(instructions map[int]string) (int, bool) {
	accumulator := 0
	completed := false
	current := 1
	executed := make(map[int]struct{})
	reNum := regexp.MustCompile("[0-9]+")

	// Start with the first instruction
	for true {
		// Before proceeding any further, check if the upcoming line is empty string - if it is, the program is done
		if instructions[current] == "" {
			completed = true
			break
		}

		// Every instruction can be split by space
		cmd := strings.Split(instructions[current], " ")[0]
		valStr := strings.Split(instructions[current], " ")[1]
		val, err := strconv.Atoi(reNum.FindString(valStr))
		if err != nil {
			log.Fatal(err)
		}
		if strings.HasPrefix(valStr, "-") {
			val *= -1
		}

		// Break if the command has been executed before; otherwise, execute the command
		if _, before := executed[current]; before {
			break
		}
		executed[current] = struct{}{}
		if cmd == "acc" {
			accumulator += val
			current++
		} else if cmd == "nop" {
			current++
		} else if cmd == "jmp" {
			current += val
		}
	}

	return accumulator, completed
}

func main() {
	// Reader
	path := "./input.txt"
	buf, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(buf)

	// Retrieve input
	input := make(map[int]string)
	lineNumber := 1
	for scanner.Scan() {
		input[lineNumber] = scanner.Text()
		lineNumber++
	}

	// P1: Trace a period of the infinite loop by storing a simple map - keys are line numbers, values are empty structs.
	// If checking the map for the given key does not return nil, then the line number was already in. At that point, stop execution of the game code.
	// At this point, return the value in the accumulator.
	finalValue, success := Execute(input)
	log.Println("P1 | Accumulator:", finalValue, "| Completed successfully:", success)

	// P2: There's no guarantee that the last jmp or nop is the one that needs to be fixed.
	// Instead, at every discovered jmp or nop, attempt flipping it to the opposite instruction in mem (the line number itself doesn't need to be returned)
	// If it doesn't work, then just proceed as normal until another jmp or nop is encountered
	// According to the puzzle, one of these is guaranteed to succeed
	for i := 1; i < len(input); i++ {
		if strings.HasPrefix(input[i], "jmp") {
			input[i] = "nop " + strings.Split(input[i], " ")[1]
		} else if strings.HasPrefix(input[i], "nop") {
			input[i] = "jmp " + strings.Split(input[i], " ")[1]
		} else {
			continue
		}
		finalValue, success = Execute(input)
		if success {
			log.Println("P2 | Modified line:", i, "| Accumulator:", finalValue, "| Completed successfully:", success)
			break
		}
		if strings.HasPrefix(input[i], "jmp") {
			input[i] = "nop " + strings.Split(input[i], " ")[1]
		} else if strings.HasPrefix(input[i], "nop") {
			input[i] = "jmp " + strings.Split(input[i], " ")[1]
		}
	}
}
