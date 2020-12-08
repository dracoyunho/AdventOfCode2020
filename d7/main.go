package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Bag is a simple struct for the Bags in the input
type Bag struct {
	Descriptor string
	Colour     string
}

// BagsEqual compares two bags and returns true if their properties match exactly
func BagsEqual(a, b Bag) bool {
	// Bags are equal if Descriptor and Colour match
	if a.Descriptor == b.Descriptor && a.Colour == b.Colour {
		return true
	}
	return false
}

// ContainsBag checks against a ruleset to see if the desired bag will be in the outer bag
// It will return true if one of the bags inside the outer bag is the target and false otherwise
// func ContainsBag(outerBag Bag, targetBag Bag, ruleset map[Bag]map[Bag]int) bool {
// 	bagsInside := ruleset[outerBag] // Set to a map[Bag]int
// 	for bagInside := range bagsInside {
// 		if BagsEqual(bagInside, targetBag) {
// 			return true
// 		}
// 	}
// 	return false
// }

// WillContainBag searches recursively for a target bag given some outermost bag and a ruleset of bags in a bag
func WillContainBag(outerBag, targetBag Bag, ruleset map[Bag]map[Bag]int) bool {
	// The terminating cases are 1) the outer bag's contents contains the target (true) or 2) the outer bag contains no bags (false)
	// If neither of these are met, then return true if the bag within is going to contain the target
	// If the bag inside won't contain the target, don't return false until all the bags inside are checked
	bagsInside := ruleset[outerBag]
	if bagsInside == nil {
		return false
	}
	for bagInside := range bagsInside {
		if BagsEqual(bagInside, targetBag) {
			return true
		}
		if WillContainBag(bagInside, targetBag, ruleset) {
			return true
		}
	}
	return false
}

// BagsInside searches recursively to find out the total number of bags contained within a given bag given some ruleset of bags in a bag
func BagsInside(outerBag Bag, ruleset map[Bag]map[Bag]int) int {
	// The terminating case is when there is no bag inside the given bag; this returns 0
	bagsInside := ruleset[outerBag]
	if bagsInside == nil {
		return 0
	}
	totalBagsInside := 0
	for bagInside, countBagInside := range bagsInside {
		totalBagsInside += countBagInside * (1 + BagsInside(bagInside, ruleset))
	}
	return totalBagsInside
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
	var input []string
	for scanner.Scan() {
		input = append(input, scanner.Text())
	}

	// Every rule exists only on one line, thankfully.
	// In regex form, (\w+ \w+) bags contain (\d+ \w+ \w+ bags?(, )?)+\.
	// Where shit gets real is that each rule is only one level deep - but **we have to go deeper.**
	// oui, je regrette tout
	// For now, just compile the ruleset one level deep. It should be good enough.
	ruleset := make(map[Bag]map[Bag]int)
	reBagsIn := regexp.MustCompile("\\w+ \\w+ \\w+") // Everything after the bag colour; if there are multiple, just FindAllString
	for _, rule := range input {
		ruleSplit := strings.Split(rule, " bags contain ")

		bagOuter := Bag{strings.Split(ruleSplit[0], " ")[0], strings.Split(ruleSplit[0], " ")[1]}

		bagsInner := reBagsIn.FindAllString(ruleSplit[1], -1)

		if bagsInner[0] != "no other bags" {
			bagsInnerMap := make(map[Bag]int)
			for _, bag := range bagsInner {
				bagInner := Bag{strings.Split(bag, " ")[1], strings.Split(bag, " ")[2]}
				bagInnerCount, err := strconv.Atoi(strings.Split(bag, " ")[0])
				if err != nil {
					log.Fatal(err)
				}
				bagsInnerMap[bagInner] = bagInnerCount
			}
			ruleset[bagOuter] = bagsInnerMap
		}
	}

	// P1: Now the deep search has to happen.
	// For every outer bag in a rule, pull up the rules for the inner bags, if they exist.
	// If for a given outer bag there is no rule, then it is time to end recursion.
	// Alternatively, once a bag of Descriptor shiny and Colour gold is found, cease recursion and increment the counter.
	matches := 0
	targetBag := Bag{"shiny", "gold"}
	for outerBag := range ruleset {
		if WillContainBag(outerBag, targetBag, ruleset) {
			matches++
		}
	}
	log.Println("P1 | Bags containing a", targetBag.Descriptor, targetBag.Colour, "bag:", matches)

	// P2: In this case there isn't a wide search but a deep search.
	// There is a single rule for shiny gold bags, but (just looking at the rule for it) there are a significant number of bags for those.
	// Recursion doesn't stop until the bag being looked up comes up with a nil map.
	// The return value for the recursive function is the number of bags inside.
	// For example:
	// Bag A contains 1 Bag B and 1 Bag C:
	//   Bag B contains 2 Bag D
	//     Bag D is empty
	//   Bag C is empty
	// Bag D returns 0
	// Bag C returns 0
	// Bag B returns 2*(D+1) = 2*1 = 2
	// Bag A returns 1*(B+1)+1*(C+1) = 1*(2+1)+1*(0+1) = 1*3+1*1 = 3+1 = 4
	// The +1 when calculating bags inside is to add the containing bag to the contained bags.
	log.Println("P2 | Bags inside a", targetBag.Descriptor, targetBag.Colour, "bag:", BagsInside(targetBag, ruleset))
}
