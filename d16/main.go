package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// ValidTicketValues holds all ranges of valid values for each field on a train ticket
// Its properties should have the same name as a Ticket, but each property is a map of 2-slices
// Each 2-slice contains the start and end values of the valid value range (there may be multiple, hence the map)
// type ValidTicketValues struct {
// 	DepartureLocation map[int][2]int
// 	DepartureStation  map[int][2]int
// 	DeparturePlatform map[int][2]int
// 	DepartureTrack    map[int][2]int
// 	DepartureDate     map[int][2]int
// 	DepartureTime     map[int][2]int
// 	ArrivalLocation   map[int][2]int
// 	ArrivalStation    map[int][2]int
// 	ArrivalPlatform   map[int][2]int
// 	ArrivalTrack      map[int][2]int
// 	Class             map[int][2]int
// 	Duration          map[int][2]int
// 	Price             map[int][2]int
// 	Route             map[int][2]int
// 	Row               map[int][2]int
// 	Seat              map[int][2]int
// 	Train             map[int][2]int
// 	Type              map[int][2]int
// 	Wagon             map[int][2]int
// 	Zone              map[int][2]int
// }

// Ticket holds the properties of a single ticket, an int per value
type Ticket struct {
	DepartureLocation int
	DepartureStation  int
	DeparturePlatform int
	DepartureTrack    int
	DepartureDate     int
	DepartureTime     int
	ArrivalLocation   int
	ArrivalStation    int
	ArrivalPlatform   int
	ArrivalTrack      int
	Class             int
	Duration          int
	Price             int
	Route             int
	Row               int
	Seat              int
	Train             int
	Type              int
	Wagon             int
	Zone              int
}

// InvalidTicketFields ingests some tickets (mapped by numeric IDs to a slice of fields) and a set of valid ticket values, mapped as
//   field names to ranges (where every range is a map of range start to range end)
// The return value is a mapping of ticket IDs to a list of its invalid field indexes
// If a ticket is valid (i.e. all of its field values are in some range, somewhere) then it is not in the returned mapping.
func InvalidTicketFields(validTicketValues map[string]map[int]int, tickets map[int][]int) map[int][]int {
	checkFails := make(map[int][]int)
	// For every ticket, record the index that failed a check - and if the ticket's values passes every check, don't record it
	for ticketNumber, ticketFields := range tickets {
		var badIndexes []int
		for fieldIndex := range ticketFields {
			pass := false
			for _, validRange := range validTicketValues {
				for start, end := range validRange {
					if ticketFields[fieldIndex] >= start && ticketFields[fieldIndex] <= end {
						pass = true
					}
				}
			}
			// Only check for pass down here, since if the field value on the ticket is in ANY range, it is OK
			if !pass {
				// Append to the list of bad field indexes for this ticket
				badIndexes = append(badIndexes, fieldIndex)
			}
		}
		// Only add the ticket's ID here if it is invalid in some way
		if badIndexes != nil {
			checkFails[ticketNumber] = badIndexes
		}
	}
	return checkFails
}

// ScanErrorRate will return the sum of all invalid values, given a map of ticket IDs to their values, as well as a map of ticket IDs to invalid field indexes for each ticket
func ScanErrorRate(tickets map[int][]int, invalidFields map[int][]int) int {
	scanErrorRate := 0
	if len(invalidFields) == 0 {
		return scanErrorRate
	}
	for ticketNumber, invalidFieldIndexes := range invalidFields {
		for _, index := range invalidFieldIndexes {
			scanErrorRate += tickets[ticketNumber][index]
		}
	}
	return scanErrorRate
}

// MapFieldIndexToNames ingests a map of field names to their ranges and a set of reference tickets; from the reference tickets, it will attempt to map field names to indexes
func MapFieldIndexToNames(validTicketValues map[string]map[int]int, referenceTickets map[int][]int) map[string]int {
	fieldFailSets := make(map[string]map[int]struct{}) // The value map grows with keys being field indexes that didn't pass checks - once any one reaches size 19, then the field name is mappable
	fieldMapping := make(map[string]int)

	// See if it's possible to map all 20 in one pass
	for fieldName, fieldRanges := range validTicketValues {
		// log.Println("DEBUG | MAPPING", fieldName)
		var failSet = make(map[int]struct{}) // Holds unique values of failed indexes for the given field name
		// It's possible to reuse InvalidTicketFields to perform the check (any returned indexes indicate any field indexes that won't work for the given field name)
		var testFieldName = map[string]map[int]int{fieldName: fieldRanges}
		if checkFails := InvalidTicketFields(testFieldName, referenceTickets); checkFails != nil {
			// The keys of checkFails don't matter - they're just ticket IDs - but the values do; the values of each map key in checkFails indicates the field indexes that don't work for the given field name
			// This works for the given field name because the call to InvalidTicketFields only has one field name at a time
			for _, failedIndexes := range checkFails {
				for _, i := range failedIndexes {
					failSet[i] = struct{}{}
				}
			}
		} // If there is no failure, that doesn't guarantee that the field name maps to this index
		// log.Println("DEBUG | FAIL SET FOR", fieldName, ":", failSet)
		if len(failSet) == len(validTicketValues)-1 {
			// At this point there are enough fails to demonstrate that this field name is mappable to an index, which is the one index that isn't present
			for i := 0; i < len(validTicketValues); i++ {
				if _, found := failSet[i]; !found {
					fieldMapping[fieldName] = i
				}
			}
		}
		fieldFailSets[fieldName] = failSet
	}

	// Set up process of elimination, if needed
	// This may be accomplished by dropping field names in fieldMapping from fieldFailSets and then adding already-used field indexes to every fail set in fieldFailSets
	for len(fieldMapping) < len(validTicketValues) {
		for mappedFieldName, mappedFieldIndex := range fieldMapping {
			delete(fieldFailSets, mappedFieldName)
			for unmappedFieldName, failSet := range fieldFailSets {
				failSet[mappedFieldIndex] = struct{}{}
				fieldFailSets[unmappedFieldName] = failSet
				if len(failSet) == len(validTicketValues)-1 {
					// At this point there are enough fails to demonstrate that this field name is mappable to an index, which is the one index that isn't present
					for i := 0; i < len(validTicketValues); i++ {
						if _, found := failSet[i]; !found {
							fieldMapping[unmappedFieldName] = i
						}
					}
				}
			}
		}
	}

	return fieldMapping
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
	validTicketValues := make(map[string]map[int]int)
	var personalTicketValues []int
	referenceTicketValues := make(map[int][]int)
	phase := 0
	referenceTicketCount := 0
	// personalTicket := make(map[string]int)
	// var referenceTickets []map[string]int
	for _, line := range input {
		if line == "" {
			phase++
			continue
		}

		if phase == 0 {
			// Each valid value line is explicitly called out as "value:", and values take the form of int ranges defined with a hyphen and separated by " or "
			// The personal ticket is preceded by the line "your ticket:" whereas the reference tickets are preceded by the line "nearby tickets:"
			lineHeader := strings.Split(line, ":")
			validRanges := strings.Split(strings.TrimSpace(lineHeader[1]), " or ")
			validTicketValues[lineHeader[0]] = make(map[int]int)
			for _, validRange := range validRanges {
				rangeValues := strings.Split(validRange, "-")
				start, err := strconv.Atoi(rangeValues[0])
				if err != nil {
					log.Fatal(err)
				}
				end, err := strconv.Atoi(rangeValues[1])
				if err != nil {
					log.Fatal(err)
				}
				validTicketValues[lineHeader[0]][start] = end
			}
			log.Println("Set valid ranges for", lineHeader[0], "as:", validTicketValues[lineHeader[0]])
		} else if phase == 1 {
			// Ignore the your ticket: header line
			// The next line is guaranteed the personal ticket
			if header, err := regexp.MatchString(`your ticket`, line); !header && err == nil {
				// All values are separated by CSV
				for _, split := range strings.Split(line, ",") {
					val, err := strconv.Atoi(split)
					if err != nil {
						log.Fatal(err)
					}
					personalTicketValues = append(personalTicketValues, val)
				}
				log.Println("Discovered personal ticket values:", personalTicketValues)
			} else if err != nil {
				log.Fatal(err)
			}
		} else if phase == 2 {
			// Ignore the nearby tickets: header line
			// Every non-empty line following is a reference ticket
			if header, err := regexp.MatchString(`nearby tickets`, line); !header && err == nil {
				// All values are separated by CSV
				var referenceTicket []int
				for _, split := range strings.Split(line, ",") {
					val, err := strconv.Atoi(split)
					if err != nil {
						log.Fatal(err)
					}
					referenceTicket = append(referenceTicket, val)
				}
				referenceTicketValues[referenceTicketCount] = referenceTicket
				referenceTicketCount++
			} else if err != nil {
				log.Fatal(err)
			}
		}
	}
	// log.Println("Discovered reference ticket values:", referenceTicketValues)

	// P1: Parse through each reference ticket, and check if any given value on a ticket fails to meet any of the valid ranges
	invalidReferenceFields := InvalidTicketFields(validTicketValues, referenceTicketValues)
	log.Println("P1 | INVALID REFERENCE TICKET FIELD INDEXES:", invalidReferenceFields)
	log.Println("P1 | SCAN ERROR RATE:", ScanErrorRate(referenceTicketValues, invalidReferenceFields))

	// Use the ticket IDs in the invalidReferenceFields mapping to toss out invalid tickets from the reference group
	if len(invalidReferenceFields) > 0 {
		for ticketNumber := range invalidReferenceFields {
			delete(referenceTicketValues, ticketNumber)
		}
	}

	// P2: Instead of comparing all valid ranges against one ticket at a time, compare one set of valid ranges against one field value from personal + reference tickets
	// While any set of field values sliced across tickets may pass more than one range check, there cannot be a situation where two sets of field values do this
	// However, if a set of field values sliced across tickets passes N range checks, there could be a situation where another set passes N-1 checks, and N does not need to be 2
	// However, at some point, one field value slice passes only set of range checks, and from there, a process of elimination may begin
	// The result may then be stored as a map of field names to indexes
	// Because the personal ticket may also be used to determine field index-to-name mappings, put it into the reference ticket map as ID -1
	referenceTicketValues[-1] = personalTicketValues
	fieldMapping := MapFieldIndexToNames(validTicketValues, referenceTicketValues)
	log.Println("P2 | FIELD NAME - INDEX MAPPING:", fieldMapping)
	fieldProduct := 1
	for fieldName, fieldIndex := range fieldMapping {
		match, err := regexp.MatchString(`^departure`, fieldName)
		if err != nil {
			log.Fatal(err)
		}
		if match {
			log.Println("P2 | FIELD:", strings.ToTitle(fieldName), "| VALUE:", personalTicketValues[fieldIndex])
			fieldProduct *= personalTicketValues[fieldIndex]
		}
	}
	log.Println("P2 | DEPARTURE FIELD PRODUCT:", fieldProduct)
}
