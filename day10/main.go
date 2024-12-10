package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

var directions = []Direction{
	Direction{-1, 0}, // UP
	Direction{0, 1},  // RIGHT
	Direction{1, 0},  // DOWN
	Direction{0, -1}, // LEFT
}

type Direction struct {
	dr, dc int
}

type Location struct {
	r, c int
}

func (l Location) add(d Direction) Location {
	return Location{
		r: l.r + d.dr,
		c: l.c + d.dc,
	}
}

type TrailArea struct {
	rows       int
	columns    int
	inputLines []string
}

func getTrailArea(lines []string) TrailArea {
	return TrailArea{
		rows:       len(lines),
		columns:    len(lines[0]),
		inputLines: lines,
	}
}

func (t TrailArea) isLocationInArea(l Location) bool {
	return l.r >= 0 && l.c >= 0 && l.r < t.rows && l.c < t.columns
}

func (t TrailArea) valueAt(l Location) int {
	if t.isLocationInArea(l) {
		value, _ := strconv.Atoi(string(t.inputLines[l.r][l.c]))
		return value
	} else {
		return -1
	}
}

func (t TrailArea) getTrails(currentLocation Location) []Location {
	curVal := t.valueAt(currentLocation)
	if curVal == 9 {
		return []Location{currentLocation}
	} else {
		foundTrailEnds := []Location{}
		for _, d := range directions {
			newLoc := currentLocation.add(d)
			if t.isLocationInArea(newLoc) {
				newLocVal := t.valueAt(newLoc)
				if newLocVal == curVal+1 {
					trails := t.getTrails(newLoc)
					for _, t := range trails {
						if !slices.Contains(foundTrailEnds, t) {
							foundTrailEnds = append(foundTrailEnds, t)
						}
					}
				}
			}
		}
		return foundTrailEnds
	}
}

func (t TrailArea) getNumTrails(currentLocation Location) int {
	curVal := t.valueAt(currentLocation)
	if curVal == 9 {
		return 1
	} else {
		sum := 0
		for _, d := range directions {
			newLoc := currentLocation.add(d)
			if t.isLocationInArea(newLoc) {
				newLocVal := t.valueAt(newLoc)
				if newLocVal == curVal+1 {
					sum += t.getNumTrails(newLoc)
				}
			}
		}
		return sum
	}
}

func main() {
	input, _ := os.ReadFile("demo-input.txt")
	inputLines := strings.Split(string(input), "\n")
	part1Sum := 0
	part2Sum := 0

	trailArea := getTrailArea(inputLines)

	for r, row := range trailArea.inputLines {
		for c, col := range row {
			if col == '0' {
				part1Sum += len(trailArea.getTrails(Location{r, c}))
				part2Sum += trailArea.getNumTrails(Location{r, c})
			}
		}
	}
	fmt.Println("Part 1:", part1Sum)
	fmt.Println("Part 2:", part2Sum)
}
