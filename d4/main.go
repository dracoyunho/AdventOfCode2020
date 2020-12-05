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

// Passport contains all possible passport fields as strings
type Passport struct {
	byr string
	iyr string
	eyr string
	hgt string
	hcl string
	ecl string
	pid string
	cid string
}

func resetPassport(passport *Passport) {
	passport.byr = ""
	passport.iyr = ""
	passport.eyr = ""
	passport.hgt = ""
	passport.hcl = ""
	passport.ecl = ""
	passport.pid = ""
	passport.cid = ""
}

func main() {
	// P1: If reading by line, then an empty line signifies the end of a record.
	// A Passport struct will be required. This allows us to toss instances out of memory when we don't need them anymore.
	// Every field is known to follow key:value syntax, separated by some kind of whitespace (until the EOL is encountered)
	// Now, a regex can be applied on every line for the expected keys.
	// If a match is returned by FindAllString, then splitting the resulting match on : and taking the second element
	//   results in the desired value.
	// If there is no second element, however, then that means the key was specified, but there is no value.
	// If any field is nil (except cid, currently optional), then the password is invalid.

	// P2: To reiterate:
	// BYR: must be 4 digits, 1920-2002
	// IYR: must be 4 digits, 2010-2020
	// EYR: must be 4 digits, 2020-2030
	// HGT: must be 2/3 digits and then cm or in; 150-193 cm or 59-76 in
	// HCL: must be # followed by 6 hex digits
	// ECL: must be one of amb, blu, brn, gry, grn, hzl, oth
	// PID: must be 9 digits including leading digits
	// CID: Ignored.

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

	// Regexes
	reByrField := regexp.MustCompile("byr:[^\\s]+")
	reIyrField := regexp.MustCompile("iyr:[^\\s]+")
	reEyrField := regexp.MustCompile("eyr:[^\\s]+")
	reHgtField := regexp.MustCompile("hgt:[^\\s]+")
	reHclField := regexp.MustCompile("hcl:[^\\s]+")
	reEclField := regexp.MustCompile("ecl:[^\\s]+")
	rePidField := regexp.MustCompile("pid:[^\\s]+")
	reCidField := regexp.MustCompile("cid:[^\\s]+")

	reByrValue := regexp.MustCompile("^\\d{4}$")
	reIyrValue := regexp.MustCompile("^\\d{4}$")
	reEyrValue := regexp.MustCompile("^\\d{4}$")
	reHgtValue := regexp.MustCompile("^\\d{3}cm$|^\\d{2}in$")
	reHclValue := regexp.MustCompile("^#[0-9a-f]{6}$")
	reEclValue := regexp.MustCompile("^(amb|blu|brn|gry|grn|hzl|oth)$")
	rePidValue := regexp.MustCompile("^\\d{9}$")

	reHgtNum := regexp.MustCompile("\\d{2,3}")

	// Search
	validPassportsP1 := 0
	validPassportsP2 := 0
	var passport Passport
	resetPassport(&passport)
	for _, line := range input {
		if line == "" {
			if passport.byr != "" && passport.ecl != "" && passport.eyr != "" && passport.hcl != "" && passport.hgt != "" && passport.iyr != "" && passport.pid != "" {
				validPassportsP1++

				// log.Print(fmt.Sprintf("P1 | DEBUG | BYR %s | IYR %s | EYR %s | HGT %s | HCL %s | ECL %s | PID %s", passport.byr, passport.iyr, passport.eyr, passport.hgt, passport.hcl, passport.ecl, passport.pid))

				passport.byr = reByrValue.FindString(passport.byr)
				passport.iyr = reIyrValue.FindString(passport.iyr)
				passport.eyr = reEyrValue.FindString(passport.eyr)
				passport.hgt = reHgtValue.FindString(passport.hgt)
				if len(strings.Split(passport.hcl, "")) == 7 {
					passport.hcl = reHclValue.FindString(passport.hcl)
				} else {
					passport.hcl = ""
				}
				if len(strings.Split(passport.ecl, "")) == 3 {
					passport.ecl = reEclValue.FindString(passport.ecl)
				} else {
					passport.ecl = ""
				}
				passport.pid = rePidValue.FindString(passport.pid)

				if passport.byr != "" && passport.ecl != "" && passport.eyr != "" && passport.hcl != "" && passport.hgt != "" && passport.iyr != "" && passport.pid != "" {
					byrValue, err := strconv.Atoi(passport.byr)
					if err != nil {
						log.Fatal(err)
					}
					iyrValue, err := strconv.Atoi(passport.iyr)
					if err != nil {
						log.Fatal(err)
					}
					eyrValue, err := strconv.Atoi(passport.eyr)
					if err != nil {
						log.Fatal(err)
					}
					hgtNum, err := strconv.Atoi(reHgtNum.FindString(passport.hgt))
					if byrValue >= 1920 && byrValue <= 2002 && iyrValue >= 2010 && iyrValue <= 2020 && eyrValue >= 2020 && eyrValue <= 2030 {
						log.Print(fmt.Sprintf("P2 | DEBUG | BYR %s | IYR %s | EYR %s | HGT %s | HCL %s | ECL %s | PID %s", passport.byr, passport.iyr, passport.eyr, passport.hgt, passport.hcl, passport.ecl, passport.pid))

						inCm, err := regexp.MatchString("cm", passport.hgt)
						if err != nil {
							log.Fatal(err)
						}
						inIn, err := regexp.MatchString("in", passport.hgt)
						if err != nil {
							log.Fatal(err)
						}
						if (hgtNum >= 150 && hgtNum <= 193 && inCm) || (hgtNum >= 59 && hgtNum <= 76 && inIn) {
							validPassportsP2++
						}
					}
				}
			}
			resetPassport(&passport)
			continue
		}

		byr := strings.Split(reByrField.FindString(line), ":")
		iyr := strings.Split(reIyrField.FindString(line), ":")
		eyr := strings.Split(reEyrField.FindString(line), ":")
		hgt := strings.Split(reHgtField.FindString(line), ":")
		hcl := strings.Split(reHclField.FindString(line), ":")
		ecl := strings.Split(reEclField.FindString(line), ":")
		pid := strings.Split(rePidField.FindString(line), ":")
		cid := strings.Split(reCidField.FindString(line), ":")

		if len(byr) > 1 {
			passport.byr = byr[1]
		}
		if len(iyr) > 1 {
			passport.iyr = iyr[1]
		}
		if len(eyr) > 1 {
			passport.eyr = eyr[1]
		}
		if len(hgt) > 1 {
			passport.hgt = hgt[1]
		}
		if len(hcl) > 1 {
			passport.hcl = hcl[1]
		}
		if len(ecl) > 1 {
			passport.ecl = ecl[1]
		}
		if len(pid) > 1 {
			passport.pid = pid[1]
		}
		if len(cid) > 1 {
			passport.cid = cid[1]
		}
	}
	log.Print(fmt.Sprintf("P1 | Valid passports: %d", validPassportsP1))
	log.Print(fmt.Sprintf("P2 | Valid passports: %d", validPassportsP2))
}
