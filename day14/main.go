package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var robotDataRe = regexp.MustCompile(`^p=(\d+),(\d+) v=(-?\d+),(-?\d+)$`)

type Position struct {
	x, y int
}

type Velocity struct {
	vx, vy int
}

type RobotData struct {
	position Position
	velocity Velocity
}

type RobotArena struct {
	width, height int
	robots        []*RobotData
}

func newArena(width, height int) *RobotArena {
	robots := []*RobotData{}
	return &RobotArena{
		width,
		height,
		robots,
	}
}

func getRobotData(line string) *RobotData {
	match := robotDataRe.FindStringSubmatch(line)
	x, _ := strconv.Atoi(match[1])
	y, _ := strconv.Atoi(match[2])
	vx, _ := strconv.Atoi(match[3])
	vy, _ := strconv.Atoi(match[4])
	return newRobot(x, y, vx, vy)
}

func newRobot(x int, y int, vx int, vy int) *RobotData {
	return &RobotData{
		position: Position{x, y},
		velocity: Velocity{vx, vy},
	}
}

func (r RobotData) PositionAfterSeconds(n int, arena *RobotArena) Position {
	x := (r.position.x + n*r.velocity.vx) % arena.width
	y := (r.position.y + n*r.velocity.vy) % arena.height
	if x < 0 {
		x += arena.width
	}

	if y < 0 {
		y += arena.height
	}

	return Position{x, y}
}

func (a RobotArena) getQuadrant(p Position) int {
	middleColumn := a.width / 2
	middleRow := a.height / 2
	quadrant := 0
	if p.x == middleColumn || p.y == middleRow {
		return -1
	} else {
		if p.x > middleColumn {
			quadrant += 2
		}
		if p.y > middleRow {
			quadrant += 1
		}
	}
	return quadrant
}

func (arena RobotArena) getFactor(seconds int) int {
	quadrantCounts := map[int]int{}
	for _, r := range arena.robots {
		q := arena.getQuadrant(r.PositionAfterSeconds(seconds, &arena))
		if q != -1 {
			quadrantCounts[q] += 1
		}
	}

	factor := 1
	for _, q := range quadrantCounts {
		factor *= q
	}
	return factor
}

func doPart1(fileName string, width, height int) {
	startTime := time.Now()
	input, _ := os.ReadFile(fileName)
	arena := newArena(width, height)
	for _, l := range strings.Split(string(input), "\n") {
		arena.robots = append(arena.robots, getRobotData(l))
	}
	factor := arena.getFactor(100)

	fmt.Println(factor)

	fmt.Println("Took", time.Since(startTime).Microseconds(), "μs")
}

func (a RobotArena) DrawArena(seconds int) {
	positions := map[Position]bool{}
	for _, r := range a.robots {
		positions[r.PositionAfterSeconds(seconds, &a)] = true
	}
	for y := range a.height {
		for x := range a.width {
			if positions[Position{x, y}] {
				fmt.Print("X")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func doPart2(fileName string, width, height int) {
	startTime := time.Now()
	input, _ := os.ReadFile(fileName)
	arena := newArena(width, height)
	for _, l := range strings.Split(string(input), "\n") {
		arena.robots = append(arena.robots, getRobotData(l))
	}
	minFactor := arena.getFactor(0)
	maxFactor := minFactor
	secondMinFactor := minFactor
	minFactorIdx := 0
	factorSum := 0
	for i := range 100000 {
		factor := arena.getFactor(i)
		factorSum += factor
		if factor < minFactor {
			secondMinFactor = minFactor
			minFactor = factor
			minFactorIdx = i
		}
		if factor > maxFactor {
			maxFactor = factor
		}
	}
	fmt.Println("Min Factor:", minFactor, "Second Min Factor:", secondMinFactor, "Max Factor:", maxFactor, "AvgFactor:", int(float64(factorSum)/10000))

	// arena.DrawArena(minFactorIdx)
	fmt.Println("Idx:", minFactorIdx)
	fmt.Println("Took", time.Since(startTime).Microseconds(), "μs")
}

func main() {
	// demo
	// doDay("demo-input.txt", 11, 7)

	doPart2("input.txt", 101, 103)

}
