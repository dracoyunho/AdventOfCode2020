package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

// Eea performs the extended Euclidean algorithm, solving Bezout's lemma ax + by = 1, given some a and b, producing Bezout coefficients x and y
// After returning x and y, it also returns the GCD
// If positive-only values are required, submit flag
func Eea(a, b int, positive bool) (int, int, int) {
	if a < b {
		return Eea(b, a, positive)
	}
	rem1, rem2 := a, b
	s1, s2, t1, t2 := 1, 0, 0, 1
	for rem2 != 0 {
		q := rem1 / rem2
		rem1, rem2 = rem2, rem1-q*rem2 // Enforce rem1 > rem2
		s1, s2 = s2, s1-q*s2
		t1, t2 = t2, t1-q*t2
	}
	return s1, t1, rem1
}

// Gcd calculates the GCD of two numbers
func Gcd(a, b int) int {
	if a < b {
		return Gcd(b, a)
	}
	if b == 0 {
		return a
	}
	return Gcd(b, a%b)
}

// Egcd implements extended Euclidean algorithm while returning the GCD of two values
// The GCD is the first returned value
// Given the modular multiplicative inverse problem ax ≡ 1 (mod m), the inverse x is the second returned value
// The third returned value is the MMI for my ≡ 1 (mod a), i.e. y
func Egcd(a, b int) (int, int, int) {
	x, y, u, v := 0, 1, 1, 0
	for a != 0 {
		q, r := b/a, b%a
		b, a, x, y, u, v = a, r, u, v, x-u*q, y-v*q
	}
	return b, x, y
}

// ModInv discovers the modular multiplicative inverse given an argument and a modulus, guaranteed positive
// That is, it solve for a^-1 in a*a^1 ≡ 1 (mod m)
// If there was something wrong with the input, the success flag is false
func ModInv(arg, mod int) (int, bool) {
	if mod == 0 {
		return 0, false // undefined; how can something be mod 0?
	}
	if arg == 0 {
		return 0, true // trivial case
	}
	gcd, x, _ := Egcd(arg, mod)
	if gcd != 1 {
		return gcd, false
	}
	// x % mod reduces the magnitude of x as much as possible but it won't stop it from being negative
	// Instead return mod + x if x < 0
	x %= mod
	if x < 0 {
		return mod + x, true
	}
	return x % mod, true
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

	// The first line is the timestamp - it's non-conformant to any timestamp standard; it's just some int representing a "time"
	earliestTime, err := strconv.Atoi(input[0])
	if err != nil {
		log.Fatal(err)
	}

	// The second line is comma-separated, representing either numbers or the letter x
	// Thus, if strconv.Atoi fails on it, just go to the next one
	splits := strings.Split(input[1], ",")
	buses := make(map[int]int)
	for delay, val := range splits {
		busID, err := strconv.Atoi(val)
		if err != nil {
			continue
		}
		buses[busID] = delay
	}

	if len(buses) == 0 {
		log.Fatal("No buses are in service!")
	}

	// P1: Time to wait for a bus is the first time to appear after the potential departure time
	// This may be calculated as earliestTime + freq - (earliestTime % freq)
	// This effectively adds a full frequency to the earliest time, which goes past the next bus in the schedule, and then rewinds by the overshoot
	takeBus := 0
	takeWaitTime := 0
	for freq := range buses {
		proposedWaitTime := freq - (earliestTime % freq)
		log.Println("Proposed bus:", freq, "| Wait time:", proposedWaitTime)
		if takeWaitTime == 0 || takeWaitTime > proposedWaitTime {
			takeBus = freq
			takeWaitTime = proposedWaitTime
		}
	}
	log.Println("P1 | TAKE BUS:", takeBus, "| WAIT:", takeWaitTime, "| PRODUCT:", takeBus*takeWaitTime)

	// P2: Consider the first bus to depart with a departure time firstDeparture: every bus after it departs with some delay after the first departure time; call this delay[i]
	// The delay[i] may then be related to the first departure as firstDeparture + delay[i] ≡ 0 (mod bus[i])
	// Alternatively, this may be expressed as firstDeparture ≡ bus[i]-delay[i] (mod bus[i])
	// If the bus frequencies are presumed to be coprime, then Bezout's Lemma applies
	// Given two buses i and j, they satisfy: firstDeparture = (bus[i]-delay[i])*bus[j]*bezout[j] + (bus[j]-delay[j])*bus[i]*bezout[i]
	// Where bezout[i] and bezout[j] are integer coefficients determined by the extended Euclidean algorithm
	// This generalizes very easily when all the bus frequencies are coprime. For each modulus bus[i]:
	//   1) Multiply all bus frequencies except bus[i], yielding a value notBus[i]
	//   2) Determine the Bezout coefficient for notBus[i], i.e. solve b[i]bus[i] + Bezout[i]notBus[i] = 1 for Bezout[i], which should be > 0
	//      Since the Bezout coeff for bus[i] doesn't actually matter, only Bezout[i] needs to be solved for
	//      Rearrangement yields Bezout[i]notBus[i] = 1 - b[i]bus[i] - and conveniently, the RS can simplify with modulus bus[i]
	//      This then results in Bezout[i]notBus[i] ≡ 1 (mod bus[i])
	//   3) Repeat for every bus[i]
	//   4) The solution is the sum of all (bus[i]-delay[i])*Bezout[i]*notBus[i]
	// However, there are a few gotchas in this problem:
	//   1) Bezout[i] is not allowed to be negative, which standard EEA can yield - so this needs to be corrected; it's as easy as doing Bezout[i] % bus[i]
	//      If this is still negative, then just take bus[i] - Bezout[i]
	//   2) If the index position in splits for a given bus[i] is less than the bus[i] itself, it would yield a delay larger than the bus[i]
	//      While this is a valid congruency, it's not valid when it comes to the timestamp
	//      Therefore, the smallest delay[i] (the most immediate bus) is the smaller of index[i] % bus[i] or delay[i], unless delay[i] is already 0, in which case, this is valid

	allBus := 1
	for bus := range buses {
		allBus *= bus
	}
	log.Println("P2 | ALL BUS PRODUCT:", allBus)
	solution := 0
	for bus, delay := range buses {
		delay = delay % bus
		notBus := allBus / bus
		bezoutNotBus, _ := ModInv(notBus, bus)
		solComponent := (bus - delay) * bezoutNotBus * notBus
		solution += solComponent
		log.Println("P2 | BUS", bus, "| DELAY", delay, "| BEZOUT", bezoutNotBus, "| ADDING", solComponent)
	}
	log.Println("P2 | Solution:", solution%allBus)
}
