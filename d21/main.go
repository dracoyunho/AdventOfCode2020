package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

const (
	// InputFilePath is the path to the input for this puzzle
	InputFilePath string = "./input.txt"
)

// IntersectLists returns the intersection of two sets (any kind)
func IntersectLists(a, b map[string]struct{}) map[string]struct{} {
	var intersection map[string]struct{} = make(map[string]struct{})

	// To determine the intersection, every element in set A must be in set B (the opposite is true, but not necessary to check)
	for ea := range a {
		if _, ok := b[ea]; ok {
			intersection[ea] = struct{}{}
		}
	}

	return intersection
}

// MatchAllergenToIngredient is given a specific allergen, a set of ingredients lists, and the corresponding allergen lists, and a set of known allergen-to-ingredient mappings, and returns the ingredient matching the allergen
func MatchAllergenToIngredient(allergen string, allergenLists, ingredientLists map[int]map[string]struct{}, knownAllergens map[string]string) string {
	var ingredient string = ""

	// First identify the IDs that have the allergen in them
	var lists []int
	for al := range allergenLists {
		if _, def := allergenLists[al][allergen]; def {
			lists = append(lists, al)
		}
	}

	if len(lists) == 0 {
		log.Println("An attempt was made to match the allergen", allergen, "to an ingredient, but no allergen lists contained", allergen, "!")
		return ingredient
	}

	// A result may only be returned if:
	// 1) There is only one ID, and the list of ingredients is 1 element in size - anything higher is probably not a good thing, as this indicates an unsolvable puzzle
	// 2) There are many IDs that, when intersected with each other, reduces to 1 element
	if len(lists) == 1 {
		if len(ingredientLists[lists[0]]) == 0 {
			log.Println("The allergen list for ID", lists[0], "contained", allergen, ", but the corresponding ingredient list had no ingredients!")
			return ingredient
		}
		for ing := range ingredientLists[lists[0]] {
			// Only declare an ingredient the match for this allergen if it isn't already mapped to some other allergen
			if _, def := knownAllergens[ing]; !def {
				ingredient = ing
				break
			}
		}
	} else {
		// Collapse all the ingredient lists into one list, which after this, should have only one ingredient
		var intersection map[string]struct{} = ingredientLists[lists[0]]
		for i := 1; i < len(lists); i++ {
			intersection = IntersectLists(intersection, ingredientLists[lists[i]])
		}
		if len(intersection) == 0 {
			log.Println("The allergen lists contained", allergen, ", but an intersection search pulled up no matching ingredients!")
			return ingredient
		}
		for ing := range intersection {
			// Only declare an ingredient the match for this allergen if it isn't already mapped to some other allergen
			if _, def := knownAllergens[ing]; !def {
				ingredient = ing
				break
			}
		}
	}

	return ingredient
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

	// Every line of the input is some list of ingredients mapped to some list of allergens
	var ingredientLists map[int]map[string]struct{} = make(map[int]map[string]struct{})
	var allergenLists map[int]map[string]struct{} = make(map[int]map[string]struct{})
	var knownAllergens map[string]string = make(map[string]string)
	for i := range input {
		splits := strings.Split(input[i], " (contains ")
		iL := strings.Split(splits[0], " ")
		aL := strings.Split(strings.Trim(splits[1], ")"), ", ")
		log.Println("Line", i, "| Ingredients:", iL, "| Allergens:", aL)
		ingredientLists[i] = make(map[string]struct{})
		allergenLists[i] = make(map[string]struct{})
		for _, in := range iL {
			ingredientLists[i][in] = struct{}{}
		}
		for _, al := range aL {
			allergenLists[i][al] = struct{}{}
		}
	}

	// The number of unique allergens is determined by looking through the allergen lists for unique values
	var allAllergens map[string]struct{} = make(map[string]struct{})
	for list := range allergenLists {
		for allergen := range allergenLists[list] {
			allAllergens[allergen] = struct{}{}
		}
	}

	// P1: For each allergen, get all ingredient lists that contain the allergen, and intersect the ingredient lists
	// This must necessarily reveal a single ingredient - after all, the allergen is only present in one ingredient
	// Then, both ingredient and allergen may be removed from all lists
	// Check all ingredient lists for a single ingredient value corresponding to a single allergen value - this also indicates a match between ingredient and allergen
	// Repeat until done
	for len(allAllergens) > len(knownAllergens) {
		log.Println("Still searching for allergen matches...")
		log.Println("The complete list of allergens:", allAllergens)
		// Only bother to search for allergens not already in knownAllergens
		for allergen := range allAllergens {
			if _, done := knownAllergens[allergen]; !done {
				log.Println("Attempting to match", allergen, "")
				ingredient := MatchAllergenToIngredient(allergen, allergenLists, ingredientLists, knownAllergens)
				knownAllergens[ingredient] = allergen
				log.Println("Matched", ingredient, "to", allergen)
			}
		}
	}

	// The answer to P1 is the sum of all objects still left behind in all ingredientLists
	p1 := 0
	for i := range ingredientLists {
		for ing := range ingredientLists[i] {
			if _, def := knownAllergens[ing]; !def {
				p1++
			}
		}
	}
	log.Println("P1 | Non-Allergen ingredient incidences:", p1)

	// P2: The known allergens list maps ingredients to allergens, not the other way around. whoops
	// Is it worth building the string with code?
	// tbh, no - it's so short, might as well do it by eye lol
	log.Println("P2 | Ingredients known to contain allergen:", knownAllergens)
}
