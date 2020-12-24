package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

const (
	// InputFilePath is the path to the input for this puzzle
	InputFilePath string = "./input.txt"
)

// DrawCard takes an input deck and draws the top card (i.e. the front of the slice), returning the top card's value and the remaining cards
func DrawCard(deck []int) (int, []int) {
	var card int
	var remains []int
	card = deck[0]
	if len(deck) > 1 {
		remains = deck[1:]
	}
	return card, remains
}

// PlayHand performs one iteration of the game given two players' decks
// It returns the state of the two decks after the round, i.e. it will assign the winner the hand's cards
func PlayHand(p1, p2 []int) ([]int, []int) {
	// Yoink the top card from both
	p1c, p1 := DrawCard(p1)
	p2c, p2 := DrawCard(p2)
	// Because the cards are all unique, there is no chance of a draw
	if p1c > p2c {
		log.Println("Player 1:", p1c, "| Player 2:", p2c, "| Winner: Player 1")
		// P1 wins, assign back the P1 card, then the P2 card
		p1 = append(p1, p1c, p2c)
	} else {
		log.Println("Player 1:", p1c, "| Player 2:", p2c, "| Winner: Player 2")
		p2 = append(p2, p2c, p1c)
	}
	return p1, p2
}

// Score calculates the score of a deck
// The score is calculated as the sum of the product of the card value by its distance from the bottom of the deck, where the top card multiplier is the length of the deck and decrements by 1 downward
func Score(deck []int) int {
	var score int = 0
	for i := range deck {
		score += (len(deck) - i) * deck[i]
	}
	return score
}

// PlayGame plays a full game of Combat, right up until there are zero cards in someone's hand, returning the player that won
// In doing so, it will calculate and print the score of the winner of the game
func PlayGame(p1, p2 []int) int {
	var round int = 1
	for len(p1) > 0 && len(p2) > 0 {
		log.Println("==== Round", round, "====")
		p1, p2 = PlayHand(p1, p2)
		round++
	}
	// Calculate the winning score of whoever has all the cards
	if len(p1) != 0 {
		log.Println("P1 | Player 1 Score:", Score(p1))
		return 1
	}
	log.Println("P1 | Player 2 Score:", Score(p2))
	return 2
}

// EquateDecks will take two decks and see if they are the same; decks are equal if the cards in each deck position between the decks are the same
func EquateDecks(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i]-b[i] != 0 {
			// This means the card that the same position between the decks is different and thus the decks are not equal
			return false
		}
	}
	return true
}

// PlayRecursiveGame plays a full game of Recursive Combat, right up until there are zero cards in someone's hand, returning the player that won
// It recurses when both players have at least as many cards remaining in their deck as the value of the card they just drew
// Every game also remembers round history and will assign victory to P1 automatically if the P1 and P2 cards have been seen before
func PlayRecursiveGame(p1, p2 []int, depth int) int {
	var p1History map[int][]int = make(map[int][]int) // Map round number to P1's hand before drawing a card
	var p2History map[int][]int = make(map[int][]int) // Map round number to P2's hand before drawing a card
	var round int = 0
	// Play hands until one of p1 or p2 becomes empty
	for len(p1) > 0 && len(p2) > 0 {
		handWinner := 0
		round++
		// Check history of both players - if the same decks have been seen in the same round before, assign P1 as GAME winner immediately
		var seen int = -1
		for h := range p1History {
			if EquateDecks(p1, p1History[h]) && EquateDecks(p2, p2History[h]) {
				seen = h
				break
			}
		}
		if seen > -1 {
			// Assign P1 as game winner without dealing out cards from the two players' decks
			log.Println("Depth", depth, "Round", round, "| Player 1 won the game by historical basis! This deck set last seen in round", seen)
			break
		}
		// If the two decks have not been seen in this combination before, then it is OK to proceed with a regular hand
		// First add these two decks to history, by mapping each player's deck to the current round
		p1History[round] = p1
		p2History[round] = p2
		// Yoink the top card from both
		p1c, p1r := DrawCard(p1)
		p1 = p1r
		p2c, p2r := DrawCard(p2)
		p2 = p2r
		log.Println("Depth", depth, "Round", round, "| Player 1 card & deck:", p1c, "&", p1, "| Player 2 card & deck:", p2c, "&", p2)
		// If both p1c and p2c are the size of p1 and p2 after drawing, then recurse, using decks that are the size of p1c and p2c
		if p1c <= len(p1) && p2c <= len(p2) {
			// Assemble new decks
			// Why this? Because passing p1[:p1c] and p2[:p2c] resulted in p2's last element being modified to be the first, for whatever wizardry reason
			// And using copy(np1, p1[:p1c]), copy(np2, p2[:p2c]) just passed empty arrays, which was even worse
			var np1, np2 []int
			for i := 0; i < p1c; i++ {
				np1 = append(np1, p1[i])
			}
			for i := 0; i < p2c; i++ {
				np2 = append(np2, p2[i])
			}
			log.Println("Depth", depth, "Round", round, "| Recursing into a new game with Player 1 deck", np1, "and Player 2 deck", np2)
			handWinner = PlayRecursiveGame(np1, np2, depth+1)
		} else {
			if p1c > p2c {
				handWinner = 1
			} else {
				handWinner = 2
			}
		}
		log.Println("Depth", depth, "Round", round, "| Hand won by player", handWinner)
		if handWinner == 1 {
			p1 = append(p1, p1c, p2c)
		} else if handWinner == 2 {
			p2 = append(p2, p2c, p1c)
		} else {
			log.Fatal("Hand winner was not set properly, it is currently:", handWinner, "| Player 1 deck", p1, "| Player 2 deck", p2)
		}
	}
	// Game over
	// The game could also end by historical basis, upon which Player 1 wins, despite not having all cards, so P1 may win as long as it has more than 0 cards
	if len(p1) > 0 {
		log.Println("Game of depth", depth, "won by Player 1; score:", Score(p1))
		return 1
	}
	log.Println("Game of depth", depth, "won by Player 2; score:", Score(p2))
	return 2
}

func main() {
	// Reader
	buf, err := os.Open(InputFilePath)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(buf)

	// Retrieve input
	var input []string
	for scanner.Scan() {
		input = append(input, scanner.Text())
	}

	// Either it's a header (indicating the player) or it's a number, indicating the card, or it's just an empty line, which can be safely ignored
	var deck1, deck2 []int
	var playerAssign int = 0
	for _, line := range input {
		if line == "Player 1:" {
			playerAssign = 1
		} else if line == "Player 2:" {
			playerAssign = 2
		} else if line == "" {
			continue
		} else {
			val, err := strconv.Atoi(line)
			if err != nil {
				log.Fatal(err)
			}
			if playerAssign == 1 {
				deck1 = append(deck1, val)
			} else if playerAssign == 2 {
				deck2 = append(deck2, val)
			} else {
				log.Fatal("Assigned player is not 1 or 2, currently:", playerAssign)
			}
		}
	}

	// P1: Regular Combat
	// PlayGame(deck1, deck2)

	// P2: Recursive Combat
	PlayRecursiveGame(deck1, deck2, 0)
}
