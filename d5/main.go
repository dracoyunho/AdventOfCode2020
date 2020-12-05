package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
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

	// P1: For each boarding pass, set up a low of 0 and a high of 127. The diff, plus 1, indicates the potential row count.
	// The same can be said for seat column: set up a low of 0 and a high of 7. The diff, plus 1, indicates the potential column count.
	// Halving the row/column count indicates the size of the next slice.

	// In P1, the number of interest is the highest seat ID - which is row * 8 + column.
	// If the plane were full, then the highest potential seat ID could be 127 * 8 + 7, or 1023.
	maxSeatID := 0
	for _, line := range input {
		rowLow := 0
		rowHigh := 127
		columnLow := 0
		columnHigh := 7
		splits := strings.Split(line, "")
		for _, char := range splits {
			rowCount := (rowHigh - rowLow + 1) / 2
			columnCount := (columnHigh - columnLow + 1) / 2
			switch c := char; c {
			case "F":
				rowHigh -= rowCount
			case "B":
				rowLow += rowCount
			case "L":
				columnHigh -= columnCount
			case "R":
				columnLow += columnCount
			}
		}
		if maxSeatID < rowLow*8+columnLow {
			maxSeatID = rowLow*8 + columnLow
		}
	}
	log.Print(fmt.Sprintf("P1 | MAX SEAT ID: %d", maxSeatID))

	// P2: There is a naive way of doing this, which is to gather all the seat IDs, sort them, and then skip along until the next ID is missing
	// Which is probably the easiest way of handling this, tbh.
	var seatIDs []int
	for _, line := range input {
		rowLow := 0
		rowHigh := 127
		columnLow := 0
		columnHigh := 7
		splits := strings.Split(line, "")
		for _, char := range splits {
			rowCount := (rowHigh - rowLow + 1) / 2
			columnCount := (columnHigh - columnLow + 1) / 2
			switch c := char; c {
			case "F":
				rowHigh -= rowCount
			case "B":
				rowLow += rowCount
			case "L":
				columnHigh -= columnCount
			case "R":
				columnLow += columnCount
			}
		}
		seatIDs = append(seatIDs, rowLow*8+columnLow)
	}
	sort.Ints(seatIDs)
	for index, seatID := range seatIDs {
		if seatIDs[index+1]-seatID > 1 {
			log.Print(fmt.Sprintf("P2 | BOOKED SEAT ID: %d", seatID+1))
			break
		}
	}
}
