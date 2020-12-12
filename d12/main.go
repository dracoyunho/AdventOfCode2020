package main

import (
	"bufio"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
)

// Ship struct represents the basic parts of the ship - its heading/bearing and X and Y positions
type Ship struct {
	Heading int // degrees, where 0 is North, increasing clockwise
	X       int // +ve - to the East, -ve - to the West
	Y       int // +ve - to the North, -ve - to the South
}

// RunInstructionsV1 takes a list of instructions, expecting a single letter followed by digits per line, and an initial Ship state
// It returns a Ship struct representing the final state
func RunInstructionsV1(input []string, ship Ship) Ship {
	if ship.Heading < 0 {
		ship.Heading += 360
	}

	reInst := regexp.MustCompile(`(?P<Direction>\w)(?P<Value>\d+)`)
	for _, line := range input {
		instruction := reInst.FindStringSubmatch(line)

		if instruction[1] == "F" {
			switch ship.Heading {
			case 0:
				instruction[1] = "N"
			case 90:
				instruction[1] = "E"
			case 180:
				instruction[1] = "S"
			case 270:
				instruction[1] = "W"
			}
		}

		val, err := strconv.Atoi(instruction[2])
		if err != nil {
			log.Fatal(err)
		}

		switch instruction[1] {
		case "E":
			ship.X += val
		case "W":
			ship.X -= val
		case "N":
			ship.Y += val
		case "S":
			ship.Y -= val
		case "R":
			ship.Heading = (ship.Heading + val) % 360
		case "L":
			ship.Heading = (ship.Heading - val) % 360
		}

		if ship.Heading < 0 {
			ship.Heading += 360
		}
	}
	return ship
}

// RunInstructionsV2 takes a list of instructions, expecting a single letter followed by digits per line, an initial waypoint (represented by a Ship struct), and an initial ship state
// It returns two Ship structs representing the final state of the waypoint and the ship itself
func RunInstructionsV2(input []string, ship, waypoint Ship) (Ship, Ship) {
	reInst := regexp.MustCompile(`(?P<Direction>\w)(?P<Value>\d+)`)
	for _, line := range input {
		instruction := reInst.FindStringSubmatch(line)
		val, err := strconv.Atoi(instruction[2])
		if err != nil {
			log.Fatal(err)
		}
		switch instruction[1] {
		case "E":
			waypoint.X += val
		case "W":
			waypoint.X -= val
		case "N":
			waypoint.Y += val
		case "S":
			waypoint.Y -= val
		case "R":
			// A positive rotation swaps the x and y values, and multiplies the new y value by -1
			// This can be done without new variables: 1) x = x+y , 2) y = y-x , 3) x = x+y
			// Do this as many times as needed (based on val)
			for i := 0; i < val/90; i++ {
				waypoint.X += waypoint.Y
				waypoint.Y -= waypoint.X
				waypoint.X += waypoint.Y
			}
		case "L":
			// A negative rotation swaps the x and y values, and multiplies the new x value by -1
			// This can be done without new variables: 1) y = y+x , 2) x = x-y , 3) y = y+x
			// Do this as many times as needed (based on val)
			for i := 0; i < val/90; i++ {
				waypoint.Y += waypoint.X
				waypoint.X -= waypoint.Y
				waypoint.Y += waypoint.X
			}
		case "F":
			// Simply multiply the waypoint's Cartesian coords - this tells how much to add to the ship's Cartesian coords
			ship.X += waypoint.X * val
			ship.Y += waypoint.Y * val
		}
	}
	return waypoint, ship
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
	var input []string
	for scanner.Scan() {
		input = append(input, scanner.Text())
	}

	// The initial Ship state - 0,0 and heading East
	initial := Ship{90, 0, 0}

	// P1: Simply execute the instructions - the two positions will be in the returned struct
	final := RunInstructionsV1(input, initial)
	log.Println("P1 | MANHATTAN DISTANCE:", math.Abs(float64(final.X))+math.Abs(float64(final.Y)))

	// P2: The waypoint is basically a ghost ship
	// The X and Y coords for this ghost ship merely designate where the marker is relative to the real ship
	// These Cartesian coords don't change even if the ship moves
	// The Heading is no longer necessary but I'm lazy lol
	initial = Ship{0, 0, 0}
	initialWaypoint := Ship{0, 10, 1}
	_, final = RunInstructionsV2(input, initial, initialWaypoint)
	log.Println("P2 | MANHATTAN DISTANCE:", math.Abs(float64(final.X))+math.Abs(float64(final.Y)))
}
