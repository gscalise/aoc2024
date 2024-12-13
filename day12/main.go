package main

import (
	"container/list"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
)

var directions = []Direction{
	{-1, 0}, // UP
	{0, 1},  // RIGHT
	{1, 0},  // DOWN
	{0, -1}, // LEFT
}

var directionPairs = [][2]Direction{
	{directions[0], directions[1]},
	{directions[1], directions[2]},
	{directions[2], directions[3]},
	{directions[3], directions[0]},
}

type Direction struct {
	dr, dc int
}

type Location struct {
	r, c int
}

func (d Direction) add(o Direction) Direction {
	return Direction{
		dc: d.dc + o.dc,
		dr: d.dr + o.dr,
	}
}

func (l Location) add(d Direction) Location {
	return Location{
		r: l.r + d.dr,
		c: l.c + d.dc,
	}
}

type FarmArea struct {
	rows       int
	columns    int
	inputLines []string
}

type FarmAreaType string

func newLocation(r, c int) Location {
	return Location{
		r,
		c,
	}
}

func getFarmArea(lines []string) FarmArea {
	return FarmArea{
		rows:       len(lines),
		columns:    len(lines[0]),
		inputLines: lines,
	}
}

func (t FarmArea) isLocationInArea(l Location) bool {
	return l.r >= 0 && l.c >= 0 && l.r < t.rows && l.c < t.columns
}

func (l Location) isInArea(t FarmArea) bool {
	return t.isLocationInArea(l)
}

func (l Location) canBeAreaCorner(t FarmArea) bool {
	return l.r >= 0 && l.c >= 0 && l.r <= t.rows && l.c <= t.columns
}

func (t FarmArea) valueAt(l Location) FarmAreaType {
	if t.isLocationInArea(l) {
		return FarmAreaType(t.inputLines[l.r][l.c])
	} else {
		return FarmAreaType("")
	}
}

func buggyMainPart1(f string) {
	startTime := time.Now()
	input, _ := os.ReadFile(f)
	farmArea := getFarmArea(strings.Split(string(input), "\n"))
	typeArea := map[FarmAreaType]int{}
	typePerimeter := map[FarmAreaType]int{}
	fmt.Println("Value at 0,0", farmArea.valueAt(newLocation(0, 0)))
	ts := time.Since(startTime)
	for r := range farmArea.rows {
		for c := range farmArea.columns {
			location := newLocation(r, c)
			r := farmArea.valueAt(location)
			typeArea[r] += 1
			for _, d := range directions {
				dLocation := location.add(d)
				if !farmArea.isLocationInArea(dLocation) || farmArea.valueAt(dLocation) != r {
					typePerimeter[r] += 1
				}
			}
		}
	}
	totalSum := 0
	for r := range typeArea {
		totalSum += typeArea[r] * typePerimeter[r]
	}
	fmt.Println(totalSum)
	fmt.Println("Took", ts.Microseconds(), "μs")
}

type FarmAreaRegion struct {
	regionId      int
	areaType      FarmAreaType
	areaSurface   int
	areaPerimeter int
	areaSides     int
	locations     []Location
}

func (l Location) isInRegion(r FarmAreaRegion) bool {
	return slices.Contains(r.locations, l)
}

func (l Location) isInRegionBorder(r FarmAreaRegion) bool {
	if !l.isInRegion(r) {
		return false
	} else {
		for _, d := range directions {
			ld := l.add(d)
			if !ld.isInRegion(r) {
				return true
			}
		}
	}
	return false
}

func (l Location) isInTopRegionBorder(r FarmAreaRegion) bool {
	if !l.isInRegion(r) {
		return false
	} else {
		return !l.add(directions[0]).isInRegion(r)
	}
}

func mainPart1(f string) {
	startTime := time.Now()
	input, _ := os.ReadFile(f)
	farmArea := getFarmArea(strings.Split(string(input), "\n"))
	visitedLocations := map[Location]bool{}
	regions := map[int]*FarmAreaRegion{}
	regionId := 0
	ts := time.Since(startTime)
	deque := list.New()
	for r := range farmArea.rows {
		for c := range farmArea.columns {
			location := newLocation(r, c)
			if !visitedLocations[location] {
				r := farmArea.valueAt(location)
				regionId += 1
				region := &FarmAreaRegion{
					regionId:      regionId,
					areaType:      r,
					areaSurface:   0,
					areaPerimeter: 0,
				}
				regions[regionId] = region
				deque.PushBack(location)
				for deque.Len() > 0 {
					qLocation := deque.Front().Value.(Location)
					deque.Remove(deque.Front())
					if !visitedLocations[qLocation] {
						visitedLocations[qLocation] = true
						region.areaSurface += 1
						for _, d := range directions {
							dLocation := qLocation.add(d)
							if !farmArea.isLocationInArea(dLocation) {
								region.areaPerimeter += 1
							} else if farmArea.valueAt(dLocation) == r {
								deque.PushBack(dLocation)
							} else {
								region.areaPerimeter += 1
							}
						}
					}
				}
			}
		}
	}
	totalSum := 0
	for _, r := range regions {
		totalSum += r.areaPerimeter * r.areaSurface
	}
	fmt.Println(totalSum)
	fmt.Println("Took", ts.Microseconds(), "μs")
}

type LocationDirection struct {
	location        Location
	movingDirection Direction
}

func (d Direction) TurnRight() Direction {
	return Direction{
		dc: -d.dr,
		dr: d.dc,
	}
}

func (d Direction) TurnLeft() Direction {
	return Direction{
		dc: d.dr,
		dr: -d.dc,
	}
}

type RegionBorderVisitor struct {
	region           *FarmAreaRegion
	location         Location
	movingDirection  Direction
	turnCount        int
	visitedLocations *(map[LocationDirection]bool)
}

func (d Direction) String() string {
	if d.dc == 1 {
		return "RIGHT"
	} else if d.dc == -1 {
		return "LEFT"
	} else if d.dr == -1 {
		return "UP"
	} else {
		return "DOWN"
	}
}

func (l Location) String() string {
	return fmt.Sprint(l.c+1, ",", l.r+1)
}

func (r *RegionBorderVisitor) Advance() bool {
	candidateLocation := r.location.add(r.movingDirection)
	ldKey := LocationDirection{
		location:        candidateLocation,
		movingDirection: r.movingDirection,
	}
	fmt.Println("Trying to advance", r.movingDirection, "from", r.location, "to", candidateLocation)
	if (*(r.visitedLocations))[ldKey] {
		fmt.Println("Been here. Exiting after", r.turnCount, "turns")
		return false
	}
	(*(r.visitedLocations))[ldKey] = true
	if candidateLocation.isInRegion(*r.region) {
		if candidateLocation.add(r.movingDirection.TurnLeft()).isInRegion(*r.region) {
			r.movingDirection = r.movingDirection.TurnLeft()
			fmt.Println("Turn left at", r.location, "new direction is", r.movingDirection)
			r.turnCount += 1
		}
		r.location = candidateLocation
		return true
	} else {
		r.movingDirection = r.movingDirection.TurnRight()
		fmt.Println("Turn right at", r.location, "new direction is", r.movingDirection)
		r.turnCount += 1
		return true
	}
}

func newRegionBorderVisitor(region *FarmAreaRegion, location Location, visitedLocations *map[LocationDirection]bool) *RegionBorderVisitor {
	fmt.Println("\n\n\nStart in region", region.areaType, "location", location)
	(*visitedLocations)[LocationDirection{location: location, movingDirection: directions[1]}] = true
	return &RegionBorderVisitor{
		region:           region,
		location:         location,
		movingDirection:  directions[1],
		turnCount:        0,
		visitedLocations: visitedLocations,
	}
}

func mainPart2(f string) {
	startTime := time.Now()
	input, _ := os.ReadFile(f)
	farmArea := getFarmArea(strings.Split(string(input), "\n"))
	visitedLocations := map[Location]bool{}
	regions := map[int]*FarmAreaRegion{}
	locationRegionMap := map[Location]*FarmAreaRegion{}
	regionId := 0
	ts := time.Since(startTime)
	deque := list.New()
	for r := range farmArea.rows {
		for c := range farmArea.columns {
			location := newLocation(r, c)
			if !visitedLocations[location] {
				r := farmArea.valueAt(location)
				regionId += 1
				region := &FarmAreaRegion{
					regionId: regionId,
					areaType: r,
				}
				regions[regionId] = region
				deque.PushBack(location)
				for deque.Len() > 0 {
					qLocation := deque.Front().Value.(Location)
					deque.Remove(deque.Front())
					if !visitedLocations[qLocation] {
						visitedLocations[qLocation] = true
						region.locations = append(region.locations, qLocation)
						locationRegionMap[qLocation] = region
						for _, d := range directions {
							dLocation := qLocation.add(d)
							if farmArea.isLocationInArea(dLocation) && farmArea.valueAt(dLocation) == r {
								deque.PushBack(dLocation)
							}
						}
					}
				}
				// fmt.Printf("%+v\n", *regions[regionId])
			}
		}
	}

	totalSum := 0
	for _, r := range regions {
		turnCount := 0
		visitedLocationDirections := map[LocationDirection]bool{}
		for _, l := range r.locations {
			if l.isInTopRegionBorder(*r) {
				visitor := newRegionBorderVisitor(r, l, &visitedLocationDirections)
				for visitor.Advance() {
					// Advance
				}
				fmt.Println("Finished visiting with", visitor.turnCount, "turns")
				turnCount += visitor.turnCount
			}
		}
		r.areaSides = turnCount
		fmt.Println("Found", r.areaSides, "corners for region", r.regionId, "of type", r.areaType)
		totalSum += r.areaSides * len(r.locations)
	}
	fmt.Println(totalSum)
	fmt.Println("Took", ts.Microseconds(), "μs")
}

func main() {
	mainPart2("input.txt")
}
