package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

const (
	// Input is the, well, input, for the puzzle.
	Input string = "469217538"
	// Test is the test string for the puzzle.
	Test string = "389125467"
)

// ListContainsValue checks if a given list contains a given value and returns correspondingly
func ListContainsValue(l []int, v int) bool {
	for i := range l {
		if v == l[i] {
			return true
		}
	}
	return false
}

// PrintCupList ingests a list of cups, where the indexes are the cups and the values are the cup that come after, as well as a starting/current cup
// The starting/current cup is indicated with parentheses
// It will also return the string that would be printed; printing may be turned off if only this string is desired
func PrintCupList(cups []int, start int, verbose bool) string {
	var output []string
	var next int = cups[start]
	output = append(output, "("+fmt.Sprint(start)+")")
	for next != start {
		output = append(output, fmt.Sprint(next))
		next = cups[next]
	}
	if verbose {
		log.Println(strings.Join(output, ","))
	}
	return strings.Join(output, ",")
}

func main() {
	input := strings.Split(Input, "")
	var inputVals []int
	for i := range input {
		v, err := strconv.Atoi(input[i])
		if err != nil {
			log.Fatal(err)
		}
		inputVals = append(inputVals, v)
	}
	// Spoilers from P2: Using the Golang ring is going to take a... long time. O(n) at the size of 1000000, 10M times, is going to take _forever_
	// Instead, consider what it means for these cups: All that matters is the value ahead of a given cup
	// Consider cups 1 to 5 in a loop; 1 is behind 2, which is behind 3, which is behind 4, which is behind 5
	// Consider now an array, designed such that:
	//   - The index represents a number printed onto a cup
	//   - The value represents the number printed on the next cup
	// For the example of cups from 1 -> 2 -> 3 -> 4 -> 5 -> 1, this would be:
	// Cup/Index: 0 1 2 3 4 5
	// Next Cup: [0,2,3,4,5,1]
	// For the example input, 3 -> 8 -> 9 -> 1 -> 2 -> 5 -> 4 -> 6 -> 7 -> 3, this is:
	// Cup/index: 0 1 2 3 4 5 6 7 8 9
	// Next Cup: [0,2,5,8,6,4,7,3,9,1]
	// Let's take this example and do one round of iteration on it.
	// The current cup value is current = 3.
	// The snip is the next three cups: snip = [cups[current], cups[cups[current], cups[cups[cups[current]]]] = [8, 9, 1]
	// The current cup value is immediately updated to be the cup after this snip, or cups[cups[cups[cups[current]]]]], in this case, 2:
	// Cup/index: 0 1 2 3 4 5 6 7 8 9
	// Next Cup: [0,2,5,2,6,4,7,3,9,1]
	// The next cup after, i.e. the target, will be initialized as target = current - 1 = 2, which is a valid target, as none of the three values in snip correspond to this target.
	// Otherwise, the target iterates downward, looping to the maximum cup value upon reaching 0 (so 0 is never a valid cup number)
	// To stitch this back together:
	//   - Set the cup indicated by the tail of the snip, i.e. cups[snip[2]], to be the current value of the target; in this example, cups[2] is 5, so set cups[snip[2]] = cups[1] = 5
	//   - Set the cup indicated by the target, i.e. cups[target], to be the head of the snip; in this example, cups[target] is 2 and snip[0] is 8, so set cups[target] = cups[2] = 8
	// This yields:
	// Cup/index: 0 1 2 3 4 5 6 7 8 9
	// Next Cup: [0,5,8,2,6,4,7,3,9,1]
	// Since current is still 3, then the chain goes: 3 --> 2 --> 8 --> 9 --> 1 --> 5 --> 4 --> 6 --> 7 --> 3
	// This is as expected from the example from Puzzle 1.

	// By the way: the 0 at 0 is just there to fill space. **please for the love of god never call cups[0], bad times will happen**

	// P1: Nine cups, from 1 to 9, over 100 iterations
	var p1cups []int = make([]int, 10)
	var current int = inputVals[0]
	for i := range input {
		// Every non-zero cup needs to be assigned such that the index of the cups array is the cup itself and the value of cups[index] is the next cup
		c := inputVals[i]
		v := inputVals[(i+1)%len(inputVals)]
		p1cups[c] = v
	}
	log.Println("Initial cups:", PrintCupList(p1cups, current, false))
	for it := 1; it <= 100; it++ {
		log.Println("Move", it, "| Cups:", PrintCupList(p1cups, current, false))

		var snip []int = []int{p1cups[current], p1cups[p1cups[current]], p1cups[p1cups[p1cups[current]]]}
		log.Println("Move", it, "| Extracted cups", snip)

		p1cups[current] = p1cups[p1cups[p1cups[p1cups[current]]]]

		var target int = current - 1
		if target == 0 {
			target = len(p1cups) - 1
		}
		for ListContainsValue(snip, target) {
			target--
			if target == 0 {
				target = len(p1cups) - 1
			}
		}
		log.Println("Move", it, "| Insert after:", target)

		// First set the cup indicated by the tail of the snip, i.e. cups[snip[2]], to be the current value of the target
		p1cups[snip[2]] = p1cups[target]
		// Second set the cup indicated by the target, i.e. cups[target], to be the head of the snip
		p1cups[target] = snip[0]
		// Third advance to the next cup in line
		current = p1cups[current]
	}
	// P1: The answer to P1 requires following cups starting from cup 1 until it wraps around
	log.Println("P1 | Cups, starting with 1:", PrintCupList(p1cups, 1, false))

	// P2: this crab is cancerous, my god
	// 1 Million cups, 1 - 1,000,000, over 10 Million iterations
	// The first nine cups follow the order of the input, and after that, they're all ascending numerical order
	p2cap := 1 * 1000 * 1000
	p2iter := 10 * 1000 * 1000
	var p2cups []int = make([]int, p2cap+1)
	current = inputVals[0]
	for i := 1; i <= p2cap; i++ {
		// Because i starts at 1, the current cup is indicated by inputVals[i-1] and the value to point the current cup to the next is indicated by inputVals[i]
		// Of course, this becomes invalid when i == 9, or when i-1 == 8, i.e. the last current cup in inputVals
		// Hence, when i-1 == 8, merely set v to 10 - it would be the next cup in the ring
		// After that, the value of i
		// At this point, instead of wrapping around to the beginning of the input, values just keep increasing, starting at 10 and up to the cap
		// Upon i == p2cap, this is now the final cup, and it should wrap back to the first cup
		if i < len(inputVals) {
			c := inputVals[i-1]
			v := inputVals[i]
			p2cups[c] = v
		} else if i == 9 {
			c := inputVals[i-1]
			v := i + 1
			p2cups[c] = v
		} else if i == p2cap {
			c := i
			v := inputVals[0]
			p2cups[c] = v
		} else {
			c := i
			v := i + 1
			p2cups[c] = v
		}
	}
	// In theory what worked for P1 should now work for P2, since there's much less "page-flipping"
	for it := 1; it <= p2iter; it++ {
		// log.Println("Move", it, "is now executing...") // Don't you dare turn this line on if either cup or iteration cap is more than, say, 10000.

		var snip []int = []int{p2cups[current], p2cups[p2cups[current]], p2cups[p2cups[p2cups[current]]]}

		p2cups[current] = p2cups[p2cups[p2cups[p2cups[current]]]]

		var target int = current - 1
		if target == 0 {
			target = len(p2cups) - 1
		}
		for ListContainsValue(snip, target) {
			target--
			if target == 0 {
				target = len(p2cups) - 1
			}
		}

		// First set the cup indicated by the tail of the snip, i.e. cups[snip[2]], to be the current value of the target
		p2cups[snip[2]] = p2cups[target]
		// Second set the cup indicated by the target, i.e. cups[target], to be the head of the snip
		p2cups[target] = snip[0]
		// Third advance to the next cup in line
		current = p2cups[current]
	}
	// The answer to P2 is the product of the two cups after cup numbered 1
	log.Println("P2 | The cups after cup #1:", p2cups[1], "&", p2cups[p2cups[1]], "| Their product:", fmt.Sprint(p2cups[1]*p2cups[p2cups[1]]))
}
