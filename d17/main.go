package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// Point3D represents a 3D point in space, integers only.
type Point3D struct {
	X int
	Y int
	Z int
}

// Point4D represents a 4D point in space, integers only.
type Point4D struct {
	W int
	X int
	Y int
	Z int
}

// Bounds3D finds the minimum and maximum X, Y, and Z values for a given set of points, returning [2]int for X, Y, and Z lower and upper bounds
// If the point map is empty, the bounds are all 0
func Bounds3D(points map[Point3D]struct{}) ([2]int, [2]int, [2]int) {
	var boundX, boundY, boundZ [2]int = [2]int{0, 0}, [2]int{0, 0}, [2]int{0, 0} // This has a gotcha, when any of X, Y, or Z has a minimum value above 0
	// To avoid this gotcha, just set min and max for all of them to be the first real point encountered
	realInit := false
	if len(points) > 0 {
		for point := range points {
			if !realInit {
				boundX, boundY, boundZ = [2]int{point.X, point.X}, [2]int{point.Y, point.Y}, [2]int{point.Z, point.Z}
				realInit = true
				continue
			}
			if point.X < boundX[0] {
				boundX[0] = point.X
			}
			if point.X > boundX[1] {
				boundX[1] = point.X
			}
			if point.Y < boundY[0] {
				boundY[0] = point.Y
			}
			if point.Y > boundY[1] {
				boundY[1] = point.Y
			}
			if point.Z < boundZ[0] {
				boundZ[0] = point.Z
			}
			if point.Z > boundZ[1] {
				boundZ[1] = point.Z
			}
		}
	}
	return boundX, boundY, boundZ
}

// Bounds4D finds the minimum and maximum W, X, Y, and Z values for a given set of points, returning [2]int for W, X, Y, and Z lower and upper bounds
// If the point map is empty, the bounds are all 0
func Bounds4D(points map[Point4D]struct{}) ([2]int, [2]int, [2]int, [2]int) {
	var boundW, boundX, boundY, boundZ [2]int = [2]int{0, 0}, [2]int{0, 0}, [2]int{0, 0}, [2]int{0, 0} // This has a gotcha, when any of W, X, Y, or Z has a minimum value above 0
	// To avoid this gotcha, just set min and max for all of them to be the first real point encountered
	realInit := false
	if len(points) > 0 {
		for point := range points {
			if !realInit {
				boundW, boundX, boundY, boundZ = [2]int{point.W, point.W}, [2]int{point.X, point.X}, [2]int{point.Y, point.Y}, [2]int{point.Z, point.Z}
				realInit = true
				continue
			}
			if point.W < boundW[0] {
				boundW[0] = point.W
			}
			if point.W > boundW[1] {
				boundW[1] = point.W
			}
			if point.X < boundX[0] {
				boundX[0] = point.X
			}
			if point.X > boundX[1] {
				boundX[1] = point.X
			}
			if point.Y < boundY[0] {
				boundY[0] = point.Y
			}
			if point.Y > boundY[1] {
				boundY[1] = point.Y
			}
			if point.Z < boundZ[0] {
				boundZ[0] = point.Z
			}
			if point.Z > boundZ[1] {
				boundZ[1] = point.Z
			}
		}
	}
	return boundW, boundX, boundY, boundZ
}

// Print3DSpace ingests a map of Point3D - for every point indicated it will return #
// Any points not indicated in the map but within the min/max X, Y, or Z will be indicated with .
// It returns X and Y grids, with slices of Z, as well as the 3D bounds
func Print3DSpace(points map[Point3D]struct{}) {
	boundX, boundY, boundZ := Bounds3D(points)

	log.Println("Bounds: X", boundX, "Y", boundY, "Z", boundZ)

	// Print one Z-plane at a time
	for z := boundZ[0]; z <= boundZ[1]; z++ {
		log.Println("Z =", z)
		for x := boundX[0]; x <= boundX[1]; x++ {
			var chars []string
			for y := boundY[0]; y <= boundY[1]; y++ {
				if _, def := points[Point3D{x, y, z}]; def {
					chars = append(chars, "#")
				} else {
					chars = append(chars, ".")
				}
			}
			log.Println(chars)
		}
	}
}

// Print4DSpace ingests a map of Point4D - for every point indicated it will return #
// Any points not indicated in the map but within the min/max W, X, Y, or Z will be indicated with .
// It returns X and Y grids, with slices of W & Z, as well as the 4D bounds
func Print4DSpace(points map[Point4D]struct{}) {
	boundW, boundX, boundY, boundZ := Bounds4D(points)

	log.Println("Bounds:", "W", boundW, "X", boundX, "Y", boundY, "Z", boundZ)

	// Print one W & Z pixel at a time
	for w := boundW[0]; w <= boundW[1]; w++ {
		for z := boundZ[0]; z <= boundZ[1]; z++ {
			log.Println("W =", w, "Z =", z)
			for x := boundX[0]; x <= boundX[1]; x++ {
				var chars []string
				for y := boundY[0]; y <= boundY[1]; y++ {
					if _, def := points[Point4D{w, x, y, z}]; def {
						chars = append(chars, "#")
					} else {
						chars = append(chars, ".")
					}
				}
				log.Println(chars)
			}
		}
	}
}

// AtoP3D ingests a map of points and a target point and returns true if the point should be added based on the rule:
// If a point is absent but exactly 3 of its neighbors are present, add the point.
// If the target submitted is already present, this function returns true to avoid accidental state changes
func AtoP3D(points map[Point3D]struct{}, target Point3D) bool {
	if _, present := points[target]; present {
		return true
	}
	neighbours := 0
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			for dz := -1; dz <= 1; dz++ {
				if dx == 0 && dy == 0 && dz == 0 {
					continue
				}
				if _, present := points[Point3D{target.X + dx, target.Y + dy, target.Z + dz}]; present {
					neighbours++
				}
			}
		}
	}
	// log.Println("DEBUG | Absent Point", target, "| Neighbours:", neighbours)
	if neighbours != 3 {
		return false
	}
	return true
}

// AtoP4D ingests a map of points and a target point and returns true if the point should be added based on the rule:
// If a point is absent but exactly 3 of its neighbors are present, add the point.
// If the target submitted is already present, this function returns true to avoid accidental state changes
func AtoP4D(points map[Point4D]struct{}, target Point4D) bool {
	if _, present := points[target]; present {
		return true
	}
	neighbours := 0
	for dw := -1; dw <= 1; dw++ {
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				for dz := -1; dz <= 1; dz++ {
					if dw == 0 && dx == 0 && dy == 0 && dz == 0 {
						continue
					}
					if _, present := points[Point4D{target.W + dw, target.X + dx, target.Y + dy, target.Z + dz}]; present {
						neighbours++
					}
				}
			}
		}
	}
	// log.Println("DEBUG | Absent Point", target, "| Neighbours:", neighbours)
	if neighbours != 3 {
		return false
	}
	return true
}

// PtoA3D ingests a map of points and a target point and returns true if the point should be absent based on the rule:
// If a point is present and exactly 2 or 3 of its neighbors are also present, the point stays. Otherwise, delete the point.
// If the target submitted is already absent, this function returns true to avoid accidental state changes
func PtoA3D(points map[Point3D]struct{}, target Point3D) bool {
	if _, present := points[target]; !present {
		return true
	}
	neighbours := 0
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			for dz := -1; dz <= 1; dz++ {
				if dx == 0 && dy == 0 && dz == 0 {
					continue
				}
				if _, present := points[Point3D{target.X + dx, target.Y + dy, target.Z + dz}]; present {
					neighbours++
				}
			}
		}
	}
	// log.Println("DEBUG | Present Point", target, "| Neighbours:", neighbours)
	if neighbours == 2 || neighbours == 3 {
		return false
	}
	return true
}

// PtoA4D ingests a map of points and a target point and returns true if the point should be absent based on the rule:
// If a point is present and exactly 2 or 3 of its neighbors are also present, the point stays. Otherwise, delete the point.
// If the target submitted is already absent, this function returns true to avoid accidental state changes
func PtoA4D(points map[Point4D]struct{}, target Point4D) bool {
	if _, present := points[target]; !present {
		return true
	}
	neighbours := 0
	for dw := -1; dw <= 1; dw++ {
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				for dz := -1; dz <= 1; dz++ {
					if dw == 0 && dx == 0 && dy == 0 && dz == 0 {
						continue
					}
					if _, present := points[Point4D{target.W + dw, target.X + dx, target.Y + dy, target.Z + dz}]; present {
						neighbours++
					}
				}
			}
		}
	}
	// log.Println("DEBUG | Present Point", target, "| Neighbours:", neighbours)
	if neighbours == 2 || neighbours == 3 {
		return false
	}
	return true
}

// Evolve3DSpace ingests a set of present points and returns the state of the space after one iteration of the evolution rules
func Evolve3DSpace(points map[Point3D]struct{}) map[Point3D]struct{} {
	evolvedPoints := make(map[Point3D]struct{})
	// The evolution rules:
	//   - If a point is present and exactly 2 or 3 of its neighbors are also present, the point stays. Otherwise, delete the point.
	//   - If a point is absent but exactly 3 of its neighbors are present, add the point.
	// Note now that this requires checking beyond the bounds of the current space by one in each direction, plus checking point values that don't exist in the active space.
	// This is most simply performed with a raster scan of a hypothetical space one beyond the min and max bounds of the current space.
	boundX, boundY, boundZ := Bounds3D(points)

	// For every absent point, check if it should now be present, and add it if so
	// For every present point, check if it should now be absent, and add it if not
	for x := boundX[0] - 1; x <= boundX[1]+1; x++ {
		for y := boundY[0] - 1; y <= boundY[1]+1; y++ {
			for z := boundZ[0] - 1; z <= boundZ[1]+1; z++ {
				if _, def := points[Point3D{x, y, z}]; def && !PtoA3D(points, Point3D{x, y, z}) {
					// log.Println("DEBUG | Point", Point3D{x, y, z}, "will stay in")
					evolvedPoints[Point3D{x, y, z}] = struct{}{}
					continue
				}
				if _, def := points[Point3D{x, y, z}]; !def && AtoP3D(points, Point3D{x, y, z}) {
					// log.Println("DEBUG | Point", Point3D{x, y, z}, "will be added")
					evolvedPoints[Point3D{x, y, z}] = struct{}{}
					continue
				}
			}
		}
	}

	return evolvedPoints
}

// Evolve4DSpace ingests a set of present points and returns the state of the space after one iteration of the evolution rules
func Evolve4DSpace(points map[Point4D]struct{}) map[Point4D]struct{} {
	evolvedPoints := make(map[Point4D]struct{})
	// The evolution rules:
	//   - If a point is present and exactly 2 or 3 of its neighbors are also present, the point stays. Otherwise, delete the point.
	//   - If a point is absent but exactly 3 of its neighbors are present, add the point.
	// Note now that this requires checking beyond the bounds of the current space by one in each direction, plus checking point values that don't exist in the active space.
	// This is most simply performed with a raster scan of a hypothetical space one beyond the min and max bounds of the current space.
	boundW, boundX, boundY, boundZ := Bounds4D(points)

	// For every absent point, check if it should now be present, and add it if so
	// For every present point, check if it should now be absent, and add it if not
	for w := boundW[0] - 1; w <= boundW[1]+1; w++ {
		for x := boundX[0] - 1; x <= boundX[1]+1; x++ {
			for y := boundY[0] - 1; y <= boundY[1]+1; y++ {
				for z := boundZ[0] - 1; z <= boundZ[1]+1; z++ {
					if _, def := points[Point4D{w, x, y, z}]; def && !PtoA4D(points, Point4D{w, x, y, z}) {
						// log.Println("DEBUG | Point", Point4D{w, x, y, z}, "will stay in")
						evolvedPoints[Point4D{w, x, y, z}] = struct{}{}
						continue
					}
					if _, def := points[Point4D{w, x, y, z}]; !def && AtoP4D(points, Point4D{w, x, y, z}) {
						// log.Println("DEBUG | Point", Point4D{w, x, y, z}, "will be added")
						evolvedPoints[Point4D{w, x, y, z}] = struct{}{}
						continue
					}
				}
			}
		}
	}

	return evolvedPoints
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

	// P1: Iterate 6 times; the number of active points is just the length of the active space
	// Construct an initial state as a map of points to empty structs
	// Presence in the map indicates activation, and removal indicates inactivation
	// The field of play may expand (the initial state does not provide a boundary on physical space)
	space3 := make(map[int]map[Point3D]struct{})
	space3[0] = make(map[Point3D]struct{})
	// For every # in input, submit it to active state with row 0 being X = 0, column 0 being Y = 0, and the whole plane being Z = 0
	for x := range input {
		points := strings.Split(input[x], "")
		for y := range points {
			if points[y] == "#" {
				space3[0][Point3D{x, y, 0}] = struct{}{}
			}
		}
	}
	for iter := 0; iter <= 6; iter++ {
		if iter != 0 {
			space3[iter] = Evolve3DSpace(space3[iter-1])
		}
		log.Println("======== ITERATION", iter)
		Print3DSpace(space3[iter])
		log.Println("P1 | Iteration", iter, "| Active:", len(space3[iter]))
	}

	// P2: Curse your sudden but inevitable fourth dimension
	space4 := make(map[int]map[Point4D]struct{})
	space4[0] = make(map[Point4D]struct{})
	// For every # in input, submit it to active state with row 0 being X = 0, column 0 being Y = 0, and W, Z = 0
	for x := range input {
		points := strings.Split(input[x], "")
		for y := range points {
			if points[y] == "#" {
				space4[0][Point4D{0, x, y, 0}] = struct{}{}
			}
		}
	}
	for iter := 0; iter <= 6; iter++ {
		if iter != 0 {
			space4[iter] = Evolve4DSpace(space4[iter-1])
		}
		log.Println("======== ITERATION", iter)
		Print4DSpace(space4[iter])
		log.Println("P2 | Iteration", iter, "| Active:", len(space4[iter]))
	}
}
