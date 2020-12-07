package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
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

	type Bag struct {
		Descriptor string
		Colour     string
	}

	// Every rule exists only on one line, thankfully.
	// In regex form, (\w+ \w+) bags contain (\d+ \w+ \w+ bags?(, )?)+\.
	// Where shit gets real is that each rule is only one level deep - but **we have to go deeper.**
	// oui, je regrette tout
	// For now, just compile the ruleset one level deep. It should be good enough.
	ruleset := make(map[Bag]map[Bag]int)
	reBagsIn := regexp.MustCompile("\\d+ \\w+ \\w+") // Everything after the bag colour; if there are multiple, just FindAllString
	for _, rule := range input {
		ruleSplit := strings.Split(rule, " bags contain ")

		bagOuter := Bag{strings.Split(ruleSplit[0], " ")[0], strings.Split(ruleSplit[0], " ")[1]}

		bagsInner := reBagsIn.FindAllString(ruleSplit[1], -1)

		for _, bag := range bagsInner {
			bagInner := Bag{strings.Split(bag, " ")[1], strings.Split(bag, " ")[2]}
			bagInnerCount, err := strconv.Atoi(strings.Split(bag, " ")[0])
			if err != nil {
				log.Fatal(err)
			}
			ruleset[bagOuter] = map[Bag]int{
				bagInner: bagInnerCount,
			}
		}
	}

	// P1: Now the deep search has to happen.
}
