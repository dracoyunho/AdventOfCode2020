package main

import (
	"bufio"
	"log"
	"math"
	"os"
	"strings"
)

const (
	// InputFilePath is the path to the input for this puzzle
	InputFilePath string = "./input.txt"
)

// A HexVec uses a hexagonal basis defined such that, from hexagon centre-to-centre, (h, k) represents two vectors with h rotated 0 rad and k rotated +π/6 rad. Details below:
// A hexagonal space may be represented succinctly with just two basis vectors:
// Consider a line drawn through two opposite vertices of a hexagon (i.e. its major diagonal) such that this line is horizontal
// Consider then on the right side of the hexagon two additional hexagons, overlapping on the right upper and right lower edges
// Then draw a vector from the centre of the reference hexagon to the centre of each hexagon to the right; this now forms a complete basis, usually marked as (h, k)
// Reference: https://www.researchgate.net/profile/Joshua_Island/publication/320407497/figure/fig2/AS:551016690970625@1508384024782/a-Real-space-graphene-lattice-with-primitive-vectors-a-1-and-a-2-The-lattice-has-a.png
// This also works if instead of being horizontal, the major diagonal is vertical, and the overlapping hexagons are directly to the right and to the upper right
// Reference: https://www.researchgate.net/publication/324477876/figure/fig1/AS:614570607538178@1523536459336/Two-basis-vectors-a-1-a-2-and-one-generating-vector-A-for-a-hexagonal-lattice-In.png
// While it is also possible to consider this basis as originating from a vertex of the hexagon and pointed at its minor diagonals, it's not individual vertices that matter in this puzzle but entire hexagons,
//   so it's easier to use the centre-of-hexagon definition
// Note that the input to this puzzle designates east, west, north(west/east), and south(west/east)
// Using an (h, k) where h represents hexagons to the east and k represents hexagons to the northeast results in a relatively simple coordinate system:
// E:  (  1,  0)
// W:  ( -1,  0)
// NE: (  0,  1)
// SW: (  0, -1)
// NW: ( -1,  1)
// SE: (  1, -1)
// In this way, any hexagonal tile may be represented as a vector from (0, 0) to (h, k), and moving from tile to tile is merely linear algebra
// I love it when my nanotechnology engineering degree is actually useful
type HexVec struct {
	H, K int
}

// Bounds discovers the max |h| and |k| spanned by a given set of tiles, so that the tile set is guaranteed to be enclosed in the grid, and adds the specified empty space padding
func Bounds(tiles map[HexVec]struct{}, padding int) (int, int) {
	var h, k int = 0, 0
	for tile := range tiles {
		if h < int(math.Abs(float64(tile.H))) {
			h = int(math.Abs(float64(tile.H)))
		}
		if k < int(math.Abs(float64(tile.K))) {
			k = int(math.Abs(float64(tile.K)))
		}
	}
	return h + padding, k + padding
}

// ActiveNeighbours ingests an active tile set and a specific tile, and checks how many of the six neighbours around the tile are in the active tile set, returning the number of active neighbours
// The requested tile does not need to be part of the active tile set itself
func ActiveNeighbours(tiles map[HexVec]struct{}, tile HexVec) int {
	var neighbours int = 0

	for dh := -1; dh <= 1; dh++ {
		for dk := -1; dk <= 1; dk++ {
			// It's not just the case of (dh, dk) = (0, 0) that needs to be skipped; see the six valid directions from the HexVec description
			if dh == 0 && dk == 0 || dh == -1 && dk == -1 || dh == 1 && dk == 1 {
				continue
			}
			if _, def := tiles[HexVec{tile.H + dh, tile.K + dk}]; def {
				neighbours++
			}
		}
	}

	return neighbours
}

// Evolve ingests an active tile set and applies the HexGOL rules to it
func Evolve(tiles map[HexVec]struct{}) map[HexVec]struct{} {
	var newTiles map[HexVec]struct{} = make(map[HexVec]struct{})
	bh, bk := Bounds(tiles, 1)
	// Check every tile from -bh-1 to +bh+1 and -bk-1 to +bk+1
	for h := -1*bh - 1; h <= bh+1; h++ {
		for k := -1*bk - 1; k <= bk+1; k++ {
			neighbours := ActiveNeighbours(tiles, HexVec{h, k})
			if _, def := tiles[HexVec{h, k}]; def && neighbours > 0 && neighbours < 3 {
				newTiles[HexVec{h, k}] = struct{}{}
			} else if _, def := tiles[HexVec{h, k}]; !def && neighbours == 2 {
				newTiles[HexVec{h, k}] = struct{}{}
			}
		}
	}
	return newTiles
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

	// P1: Every line is merely a set of vectors with which to perform a linear combination; the result is a single vector originating from (0, 0),
	//   upon which it may be checked whether or not this tuple is contained within a set of active points, and if not, it is added, and if it is, it is removed
	// Parsing each line may be performed by char - notice how n and s require parsing the next char, i.e. a HexVec is determined only when e and w are encountered
	var defs map[int][]HexVec = make(map[int][]HexVec)
	for i := range input {
		var def []HexVec = make([]HexVec, 0)
		chars := strings.Split(input[i], "")
		north, south := false, false
		for c := range chars {
			switch chars[c] {
			case "n":
				north = true
			case "s":
				south = true
			case "e":
				if north {
					def = append(def, HexVec{0, 1})
				} else if south {
					def = append(def, HexVec{1, -1})
				} else {
					def = append(def, HexVec{1, 0})
				}
				north, south = false, false
			case "w":
				if north {
					def = append(def, HexVec{-1, 1})
				} else if south {
					def = append(def, HexVec{0, -1})
				} else {
					def = append(def, HexVec{-1, 0})
				}
				north, south = false, false
			}
		}
		defs[i] = def
	}
	// Now simply determine the single HexVec resulting from each linear combination, which is produced by vector addition
	// To the list of active tiles, flip state accordingly (i.e. add or delete)
	// var tiles map[int]HexVec = make(map[int]HexVec)
	var tiles map[HexVec]struct{} = make(map[HexVec]struct{})
	for i := range defs {
		tile := HexVec{0, 0}
		for vec := range defs[i] {
			tile = HexVec{tile.H + defs[i][vec].H, tile.K + defs[i][vec].K}
		}
		if _, def := tiles[tile]; def {
			delete(tiles, tile)
		} else {
			tiles[tile] = struct{}{}
		}
	}
	log.Println("P1 | Active tile count:", len(tiles))

	// P2: The tiles after P1 are the initial state for a GOL-like hex grid
	// Every hex tile (h, k) has six neighbours, based on the six directions indicated in the description for HexVec
	// Because GOL has a rule where inactive pixels may be flipped to active, then it is clearly not sufficient to iterate over the tiles elements, as this only serves for the A-to-I rule
	// A rather naive method of accounting for all inactive tiles in the vicinity of the active tiles is by determining the hexagon grid that leaves at minimum a 1-tile border around all active tiles
	// This may be done by checking the absolute value of every active tile and taking the largest |h|+1 and largest |k|+1, and iterating over every such tile in a raster scan
	// For every hypothetical tile not in the list of active tiles, an I-to-A check is done, whereas for those in the list, an A-to-I chneck is done
	days := 100
	for day := 0; day < days; day++ {
		log.Println("P2 | Applying HexGOL to day", day+1, "...")
		tiles = Evolve(tiles)
		log.Println("P2 | Active tiles after Day", day+1, ":", len(tiles))
	}
}
