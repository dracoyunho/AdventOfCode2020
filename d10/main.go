package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

// FindOneThreeDiffs takes a map of ints and will return the number of times any two numbers are apart by 1 and by 3
// However, if a by-one is found, then by-three is bypassed
// If the given input is less than 2 elements, it just returns 0, 0
func FindOneThreeDiffs(input map[int]int) (int, int) {
	byOne := 0
	byThree := 0
	for k := range input {
		if _, ok := input[k+1]; ok {
			byOne++
		} else if _, ok := input[k+3]; ok {
			byThree++
		}
	}
	return byOne, byThree
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
	// The input file won't necessarily contain 0, but it should be there as the charging port is designated as 0
	var input = make(map[int]int)
	input[0] = 0
	for scanner.Scan() {
		val, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		input[val] = 0
	}

	// P1: Finding the one-diffs and three-diffs is as simple as sorting and iterating
	// Remember that diffs by 3 should be one higher than the return value, as this diff always exists at the top end
	// byOne, byThree := FindDiffs(input)
	byOne, byThree := FindOneThreeDiffs(input)
	log.Println("P1 | One-Diffs x Three-Diffs:", byOne*(byThree+1))

	// P2: Some relevant lines to P2:
	//   * Any given adapter can take an input 1, 2, or 3 jolts lower than its rating
	//   * In addition, your device has a built-in joltage adapter rated for 3 jolts higher than the highest-rated adapter in your bag
	// Therefore, the target joltage (the value every valid adapter chain must end with) is 3 + this maximum
	maxJolt := 0
	for k := range input {
		if k > maxJolt {
			maxJolt = k
		}
	}
	log.Println("P2 | Maximum Joltage:", maxJolt)

	// Now consider all the valid paths
	// Any given adapter can be the head of an adapter 1, 2, or 3 greater than itself
	// Example: the set of adapters are [1,2,3,4], then the target adapter is 4, and there are 7 chains:
	// 4 - 3 - 2 - 1 - 0
	//           \ 0
	//       \ 1 - 0
	//       \ 0
	//   \ 2 - 1 - 0
	//       \ 0
	//   \ 1 - 0
	// But if the 2-rated adapter were missing, then the number of paths drops to 3. Visualized:
	// 4
	// |\
	// 1 3
	// | |\
	// 0 0 1
	//     |
	//     0
	// Thus, every non-terminal leaf returns the sum of children below them.
	// This leads to an interesting phenomenon: because the children of any element must follow the within-3 rule, and because the available adapters don't change,
	//   then it's actually possible to substitute parts of the tree wholesale
	// Returning to the [1,2,3,4] example, the 1-adapter returns 1, and the 2-adapter returns 2:
	// 4 - 3 - {2}
	//       \ {1}
	//       \ 0
	//   \ {2}
	//   \ {1}
	// Thus, the 3-adapter returns the 0-path, the 1-adapter, and the 2-adapter, or 1+1+2=4, yielding the 7 valid paths:
	// 4 - {4}
	//   \ {2}
	//   \ {1}
	// Extending to 5 adapters [1,2,3,4,5], there are therefore 13 valid paths, as the 5-adapter's valid paths is just the sum of the 2-adapter, 3-adapter, and 4-adapter
	// This also works when the adapters are not consecutive. Example: [1,3,4,5]
	// With the subset [1,3,4], the 3-adapter has 2 valid paths (from 0 and 1) and the 4-adapter has 3 valid paths (from 1 and 3)
	// And therefore the 5-adapter, which may connect to the 3- and 4-adapters, has 2+3=5 valid paths
	// Neat!
	// Therefore, take this action:
	//   - Map keys -2 and -1 must have value 0, to prevent input[1-3] and the like from freaking out
	//   - Map key 0 must have value 1
	//   - For each map key above 0 and up to and including maxJolt, the key's value is input[key-3] + input[key-2] + input[key-1]
	//   - If any map key does not exist, set its value to 0
	//   and on and on incrementing the keys by 1 until the target is reached.
	input[-2] = 0
	input[-1] = 0
	input[0] = 1
	for i := 1; i <= maxJolt; i++ {
		if _, ok := input[i]; !ok {
			input[i] = 0
			continue
		}
		input[i] = input[i-3] + input[i-2] + input[i-1]
	}
	log.Print(fmt.Sprintf("P2 | Valid paths to Adapter %d: %d", maxJolt, input[maxJolt]))
}
