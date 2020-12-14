package main

import (
	"bufio"
	"fmt"
	"log"
	"math/bits"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Puzzle1 literally solves the entirety of the day's puzzle 1
func Puzzle1(input []string) {
	// Form a map of memory locations to their values
	mem := make(map[uint64]uint64)

	// P1: Because the bit mask is not completely defined for every bit, there isn't really a convenient way of doing this without being bit-iterative
	// Interpret every incoming line: if it starts with mask, then everything after "mask = " is the bitmask
	// If it instead starts with mem, then interpret what is in [] and assign a holding variable to the number after the equals sign
	reMask := regexp.MustCompile(`mask = (?P<Mask>\w+)`)
	reMemAssign := regexp.MustCompile(`mem\[(?P<Address>\w+)\] = (?P<Value>\w+)`)
	var mask string
	for _, line := range input {
		if matched, _ := regexp.MatchString(`^mask = `, line); matched {
			maskSplit := reMask.FindStringSubmatch(line)
			if len(maskSplit) == 0 {
				log.Fatal(line, " was not split by reMask properly")
			}
			mask = maskSplit[1]
			// log.Println("P1 | MASK:", mask)
		} else if matched, _ := regexp.MatchString(`^mem`, line); matched {
			memAssignSplit := reMemAssign.FindStringSubmatch(line)
			if len(memAssignSplit) == 0 {
				log.Fatal(line, " was not split by reMemAssign properly")
			}
			address, err := strconv.ParseUint(memAssignSplit[1], 10, 64)
			if err != nil {
				log.Fatal("Could not parse ", memAssignSplit[1], " to uint64; base error: ", err)
			}
			value, err := strconv.ParseUint(memAssignSplit[2], 10, 64)
			if err != nil {
				log.Fatal("Could not parse ", memAssignSplit[2], " to uint64; base error: ", err)
			}
			// log.Println("P1 | ADDRESS:", address, "| INSTRUCT:", value)

			// Just set the address value to the desired value - the mask can be applied after
			mem[address] = value
			// Since data length is guaranteed to be 36 bits, just go through the mask, and if not X, then force the value to be that given bit
			maskBits := strings.Split(mask, "")
			for index := range maskBits {
				// Specific bits can be checked by rotating the value in mem right and then back left the same amount
				// The amount to rotate by is the data length (36 bits) - the current index - 1
				// e.g. examining the MSB requires rotating rightward by 36 - 0 - 1 = 35 bits
				mem[address] = bits.RotateLeft64(mem[address], -1*(len(maskBits)-index-1))
				if maskBits[index] == "1" && mem[address]%2 == 0 {
					mem[address]++
				} else if maskBits[index] == "0" && mem[address]%2 == 1 {
					mem[address]--
				}
				mem[address] = bits.RotateLeft64(mem[address], len(maskBits)-index-1)
			}
			// log.Println("P1 | WROTE:", mem[address], "AT ADDRESS", address)
		} else {
			log.Fatal("Unknown line: ", line)
		}
	}
	// For every element in the memory map, if its value is non-zero, add to sum
	var p1Sum uint64 = 0
	for _, val := range mem {
		if val != 0 {
			p1Sum += val
		}
	}
	log.Println("P1 | Solution:", p1Sum)
}

// Puzzle2 literally solves the entirety of the day's puzzle 2
func Puzzle2(input []string) {
	// Form a map of memory locations to their values
	mem := make(map[uint64]uint64)

	// P1: Because the bit mask is not completely defined for every bit, there isn't really a convenient way of doing this without being bit-iterative
	// Interpret every incoming line: if it starts with mask, then everything after "mask = " is the bitmask
	// If it instead starts with mem, then interpret what is in [] and assign a holding variable to the number after the equals sign
	reMask := regexp.MustCompile(`mask = (?P<Mask>\w+)`)
	reMemAssign := regexp.MustCompile(`mem\[(?P<Address>\w+)\] = (?P<Value>\w+)`)
	var mask string
	for _, line := range input {
		if matched, _ := regexp.MatchString(`^mask = `, line); matched {
			maskSplit := reMask.FindStringSubmatch(line)
			if len(maskSplit) == 0 {
				log.Fatal(line, " was not split by reMask properly")
			}
			mask = maskSplit[1]
		} else if matched, _ := regexp.MatchString(`^mem`, line); matched {
			memAssignSplit := reMemAssign.FindStringSubmatch(line)
			if len(memAssignSplit) == 0 {
				log.Fatal(line, " was not split by reMemAssign properly")
			}
			address, err := strconv.ParseUint(memAssignSplit[1], 10, 64)
			if err != nil {
				log.Fatal("Could not parse ", memAssignSplit[1], " to uint64; base error: ", err)
			}
			value, err := strconv.ParseUint(memAssignSplit[2], 10, 64)
			if err != nil {
				log.Fatal("Could not parse ", memAssignSplit[2], " to uint64; base error: ", err)
			}

			// First determine the address mask
			// Priority goes to the mask on bits X and 1
			// The actual address value indicated only needs to be checked when the mask is 0
			// A specific bit may be pulled up by rotating the address value to the right by 36 - index - 1 (e.g. the 3rd MSB by rotating 36 - 2 - 1, the 2nd LSB by rotating 36 - 34 - 1)
			var maskedAddressBits [36]string
			maskBits := strings.Split(mask, "")
			inputAddressBits := strings.Split(fmt.Sprintf("%036b", address), "")
			for index := range maskBits {
				if maskBits[index] != "0" {
					maskedAddressBits[index] = maskBits[index]
				} else {
					maskedAddressBits[index] = inputAddressBits[index]
				}
			}

			// log.Println("P2 | INPUT ADDRESS", fmt.Sprintf("%036b", address), "| INPUT MASK", mask, "| FINAL MASK", strings.Join(maskedAddressBits[0:], ""))

			// Generate a slice of all valid addresses
			// Every such valid address may be generated bit-wise:
			//   Rotate the addresses in the slice left by 1
			//   If the upcoming mask bit is 0, continue
			//   If the upcoming mask bit is 1, add one to every value in the slice and continue
			//   If the upcoming mask bit is X, then for every current address in the slice (they have already been rotated), append a new value that is incremented from the original
			var addressesToWrite []uint64
			addressesToWrite = append(addressesToWrite, 0)
			for bitIndex := range maskedAddressBits {
				currAddressCount := len(addressesToWrite) // Only directly examine the current addresses in this pass and not any new addresses added by an X mask bit
				for i := 0; i < currAddressCount; i++ {
					addressesToWrite[i] = bits.RotateLeft64(addressesToWrite[i], 1)
					if maskedAddressBits[bitIndex] == "X" {
						addressesToWrite = append(addressesToWrite, addressesToWrite[i]+1)
					} else if maskedAddressBits[bitIndex] == "1" {
						addressesToWrite[i]++
					}
				}
			}

			// Now write the value to every address generated
			for _, address := range addressesToWrite {
				mem[address] = value
			}
		} else {
			log.Fatal("Unknown line: ", line)
		}
	}
	// For every element in the memory map, if its value is non-zero, add to sum
	var p2Sum uint64 = 0
	for _, val := range mem {
		if val != 0 {
			p2Sum += val
		}
	}
	log.Println("P2 | Solution:", p2Sum)
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

	Puzzle1(input)

	// P2: ugh
	// While the value being written is more firm, the address is not
	// Every write will require generating a list of valid addresses to write the value to
	// That list of valid addresses is performed by masking the target address with the mask
	// From there, every possible address must be generated
	// The least significant bit value in the mask is now 0 instead of X, and X is the most significant (i.e X overwrites both 1 and 0, 1 overwrites only 0, 0 overwrites nothing)
	Puzzle2(input)
}
