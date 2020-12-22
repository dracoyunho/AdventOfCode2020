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

const (
	// TileDim is the square side length of a tile, by observation
	TileDim int = 10
	// ImageDim is the square side length of the whole image of tiles, by observation
	ImageDim int = 12
	// ImagePixelDim is the square side length of the whole image by its individial pixels, by observation
	ImagePixelDim int = ImageDim * (TileDim - 2)
	// InputFilePath is the path to the input for this puzzle
	InputFilePath string = "./input.txt"
)

// Tile represents a whole tile with its painted pixels (or just pixels)
type Tile struct {
	ID     string
	Pixels map[Pixel]struct{}
}

// Pixel represents an active (#) position on a tile
type Pixel struct {
	X int
	Y int
}

// GetTileEdge returns the specified edge of the given tile's pixels, ignoring the tile's rotation
func GetTileEdge(tiles map[string]Tile, id string, edge int) map[Pixel]struct{} {
	// Edge must be between 0 and 3 inclusive
	if edge < 0 || edge > 3 {
		log.Fatal("Requested edge should be between 0 and 3, it was", edge)
	}
	// Because rotation causes indexes to get all subtracted and whatnot, it is actually easier to return the edge based on what is the "free" axis
	// The free axis is the axis parallel to the edge, e.g. the top edge has a free Y axis, so points on the top edge may vary on Y but are fixed to X = 0
	edgePx := make(map[Pixel]struct{})
	for px := range tiles[id].Pixels {
		if edge == 0 && px.X == 0 || edge == 1 && px.Y == TileDim-1 || edge == 2 && px.X == TileDim-1 || edge == 3 && px.Y == 0 {
			edgePx[px] = struct{}{}
		}
	}
	return edgePx
}

// GetPixelsEdge does the same thing as GetTileEdge but accepts a pixel map instead of a whole tile
func GetPixelsEdge(pixels map[Pixel]struct{}, edge int) map[Pixel]struct{} {
	// Edge must be between 0 and 3 inclusive
	if edge < 0 || edge > 3 {
		log.Fatal("Requested edge should be between 0 and 3, it was", edge)
	}
	// Because rotation causes indexes to get all subtracted and whatnot, it is actually easier to return the edge based on what is the "free" axis
	// The free axis is the axis parallel to the edge, e.g. the top edge has a free Y axis, so points on the top edge may vary on Y but are fixed to X = 0
	edgePx := make(map[Pixel]struct{})
	for px := range pixels {
		if edge == 0 && px.X == 0 || edge == 1 && px.Y == TileDim-1 || edge == 2 && px.X == TileDim-1 || edge == 3 && px.Y == 0 {
			edgePx[px] = struct{}{}
		}
	}
	return edgePx
}

// RotatePixels returns the rotated form of the given collection of pixels
// Rotation is the same as a Tile rotation, so that if rotation is 1, what was the right side of the original tile is now pointed upward
// Described in terms of the pixels, every pixel is rotated around the centre-point to the left 90-degrees
func RotatePixels(pixels map[Pixel]struct{}, rotation, dimSize int) map[Pixel]struct{} {
	// In terms of axes, where +X was down and +Y was right, rotation of pixels left is equivalent to rotating axes right
	// From this, +X would then point left and +Y would point down, so that new Y = old X and new X = (TileDim-1) - old Y
	if rotation < 0 {
		// It doesn't really matter if rotation > 3, as the rotation means it'll just be % 4 anyway
		log.Fatal("Rotation should be > 0; it was", rotation)
	}
	for r := rotation; r >= 0; r-- {
		rPixels := make(map[Pixel]struct{})
		for pixel := range pixels {
			rPixels[Pixel{(dimSize - 1) - pixel.Y, pixel.X}] = struct{}{}
		}
		pixels = rPixels
	}
	return pixels
}

// ReflectPixels will reflect a collection of pixels about a described axis, either "x", "y", or "xy"
// Reflection on "xy" just does a reflection on x, then y
func ReflectPixels(pixels map[Pixel]struct{}, axis string, dimSize int) map[Pixel]struct{} {
	rPixels := make(map[Pixel]struct{})
	if len(pixels) == 0 || axis == "" {
		return pixels
	}
	if axis == "x" || axis == "xy" {
		for px := range pixels {
			rPixels[Pixel{dimSize - 1 - px.X, px.Y}] = struct{}{}
		}
	} else if axis == "y" || axis == "xy" {
		for px := range pixels {
			rPixels[Pixel{px.X, dimSize - 1 - px.Y}] = struct{}{}
		}
	} else {
		log.Fatal("Unrecognized axis for reflection:", axis)
	}
	return rPixels
}

// FindCommonEdges ingests a tile set and a reference tile
// The return value is a map of the reference tile's edge to a collection of strings indicating the matching tile, plus the transformations required
// Rotations indicate the number that needed to be passed to RotateTile to get the indicated rotation
// Reflections indicate the axes on which a reflection would need to be performed
// The transformations performed on the matching tile would get it to a state where it may be placed next to the reference edge
// The reference tile is never rotated during this procedure
// If the returned map is empty/nil, then there are no discovered matches
func FindCommonEdges(tiles map[string]Tile, ref string) (map[int]Tile, map[int]map[string]string) {
	// log.Println("Now comparing tile", ref, "to find edge matches...")

	matches := make(map[int]Tile)
	matchTransforms := make(map[int]map[string]string)
	// This double map looks like:
	// <ref edge>:
	//   id: <matching tile id>
	//   reflect: <"" | "x" | "y" | "xy">
	//   rotate: <"0" | "1" | "2" | "3">
	for edge := 0; edge < 4; edge++ {
		rEdge := GetTileEdge(tiles, ref, edge)
		// log.Println("Reference: ID", ref, "| Edge", edge, "| Pixels:", PrintPixels(rEdge, false))

		// Now with this reference edge, attempt to find a match
		// Once a match is found with this candidate, it no longer needs to be checked - it should only have one match to the reference
		for candidate := range tiles {
			if candidate == ref {
				continue
			}
			// It should be noted that for each rEdge, there's really only one possible comparison, which is the opposite edge of the reference
			// For example, if using reference edge 0, the only candidate edge needed is edge 2
			// However, it is possible to rotate and flip the candidates, but it'll still be edge 2
			// By the way - it is possible for edges to be common even if the pixel edge is an empty map!
			axes := []string{"", "x", "y", "xy"}
			for _, axis := range axes {
				for rotation := 0; rotation < 4; rotation++ {
					cEdge := GetPixelsEdge(RotatePixels(ReflectPixels(tiles[candidate].Pixels, axis, TileDim), rotation, TileDim), (edge+2)%4)
					if MatchEdges(rEdge, cEdge, edge, (edge+2)%4) {
						// log.Println("Match found:", "ID", candidate, "| Rotation", rotation, "| Pixels:", PrintPixels(cEdge, false))
						matches[edge] = Tile{candidate, RotatePixels(ReflectPixels(tiles[candidate].Pixels, axis, TileDim), rotation, TileDim)}
						matchTransforms[edge] = make(map[string]string)
						matchTransforms[edge]["id"] = tiles[candidate].ID
						matchTransforms[edge]["reflect"] = axis
						matchTransforms[edge]["rotate"] = fmt.Sprint(rotation)
						break
					}
				}
				if _, matched := matchTransforms[edge]; matched {
					break
				}
			}
			if _, matched := matchTransforms[edge]; matched {
				break
			}
		}
	}

	return matches, matchTransforms
}

// MatchEdges will return wehether or not two Pixel maps would match
// It is assigned two pixel maps and two edge ints, as it is possible for an edge to look like another
func MatchEdges(a, b map[Pixel]struct{}, ea, eb int) bool {
	// For edges to match, they should have the same number of points - if not even this can be satisfied, clearly they do not match
	if len(a) != len(b) {
		return false
	}
	// If a and b are empty maps, well, technically they do match
	if len(a) == 0 && len(b) == 0 {
		return true
	}
	// If these edges are not opposite each other, return false
	if (ea+2)%4 != eb {
		return false
	}
	// Now check free axis - Y is free for edges 0 and 2, X is free for edges 1 and 3
	if ea == 0 {
		for y := 0; y < TileDim; y++ {
			// If there is ever a situation where a Y value is found on one edge but not the other, return false
			_, fa := a[Pixel{0, y}]
			_, fb := b[Pixel{TileDim - 1, y}]
			if fa != fb {
				return false
			}
		}
	} else if ea == 2 {
		for y := 0; y < TileDim; y++ {
			// If there is ever a situation where a Y value is found on one edge but not the other, return false
			_, fa := a[Pixel{TileDim - 1, y}]
			_, fb := b[Pixel{0, y}]
			if fa != fb {
				return false
			}
		}
	} else if ea == 1 {
		for x := 0; x < TileDim; x++ {
			// If there is ever a situation where an X value is found on one edge but not the other, return false
			_, fa := a[Pixel{x, TileDim - 1}]
			_, fb := b[Pixel{x, 0}]
			if fa != fb {
				return false
			}
		}
	} else if ea == 3 {
		for x := 0; x < TileDim; x++ {
			// If there is ever a situation where an X value is found on one edge but not the other, return false
			_, fa := a[Pixel{x, 0}]
			_, fb := b[Pixel{x, TileDim - 1}]
			if fa != fb {
				return false
			}
		}
	}
	return true
}

// EdgeIndex returns the edge index represented by the given pixels
// If it is an indeterminate edge, it returns -1, not 0, as 0 is a valid edge index
func EdgeIndex(pixels map[Pixel]struct{}) int {
	// If the given pixel map is empty, return -1
	if len(pixels) == 0 {
		return -1
	}

	edges := map[int]struct{}{0: struct{}{}, 1: struct{}{}, 2: struct{}{}, 3: struct{}{}}
	// Once a point causes an edge check to fail, yoink it from the map
	for pixel := range pixels {
		if pixel.X != 0 { // Can't be up edge
			delete(edges, 0)
		} else if pixel.Y != 9 { // Can't be right edge
			delete(edges, 1)
		} else if pixel.X != 9 { // Can't be down edge
			delete(edges, 2)
		} else if pixel.Y != 0 { // Can't be left edge
			delete(edges, 3)
		}
	}
	if len(edges) != 1 {
		log.Println("Trying to determine edge index of", PrintPixels(pixels, false), "resulted in these valid edge indexes:", edges)
		return -1
	}
	ei := -1
	for i := range edges {
		ei = i
	}
	return ei
}

// PrintTiles nicely prints all tiles passed to it
func PrintTiles(tiles map[string]Tile) {
	for id := range tiles {
		PrintTile(tiles, id)
	}
}

// PrintTile nicely prints a single tile
func PrintTile(tiles map[string]Tile, id string) {
	log.Println("Tile", tiles[id].ID, ":")
	for x := 0; x < TileDim; x++ {
		var chars []string
		for y := 0; y < TileDim; y++ {
			if _, ok := tiles[id].Pixels[Pixel{x, y}]; ok {
				chars = append(chars, "#")
			} else {
				chars = append(chars, ".")
			}
		}
		log.Println(strings.Join(chars, " "))
	}
	log.Println("")
}

// PrintPixels nicely prints a collection of pixels instead of relying on the default map output
// It also returns what it would have printed, and will stay hushed if all that is desired is the return value
func PrintPixels(pixels map[Pixel]struct{}, verbose bool) string {
	// Assemble a string of (x,y); (x,y); ... and print that
	var builder []string
	var output string
	for px := range pixels {
		builder = append(builder, "("+fmt.Sprint(px.X)+","+fmt.Sprint(px.Y)+")")
	}
	output = strings.Join(builder, "; ")
	if verbose {
		log.Println(output)
	}
	return output
}

// PrintImageTileIDs nicely prints an image's tile IDs
func PrintImageTileIDs(image map[Pixel]Tile) {
	log.Println("Image Tile IDs:")
	for ix := 0; ix < ImageDim; ix++ {
		var rowIDs []string
		for iy := 0; iy < ImageDim; iy++ {
			rowIDs = append(rowIDs, image[Pixel{ix, iy}].ID)
		}
		log.Println(rowIDs)
	}
}

// CornerIDProduct returns the product of the Tile IDs in the four corners of an image
func CornerIDProduct(image map[Pixel]Tile) int {
	product := 1
	ic := []int{0, ImageDim - 1}
	for _, ix := range ic {
		for _, iy := range ic {
			id, err := strconv.Atoi(image[Pixel{ix, iy}].ID)
			if err != nil {
				log.Fatal(err)
			}
			product *= id
		}
	}
	return product
}

// PrintImage nicely prints the Pixels from all of its tiles as one complete image without gaps between tiles
// During this process, it will also generate and return the Pixel map that corresponds to this image - tile borders removed too
// Printing to log may be silenced if only the image pixel map is desired
func PrintImage(image map[Pixel]Tile, verbose bool) map[Pixel]struct{} {
	var ipx map[Pixel]struct{} = make(map[Pixel]struct{})

	for ix := 0; ix < ImageDim; ix++ {
		for iy := 0; iy < ImageDim; iy++ {
			// ix and iy determine the tile being selected, but not the individual pixel on that tile inside
			// Additionally, the edges are trimmed from the tile
			// For the tiles on the left edge of the image, every pixel in the tile is shifted left only one, and for every other tile to the right, their pixels are shifted 2x their index + 1
			// For the tiles on the up edge of the image, every pixel in the tile is shifted up only one, and for every other tile below, their pixels are shifted 2x their index + 1
			// e.g. Tile (0, 0) Pixel (1, 1) would go to Image Pixel (ix*(TileDim-2)+(tpx.X-1), iy*(TileDim-2)+(tpx.Y-1)) --> (0*8+(1-1), 0*8+(1-1)) --> (0, 0)
			// e.g. Tile (0, 1) Pixel (3, 3) would go to Image Pixel (ix*(TileDim-2)+(tpx.X-1), iy*(TileDim-2)+(tpx.Y-1)) --> (0*8+(3-1), 1*8+(3-1)) --> (2, 10)
			// e.g. Tile (1, 2) Pixel (3, 4) would go to Image Pixel (ix*(TileDim-2)+(tpx.X-1), iy*(TileDim-2)+(tpx.Y-1)) --> (1*8+(3-1), 2*8+(4-1)) --> (10, 19)
			for tpx := range image[Pixel{ix, iy}].Pixels {
				if tpx.X > 0 && tpx.X < TileDim-1 && tpx.Y > 0 && tpx.Y < TileDim-1 {
					ipx[Pixel{ix*(TileDim-2) + (tpx.X - 1), iy*(TileDim-2) + (tpx.Y - 1)}] = struct{}{}
				}
			}
		}
	}

	if verbose {
		log.Println("Image Pixels:")
		PrintPixelMap(ipx, ImagePixelDim)
	}

	return ipx
}

// PrintPixelMap will nicely print a field of pixels as active or inactive dots
func PrintPixelMap(pixels map[Pixel]struct{}, dimSize int) {
	for x := 0; x < dimSize; x++ {
		var px []string
		for y := 0; y < dimSize; y++ {
			if _, a := pixels[Pixel{x, y}]; a {
				px = append(px, "#")
			} else {
				px = append(px, ".")
			}
		}
		log.Println(strings.Join(px, ""))
	}
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

	// Turn the input strings into tiles
	reID := regexp.MustCompile(`Tile (?P<ID>[0-9]+):`)
	tiles := make(map[string]Tile)
	currentTile := ""
	currentX := 0
	currentY := 0
	for _, line := range input {
		if reID.MatchString(line) {
			currentTile = reID.FindStringSubmatch(line)[1]
			currentX = 0
			currentY = 0
			tiles[currentTile] = Tile{currentTile, make(map[Pixel]struct{})}
		} else if line != "" {
			for _, char := range strings.Split(line, "") {
				if char == "#" {
					tiles[currentTile].Pixels[Pixel{currentX, currentY}] = struct{}{}
				}
				currentY++
			}
			currentX++
			currentY = 0
		}
	}

	var matches map[string]map[int]Tile = make(map[string]map[int]Tile)
	var matchTransforms map[string]map[int]map[string]string = make(map[string]map[int]map[string]string)
	for tile := range tiles {
		matches[tile], matchTransforms[tile] = FindCommonEdges(tiles, tile)
	}

	// hello darkness my old friend
	// To produce the final image, first start by finding the corner piece with matches on edges 1 and 2 - this goes into the image at Image Pixel (0,0)
	// Then consider the piece matching edge 1 (i.e. rightward): replace this tile in the tile set with its transformed version, perform another match search upon the tile set, then assign it to the image at (0,1)
	// This forms a repeatable chain until a tile is placed into (0, 11) - because the first row of the image is now full
	// For the edge of this end piece (which should be a corner), there won't be a match on edge 1, but there should be one for edge 2, so use that to drop down to the next row at (1, 11)
	// Then, continue as before, but instead of matching on edge 1, match on edge 3 and proceed leftward to (1,0)
	// Then drop down using edge 2 and repeat with this snaking process (hence, even row indexes move to the right, odd row indexes move to the left)
	// Repeat until the size of the image map is the same as the size of the tile set, as this means all tiles have been assigned a position in the image, arranged so that there is alignment
	var image map[Pixel]Tile = make(map[Pixel]Tile)
	var ix, iy int = 0, 0
	var nextTile Tile
	for tile := range matches {
		// Pull up the tile with matches on edge 1 and 2; this is the starting corner piece
		_, edge1 := matches[tile][1]
		_, edge2 := matches[tile][2]
		if edge1 && edge2 && len(matches[tile]) == 2 {
			nextTile = tiles[tile]
		}
	}
	for len(image) != len(tiles) {
		// Place tile i into image
		log.Println("Inserting tile", nextTile.ID, "into (", ix, ",", iy, ")")
		image[Pixel{ix, iy}] = nextTile
		log.Println("Image now contains", len(image), "of", len(tiles), "tiles")
		if len(image) != len(tiles) {
			// The image is not yet complete
			// Reference the next tile i+1 based on matches and image position, if there are any left to do (i.e. if len(image) != len(tiles))
			// The new next tile is obtained through looking for the edge needed in this iteration:
			//   Edge 2 if ix is an even number and iy is ImageDim-1 or if ix is an odd number and iy is 0 (i.e. the down edge)
			//   Edge 1 if ix is an even number and iy is not ImageDim-1 (i.e. the right edge)
			//   Edge 3 if ix is an odd number and iy is not 0 (i.e. the left edge)
			// The next tile is already provided in transformed form by FindCommonEdges as part of the matches variable return, so it should be assigned now to the tile set
			// Use this as an opportunity to determine the next ix and iy as well
			if ix%2 == 0 && iy == ImageDim-1 || ix%2 == 1 && iy == 0 {
				nextTile = matches[nextTile.ID][2]
				ix++
			} else if ix%2 == 0 {
				nextTile = matches[nextTile.ID][1]
				iy++
			} else if ix%2 == 1 {
				nextTile = matches[nextTile.ID][3]
				iy--
			}
			tiles[nextTile.ID] = nextTile
			// Then, redo matches for this next tile - this will refresh the matches so that the above process can be easily done
			matches[nextTile.ID], matchTransforms[nextTile.ID] = FindCommonEdges(tiles, nextTile.ID)
		}
	}
	PrintImageTileIDs(image)
	log.Println("P1: Corner ID Product:", CornerIDProduct(image))

	imagePixels := PrintImage(image, true)

	// The monster pattern is 20 long and 3 high, and contains 15 dots:
	//            1111111111
	//  01234567890123456789
	// 0                  #
	// 1#    ##    ##    ###
	// 2 #  #  #  #  #  #
	mons := 0
	monDef := []Pixel{
		Pixel{0, 18},
		Pixel{1, 0},
		Pixel{1, 5},
		Pixel{1, 6},
		Pixel{1, 11},
		Pixel{1, 12},
		Pixel{1, 17},
		Pixel{1, 18},
		Pixel{1, 19},
		Pixel{2, 1},
		Pixel{2, 4},
		Pixel{2, 7},
		Pixel{2, 10},
		Pixel{2, 13},
		Pixel{2, 16},
	}
	axes := []string{"", "x", "y", "xy"}
	for _, axis := range axes {
		for rotation := 0; rotation < 4; rotation++ {
			imagePixels := RotatePixels(ReflectPixels(imagePixels, axis, ImagePixelDim), rotation, ImagePixelDim)
			// Instead of iterating over the image pixels, iterate simply over indexes 0 to ImagePixelDim in each direction
			// This is because imagePixels isn't a complete array, just a sparse map
			for x := 0; x < ImagePixelDim; x++ {
				for y := 0; y < ImagePixelDim; y++ {
					// To confirm that a monster is not present, check that, offset from the current X and Y, all of the pixels defined by monDef are in the image pixel set
					// If any one of them is not present, then it's not a match
					var match bool = true
					for _, mdpx := range monDef {
						if _, found := imagePixels[Pixel{x + mdpx.X, y + mdpx.Y}]; !found {
							match = false
							break
						}
					}
					if match {
						mons++
					}
				}
			}
		}
		// There's no need to progress any further if mons > 0 at this point
		if mons > 0 {
			break
		}
	}
	log.Println("P2 | Sea monsters:", mons, "| Water roughness:", len(imagePixels)-mons*len(monDef))
}
