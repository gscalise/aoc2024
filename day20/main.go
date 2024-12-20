package main

import (
	"fmt"
	"os"
	"strings"

	. "bautik.net/advent2024/helpers"
)

type RaceTrack struct {
	width, height   int
	wallMap         map[Position]bool
	stepsInPosition map[Position]int
	startTile       Position
	endTile         Position
}

func (m *RaceTrack) addWall(p Position) {
	m.wallMap[p] = true
}

func (m *RaceTrack) inBounds(p Position) bool {
	return p.X >= 0 && p.X < m.width && p.Y >= 0 && p.Y < m.height
}

func (m *RaceTrack) String() string {
	var sb strings.Builder

	for r := range m.height {
		for c := range m.width {
			position := Position{X: c, Y: r}
			if m.wallMap[position] {
				sb.WriteRune('#')
			} else if position == m.startTile {
				sb.WriteRune('S')
			} else if position == m.endTile {
				sb.WriteRune('E')
			} else {
				sb.WriteRune(' ')
			}
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

func newRaceTrack(input string) *RaceTrack {
	lines := strings.Split(input, "\n")
	height := len(lines)
	width := len(lines[0])

	raceTrack := RaceTrack{
		width:           width,
		height:          height,
		wallMap:         map[Position]bool{},
		stepsInPosition: map[Position]int{},
	}

	for r, l := range lines {
		for c, o := range l {
			if o == '#' {
				raceTrack.addWall(Position{X: c, Y: r})
			}
			if o == 'S' {
				raceTrack.startTile = Position{X: c, Y: r}
			}
			if o == 'E' {
				raceTrack.endTile = Position{X: c, Y: r}
			}
		}
	}

	return &raceTrack
}

func (raceTrack *RaceTrack) isValidPosition(p Position) bool {
	wall := raceTrack.wallMap[p]
	return raceTrack.inBounds(p) && !wall
}

func (raceTrack *RaceTrack) isInBounds(p Position) bool {
	return p.X >= 0 && p.Y >= 0 && p.X <= raceTrack.width-1 && p.Y <= raceTrack.height-1
}

type TripNode struct {
	position Position
	path     []Position
}

func (raceTrack *RaceTrack) getNeighbors(p Position) (neighbors []Position) {
	neighbors = []Position{}
	for _, d := range Directions {
		pd := p.Displace(*d)
		if raceTrack.isValidPosition(pd) {
			neighbors = append(neighbors, pd)
		}
	}
	return neighbors
}

func (raceTrack *RaceTrack) solve() (pathLen int, pathFound bool) {
	endPosition := raceTrack.endTile
	visited := map[Position]bool{}
	pathLen = 0
	queue := []Position{
		raceTrack.startTile,
	}
	for len(queue) > 0 {
		levelSize := len(queue)

		for range levelSize {
			current := queue[0]
			queue = queue[1:]
			raceTrack.stepsInPosition[current] = pathLen
			if current == endPosition {
				return pathLen, true
			}
			for _, neighbour := range raceTrack.getNeighbors(current) {
				if visited[neighbour] {
					continue
				}
				visited[neighbour] = true
				queue = append(queue, neighbour)
			}
		}
		pathLen++
	}
	return -1, false
}

func (raceTrack *RaceTrack) GetCheatPositions(p Position, limit int) []Position {
	positions := []Position{}
	for c := -limit; c <= limit; c++ {
		for r := -limit; r <= limit; r++ {
			if Abs(c)+Abs(r) > limit {
				continue
			}
			candPos := Position{X: p.X + c, Y: p.Y + r}
			if candPos != p && raceTrack.isValidPosition(candPos) {
				positions = append(positions, candPos)
			}
		}
	}
	return positions
}

func getStepsToPos(o Position, t Position) int {
	return Abs(o.X-t.X) + Abs(o.Y-t.Y)
}

func (raceTrack *RaceTrack) countCheats(limit int) (cheatSaveMap map[Position]int, cheatSaveCounts map[int]int) {
	cheatSaveMap = map[Position]int{}
	cheatSaveCounts = map[int]int{}
	for p := range raceTrack.stepsInPosition {
		maxSavedSteps := 0
		for _, cp := range raceTrack.GetCheatPositions(p, limit) {
			if raceTrack.stepsInPosition[cp]+getStepsToPos(p, cp) > raceTrack.stepsInPosition[p] {
				savedSteps := raceTrack.stepsInPosition[cp] - raceTrack.stepsInPosition[p] - getStepsToPos(p, cp)
				cheatSaveCounts[savedSteps]++
				if savedSteps > maxSavedSteps {
					maxSavedSteps = savedSteps
				}
			}
		}
		cheatSaveMap[p] = maxSavedSteps
	}
	return cheatSaveMap, cheatSaveCounts
}

func loadRaceTrack(file string) *RaceTrack {
	fileBytes, _ := os.ReadFile(file)
	track := newRaceTrack(string(fileBytes))
	track.solve()
	return track
}

func (raceTrack *RaceTrack) PrintWithSteps() string {
	var sb strings.Builder

	line := strings.Repeat("+-----", raceTrack.width)

	for r := range raceTrack.height {
		sb.WriteString(line)
		sb.WriteString("+\n")
		for c := range raceTrack.width {
			p := Position{X: c, Y: r}
			if raceTrack.wallMap[p] {
				sb.WriteString("|#####")
			} else if p == raceTrack.startTile {
				sb.WriteString("| SSS ")
			} else if p == raceTrack.endTile {
				sb.WriteString(fmt.Sprintf("|E:%3d", raceTrack.stepsInPosition[p]))
			} else {
				sb.WriteString(fmt.Sprintf("|%4d ", raceTrack.stepsInPosition[p]))
			}
		}
		sb.WriteString("|\n")
	}
	sb.WriteString(line)
	sb.WriteString("+")

	return sb.String()
}
func (raceTrack *RaceTrack) PrintWitCheatSteps(limit int) string {

	cheatMap, _ := raceTrack.countCheats(limit)
	var sb strings.Builder

	line := strings.Repeat("+-----", raceTrack.width)

	for r := range raceTrack.height {
		sb.WriteString(line)
		sb.WriteString("+\n")
		for c := range raceTrack.width {
			p := Position{X: c, Y: r}
			if raceTrack.wallMap[p] {
				sb.WriteString("|#####")
			} else if p == raceTrack.startTile {
				sb.WriteString("| SSS ")
			} else if p == raceTrack.endTile {
				sb.WriteString(fmt.Sprintf("|E:%3d", cheatMap[p]))
			} else {
				sb.WriteString(fmt.Sprintf("|%4d ", cheatMap[p]))
			}
		}
		sb.WriteString("|\n")
	}
	sb.WriteString(line)
	sb.WriteString("+")

	return sb.String()
}

func main() {
	var raceTrack *RaceTrack
	initDuration := MeasureRuntime(
		func() {
			raceTrack = loadRaceTrack("input.txt")
		},
	)
	fmt.Println("Init took", initDuration.Microseconds(), "μs")

	part1duration := MeasureRuntime(func() {
		_, part1cheatCount := raceTrack.countCheats(2)
		part1sum := 0
		for k := range part1cheatCount {
			if k >= 100 {
				part1sum += part1cheatCount[k]
			}
		}
		fmt.Println("Part 1:", part1sum)
	})

	fmt.Println("Part 1 took", part1duration.Microseconds(), "μs")

	part2duration := MeasureRuntime(func() {
		_, part2cheatCount := raceTrack.countCheats(20)

		part2sum := 0
		for k := range part2cheatCount {
			if k >= 100 {
				part2sum += part2cheatCount[k]
			}
		}
		fmt.Println("Part 2:", part2sum)
	})

	fmt.Println("Part 2 took", part2duration.Microseconds(), "μs")
}
