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

// ResolveExpression takes a slice of chars and returns a value corresponding to the solving of the expression where + and * have equal precedence
func ResolveExpression(chars []string) int {
	// The possible chars to be encountered are:
	//   1) A digit, "0" to "9"
	//   2) An operator, "+" or "*"
	//   3) An opening parenthesis, "("
	//   4) A closing parenthesis, ")"
	//   5) Any other char
	reDigit := regexp.MustCompile(`[0-9]`)
	reOp := regexp.MustCompile(`[\+\*]`)
	reOpen := regexp.MustCompile(`\(`)
	reClose := regexp.MustCompile(`\)`)
	var registers map[int]int = make(map[int]int)
	var results map[int]int = make(map[int]int)
	var operators map[int]string = make(map[int]string)
	var depth int = 0
	registers[depth] = 0
	results[depth] = 0
	operators[depth] = "+"
	for index := range chars {
		// Case 1: Digit. Multiply the current register by 10 and add the current digit, and assign back to the register.
		// Case 2: Operator. Evaluate the operator currently in memory between the result and the register, setting the value to result, and set the register to 0. Then, set operator memory to this operator.
		// Case 3: Opening parenthesis. Increment depth by 1. Initialize the new memory depths (register, result = 0 and operator = +)
		// Case 4: Closing parenthesis. Perform the operation defined by the current operator depth. Assign the value to the register above this depth by one level. Then, delete this depth from all memory stores.
		// Case 5: Any other char. Ignore.
		if reDigit.MatchString(chars[index]) {
			d, err := strconv.Atoi(chars[index])
			if err != nil {
				log.Fatal(err)
			}
			registers[depth] = registers[depth]*10 + d
		} else if reOp.MatchString(chars[index]) {
			if operators[depth] == "+" {
				results[depth] += registers[depth]
			} else if operators[depth] == "*" {
				results[depth] *= registers[depth]
			} else {
				log.Fatal("Operator", operators[depth], "is not recognized as a valid operator!")
			}
			registers[depth] = 0
			operators[depth] = chars[index]
		} else if reOpen.MatchString(chars[index]) {
			depth++
			registers[depth] = 0
			results[depth] = 0
			operators[depth] = "+"
		} else if reClose.MatchString(chars[index]) {
			if operators[depth] == "+" {
				registers[depth-1] = results[depth] + registers[depth]
			} else if operators[depth] == "*" {
				registers[depth-1] = results[depth] * registers[depth]
			} else {
				log.Fatal("Operator", operators[depth], "is not recognized as a valid operator!")
			}
			delete(registers, depth)
			delete(results, depth)
			delete(operators, depth)
			depth--
		}
	}
	// At end of line, perform the final operation
	if operators[depth] == "+" {
		results[depth] += registers[depth]
	} else if operators[depth] == "*" {
		results[depth] *= registers[depth]
	} else {
		log.Fatal("Operator", operators[depth], "is not recognized as a valid operator!")
	}

	// Sanity check
	if depth != 0 {
		log.Fatal("At the conclusion of this expression, the depth was not 0 - this program needs debugging!")
	}

	return results[depth]
}

// AdvResolveExpression takes a slice of chars and returns a value corresponding to the solving of the expression where + is performed before *
func AdvResolveExpression(chars []string) int {
	reDigit := regexp.MustCompile(`[0-9]`)
	reAdd := regexp.MustCompile(`\+`)
	reMult := regexp.MustCompile(`\*`)
	reOp := regexp.MustCompile(`[\+\*]`)
	reOpen := regexp.MustCompile(`\(`)
	reClose := regexp.MustCompile(`\)`)

	var subexpressions map[int][]string = make(map[int][]string)
	var depth int = 0
	// Identify subexpressions created by parentheses
	// Upon closure of a subexpression, resolve it and return its result
	for index := range chars {
		if reDigit.MatchString(chars[index]) || reOp.MatchString(chars[index]) {
			subexpressions[depth] = append(subexpressions[depth], chars[index])
		} else if reOpen.MatchString(chars[index]) {
			depth++
		} else if reClose.MatchString(chars[index]) {
			// Upon subexp closure, resolve it and append its value as string to the subexp above
			// The subexp to be closed is identified with the current depth
			seVal := AdvResolveExpression(subexpressions[depth])
			delete(subexpressions, depth)
			depth--
			subexpressions[depth] = append(subexpressions[depth], fmt.Sprintf("%d", seVal))
		}
	}
	if depth != 0 {
		log.Fatal("Depth is currently", depth, "when it should be 0!")
	}

	// Now resolve all additions by again iterating left to right; ignore multiplication
	// Because there is only one operation after addition to be done, int storage is OK
	var postAdditionValues []int
	var result, registry int = 0, 0
	for index := range subexpressions[depth] {
		if reDigit.MatchString(subexpressions[0][index]) {
			val, err := strconv.Atoi(subexpressions[0][index])
			if err != nil {
				log.Fatal(err)
			}
			registry = registry*10 + val
		} else if reAdd.MatchString(subexpressions[0][index]) {
			result += registry
			registry = 0
		} else if reMult.MatchString(subexpressions[0][index]) {
			postAdditionValues = append(postAdditionValues, result+registry)
			registry = 0
			result = 0
		}
	}
	// At EOL, perform a final addition and submit to postAdditionValues
	// e.g. 2 * 23 EOL --> 23 is in registry; it is combined with result (which is 0) and submitted to postAdditionValues
	// e.g. 3 * 5 + 3 EOL --> 3 is in registry; it is combined with result (which is 5) and submitted to postAdditionValues
	postAdditionValues = append(postAdditionValues, result+registry)

	// Now perform final multiplications
	result = 1
	for _, val := range postAdditionValues {
		result *= val
	}

	return result
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

	var p1Sum int = 0
	for _, line := range input {
		result := ResolveExpression(strings.Split(line, ""))
		p1Sum += result
		log.Println("P1 | Result:", result, "| Expression:", line)
	}
	log.Println("P1 | Result sum:", p1Sum)

	var p2Sum int = 0
	for _, line := range input {
		result := AdvResolveExpression(strings.Split(line, ""))
		p2Sum += result
		log.Println("P2 | Result:", result, "| Expression:", line)
	}
	log.Println("P2 | Result sum:", p2Sum)
}
