package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

func main() {
	/*
	 * Strategy:
	 *   The naive method is to parse systematically through every pair until a hit is found - there's only one, anyway
	 *   Slightly less naive is to:
	 *     1. Sort low to high
	 *     2. Evaluate the sum of index 0 and index length-1
	 *     3. If the sum is > 2020, then one of the two numbers must be lowered - and the only way to accomplish this is to decrease the higher end; -- the higher end index
	 *     4. If the sum is < 2020, then one of the two numbers must be raised - and the only way to accomplish this is to increase the lower end; ++ the lower end index
	 */

	// Reader
	path := "./input.txt"
	buf, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(buf)

	// Gather input
	var input []int
	for scanner.Scan() {
		inputInt, err := strconv.ParseInt(scanner.Text(), 10, 0) // return should be dependent on bitSize and put out plain int but I guess not?
		if err != nil {
			log.Fatal(err)
		}
		input = append(input, int(inputInt))
	}
	sort.Ints(input)

	// Now search for the two-value problem
	low := 0
	high := len(input) - 1
	for input[low]+input[high] != 2020 {
		if input[low]+input[high] < 2020 {
			low++
		} else {
			high--
		}
	}

	// Out with it
	log.Print(fmt.Sprintf("P1: Low: %d | High: %d | Result: %d", input[low], input[high], input[low]*input[high]))

	/*
	 * The three-value problem is a little trickier.
	 * Let's consider a "base" index. Subtracting the value for the base index from 2020 reveals the target sum for the low and high indexes.
	 * However, it is no longer a guarantee that there will be a match. The search can stop when the low and high indexes collide.
	 * If no match is found at that point, then the base index is raised by one, the low index is set to one above the base index, and the high index is reset.
	 * This is repeated until a match is found.
	 */

	base := 0
	low = base + 1
	high = len(input) - 1
	for input[low]+input[high] != 2020-input[base] {
		if high-low == 1 {
			// At this point changing low up or high down will make the indexes collide, so it is time to reset
			base++
			low = base + 1
			high = len(input) - 1
		} else if input[low]+input[high] < 2020-input[base] {
			low++
		} else {
			high--
		}
	}

	// Out with it
	log.Print(fmt.Sprintf("P2: Base: %d | Low: %d | High: %d | Result: %d", input[base], input[low], input[high], input[base]*input[low]*input[high]))
}
