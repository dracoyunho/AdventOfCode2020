package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
)

// BuildRule returns the regexp pattern for a given rule as defined by a rulebook
func BuildRule(rulebook map[string]string, id string) string {
	// If the rule's definition in rulebook doesn't start with a letter, this is a simple return
	if match, err := regexp.MatchString(`^[^0-9]`, rulebook[id]); match {
		return rulebook[id]
	} else if err != nil {
		log.Fatal(err)
	}
	// Otherwise, substitute rule IDs in the definition string, and continue substituting until the definition string no longer contains digits
	pattern := rulebook[id]
	log.Println("Returned pattern is now:", pattern)
	reDigits := regexp.MustCompile(`[0-9]+`)
	rePipe := regexp.MustCompile(`\|`)
	for reDigits.MatchString(pattern) {
		splits := strings.Split(pattern, " ")
		for i := range splits {
			if reDigits.MatchString(splits[i]) {
				// If the rule def for this rule ID has a pipe, it must go into parentheses
				if rePipe.MatchString(rulebook[splits[i]]) {
					splits[i] = "( " + rulebook[splits[i]] + " )"
				} else {
					splits[i] = rulebook[splits[i]]
				}
			}
		}
		pattern = strings.Join(splits, " ")
		log.Println("Returned pattern is now:", pattern)
	}
	// The resulting pattern has way too much whitespace for a proper pattern
	pattern = strings.Join(strings.Split(pattern, " "), "")
	// The pattern for every rule also mandates that the first subrule/letter is the start of line and the pattern ends at EOL
	pattern = "^" + pattern + "$"
	return pattern
}

// P1 solves Puzzle 1
func P1() {
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

	// The input is divided into rules and entries by an empty line
	// Every rule is either a list of subrule IDs or an actual letter (in "")
	// Determine the base rules, those using a letter for definition, and store its ID
	// For every other rule, split it by | and assign that split to its current rule definition
	var rules map[string]string = make(map[string]string) // []string because of subrule ID pipe
	var lines []string
	for _, line := range input {
		if match, _ := regexp.MatchString(`^[0-9]`, line); match {
			rule := strings.Split(line, ":")
			id, def := rule[0], rule[1]
			if strings.Contains(def, "\"") {
				rules[id] = strings.Trim(def, " \"")
			} else {
				// Tidy up whitespace
				rules[id] = strings.TrimSpace(def)
			}
		} else if line != "" {
			lines = append(lines, line)
		}
	}
	log.Println("Rules:")
	for id := range rules {
		log.Println(id, ":", rules[id])
	}

	// Determine the regexp pattern for Rule 0
	pattern := BuildRule(rules, "0")
	reZero := regexp.MustCompile(pattern)
	log.Println("P1 | Rule 0 Pattern:", pattern)

	// Determine the number of lines matching the Rule 0 pattern
	matches := 0
	for _, line := range lines {
		if reZero.MatchString(line) {
			matches++
		}
	}
	log.Println("P1 | Matches to Rule 0:", matches)
}

// P2 solves Puzzle 2
func P2() {
	// Reader
	path := "./amended.txt"
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

	// The input is divided into rules and entries by an empty line
	// Every rule is either a list of subrule IDs or an actual letter (in "")
	// Determine the base rules, those using a letter for definition, and store its ID
	// For every other rule, split it by | and assign that split to its current rule definition
	var rules map[string]string = make(map[string]string) // []string because of subrule ID pipe
	var lines []string
	for _, line := range input {
		if match, _ := regexp.MatchString(`^[0-9]`, line); match {
			rule := strings.Split(line, ":")
			id, def := rule[0], rule[1]
			if strings.Contains(def, "\"") {
				rules[id] = strings.Trim(def, " \"")
			} else {
				// Tidy up whitespace
				rules[id] = strings.TrimSpace(def)
			}
		} else if line != "" {
			lines = append(lines, line)
		}
	}
	log.Println("Rules:")
	for id := range rules {
		log.Println(id, ":", rules[id])
	}

	// Determine the regexp pattern for Rule 0
	pattern := BuildRule(rules, "0")
	reZero := regexp.MustCompile(pattern)
	log.Println("P2 | Rule 0 Pattern:", pattern)

	// Determine the number of lines matching the Rule 0 pattern
	matches := 0
	for _, line := range lines {
		if reZero.MatchString(line) {
			matches++
		}
	}
	log.Println("P2 | Matches to Rule 0:", matches)
}

func main() {
	P1()
	P2()
}
