package main

import (
	"fmt"
	"os"
	"strings"
)

type AntennaLocation struct {
	x, y int
}

type AntennaType = rune

type Antenna struct {
	antennaType AntennaType
	location    AntennaLocation
}

type AntennaArea struct {
	height           int
	width            int
	inputLines       []string
	antennas         []Antenna
	antennaLocations map[AntennaType][]Antenna
}

//	012345678
//
// 0 ..........
// 1 ..........
// 2 ...l......
// 3 ..........
// 4 .....o....
// 5 ..........
// 6 .......a..
// 7 ..........
// 8 .........a
// 9 ..........

// l = 2,3
// o = 4,5
// a = o+(o-l)
// a2 = o+2*(o-l)
func (o AntennaLocation) getAntinodeLocation(l AntennaLocation, area AntennaArea) *AntennaLocation {
	loc := AntennaLocation{
		x: 2*o.x - l.x,
		y: 2*o.y - l.y,
	}
	if area.isLocationInArea(loc) {
		return &loc
	} else {
		return nil
	}
}

func (o AntennaLocation) getMultiAntinodeLocation(l AntennaLocation, area AntennaArea) []AntennaLocation {
	locations := []AntennaLocation{
		l,
	}
	for i := 1; ; i++ {
		l := AntennaLocation{
			x: o.x + i*(o.x-l.x),
			y: o.y + i*(o.y-l.y),
		}
		if area.isLocationInArea(l) {
			locations = append(locations, l)
		} else {
			break
		}
	}
	return locations
}

func (area AntennaArea) isLocationInArea(l AntennaLocation) bool {
	return l.x >= 0 && l.y >= 0 && l.x < area.width && l.y < area.height
}

func newAntennaArea(inputLines []string) AntennaArea {
	area := AntennaArea{}
	area.height = len(inputLines)
	area.width = len(inputLines[0])
	area.inputLines = inputLines
	area.antennaLocations = map[AntennaType][]Antenna{}
	area.antennas = []Antenna{}
	for y, l := range inputLines {
		for x, p := range l {
			if p == '.' {
				continue
			}
			location := AntennaLocation{
				x: x,
				y: y,
			}
			antenna := Antenna{
				antennaType: p,
				location:    location,
			}
			area.antennas = append(area.antennas, antenna)
			antennaList, ok := area.antennaLocations[p]
			if !ok {
				antennaList = []Antenna{}
			}
			area.antennaLocations[p] = append(antennaList, antenna)
		}
	}
	return area
}

func main() {
	input, _ := os.ReadFile("input.txt")
	inputLines := strings.Split(string(input), "\n")
	area := newAntennaArea(inputLines)
	antinodes := map[AntennaLocation]bool{}

	for _, antenna := range area.antennas {
		for _, o := range area.antennaLocations[antenna.antennaType] {
			if antenna != o {
				antinodeLocation := antenna.location.getAntinodeLocation(o.location, area)
				if antinodeLocation != nil {
					antinodes[*antinodeLocation] = true
				}
			}
		}
	}
	fmt.Println("Part 1:", len(antinodes))

	for _, antenna := range area.antennas {
		for _, o := range area.antennaLocations[antenna.antennaType] {
			if antenna != o {
				antinodeLocations := antenna.location.getMultiAntinodeLocation(o.location, area)
				for _, l := range antinodeLocations {
					antinodes[l] = true
				}
			}
		}
	}
	fmt.Println("Part 2:", len(antinodes))
}
