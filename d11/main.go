package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// Position is the coordinates of a deck spot where the top-left corner is (0, 0) and increases rightward and downward
type Position struct {
	X int
	Y int
}

// PrintDeck neatly prints a deck - and also returns the number of filled seats
// If only the number of filled seats is desired, specify false to print flag
func PrintDeck(deck map[Position]string, doPrint bool) int {
	// First check that the first element in a new row exists
	// Then generate a new slice for that line, and append to it until the (X, Y) position is no longer valid
	line := 0
	done := false
	filled := 0
	for !done {
		col := 0
		pos := Position{line, col}
		if _, ok := deck[pos]; ok {
			lineDone := false
			var lineSlice []string
			for !lineDone {
				pos = Position{line, col}
				if _, ok := deck[pos]; ok {
					lineSlice = append(lineSlice, deck[pos])
					if deck[pos] == "#" {
						filled++
					}
					col++
				} else {
					lineDone = true
				}
			}
			if doPrint {
				log.Println("DECK ROW", line, "|\t", lineSlice)
			}
			line++
		} else {
			done = true
		}
	}
	return filled
}

// ResolveDeck ingests a deck map and iterates through the seating rules until nothing changes - if the deck is chaotic, this will never end!
func ResolveDeck(deck map[Position]string, version int) map[Position]string {
	run := 0
	changes := 1
	var newDeck = make(map[Position]string)
	for changes > 0 {
		run++
		changes = 0
		newDeck = make(map[Position]string)
		for pos, state := range deck {
			if state == "." {
				newDeck[pos] = deck[pos]
			} else if state == "L" {
				if TransitionEmpty(deck, pos, version) {
					changes++
					newDeck[pos] = "#"
				} else {
					newDeck[pos] = "L"
				}
			} else if state == "#" {
				if TransitionFilled(deck, pos, version) {
					changes++
					newDeck[pos] = "L"
				} else {
					newDeck[pos] = "#"
				}
			}
		}
		log.Print(fmt.Sprintf("P%d | Run %d | Changes: %d", version, run, changes))
		deck = newDeck

		// DEBUGGING
		// if version == 2 {
		// 	PrintDeck(deck, true)
		// }
	}
	return newDeck
}

// TransitionEmpty returns:
// V1: True if, given a deck map and an unfilled seat position, there are no filled seats around the given unfilled seat position
//     It will return false in any other scenario (i.e. a filled seat around the unfilled seat or the position isn't an unfilled seat)
// V2: True if, given a deck map and an unfilled seat position, there are no filled seats in the line of sight of the given unfilled seat position
//     It will return false in any other scenario (i.e. a filled seat is visible from the unfilled seat or the position isn't an unfilled seat)
func TransitionEmpty(deck map[Position]string, pos Position, version int) bool {
	for relX := -1; relX <= 1; relX++ {
		for relY := -1; relY <= 1; relY++ {
			if deck[pos] != "L" {
				return false // Just in case lol
			}
			if relX == 0 && relY == 0 {
				continue
			}
			if version == 1 {
				relPos := Position{pos.X + relX, pos.Y + relY}
				if val, ok := deck[relPos]; ok {
					if val == "#" {
						return false
					}
				}
			} else if version == 2 {
				mult := 1
				for true {
					relPos := Position{pos.X + relX*mult, pos.Y + relY*mult}
					val, ok := deck[relPos]
					if ok {
						if val == "L" {
							break
						} else if val == "#" {
							return false
						} else {
							mult++
						}
					} else {
						break
					}
				}
			}
		}
	}
	return true
}

// TransitionFilled returns
// V1: True if, given a deck map and a filled seat position, there are 4 or more filled seats around the given filled seat position
//     It will return false in any other scenario (i.e. a filled seat with 3 or less filled seats around it or the position isn't a filled seat)
// V2: True if, given a deck map and a filled seat position, there are 5 or more filled seats in the line of sight of the given filled seat position
//     It will return false in any other scenario (i.e. a filled seat with 4 or less filled seats in view or the position isn't a filled seat)
func TransitionFilled(deck map[Position]string, pos Position, version int) bool {
	filledNeighbours := 0
	for relX := -1; relX <= 1; relX++ {
		for relY := -1; relY <= 1; relY++ {
			if deck[pos] != "#" {
				return false // Just in case
			}
			if relX == 0 && relY == 0 {
				continue
			}
			if version == 1 {
				relPos := Position{pos.X + relX, pos.Y + relY}
				if val, ok := deck[relPos]; ok {
					if val == "#" {
						filledNeighbours++
					}
				}
			} else if version == 2 {
				mult := 1
				for true {
					relPos := Position{pos.X + relX*mult, pos.Y + relY*mult}
					val, ok := deck[relPos]
					if ok {
						if val == "L" {
							break
						} else if val == "#" {
							filledNeighbours++
							break
						} else {
							mult++
						}
					} else {
						break
					}
				}
			}
		}
	}
	if filledNeighbours < 4 && version == 1 || filledNeighbours < 5 && version == 2 {
		return false
	}
	return true
}

func main() {
	// Reader
	path := "./input.txt"
	buf, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(buf)

	// Interpret input
	// Record the ferry deck as a map of Positions to states
	var input []string
	for scanner.Scan() {
		input = append(input, scanner.Text())
	}
	var deck = make(map[Position]string)
	for i := range input {
		lineSplit := strings.Split(input[i], "")
		for j := range lineSplit {
			pos := Position{i, j}
			deck[pos] = lineSplit[j]
		}
	}

	// P1: Modifying the deck object directly would break checks against surrounding zones, so unfortunately a new map is required each time
	// For each spot, check:
	//   - If state is . then ignore it
	//   - If state is L then check the eight (or however many) around it for # - if no #, then change L to #
	//   - If state is # then check the eight (or however many) around it for # - if >= 4 #, then change # to L
	// Upon any state change, flag - if a pass contains no state changes, it is done; return the final deck state
	resolvedDeck := ResolveDeck(deck, 1)
	log.Println("P1 | FILLED SEATS:", PrintDeck(resolvedDeck, false))

	// P2: Rules are updated:
	//   - If state is . then ignore it
	//   - If state is L then check the eight directions for non-floor; if no #, transition the L to #
	//   - If state is # then check the eight directions for non-floor; if >= 5 #, transition the # to L
	resolvedDeckV2 := ResolveDeck(deck, 2)
	log.Println("P2 | FILLED SEATS:", PrintDeck(resolvedDeckV2, false))
}
