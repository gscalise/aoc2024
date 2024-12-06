package main

import (
	"fmt"
	"os"
	"strings"
)

var directions = [4][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

func newPosition(currentPosition [2]int, direction [2]int) [2]int {
	return [2]int{
		currentPosition[0] + direction[0],
		currentPosition[1] + direction[1],
	}
}

func nextToVisitedPosition(rows int, columns int, position [2]int, visitedPositions map[[2]int]bool) bool {
	if visitedPositions[position] {
		return true
	}

	for _, d := range directions {
		p := [2]int{
			position[0] + d[0],
			position[1] + d[1],
		}
		if p[0] < 0 || p[0] == rows || p[1] < 0 || p[1] == columns {
			break
		}
		if visitedPositions[p] {
			return true
		}
	}
	return false
}

func main() {
	fileBytes, _ := os.ReadFile("input.txt")
	fileText := string(fileBytes)
	inputLines := strings.Split(fileText, "\n")
	currentDirection := 0
	var currentPosition [2]int
	var initialPosition [2]int
	visitedPositions := map[[2]int]bool{}
	nr := len(inputLines)
	nc := len(inputLines[0])

	for r, line := range inputLines {
		c := strings.Index(line, "^")
		if c != -1 {
			currentPosition = [2]int{r, c}
			initialPosition = currentPosition
			visitedPositions[currentPosition] = true
			break
		}
	}

	for {
		newPos := newPosition(currentPosition, directions[currentDirection])
		if newPos[0] >= nr || newPos[0] < 0 || newPos[1] >= nc || newPos[1] < 0 {
			break
		} else if inputLines[newPos[0]][newPos[1]] == '#' {
			currentDirection += 1
			currentDirection %= 4
		} else {
			currentPosition = newPos
			visitedPositions[currentPosition] = true
		}
	}

	fmt.Println("Part1:", len(visitedPositions))
	blockCount := 0
	for pos := range visitedPositions {
		r := pos[0]
		c := pos[1]
		visitedPaths := map[[3]int]bool{}
		currentPosition = initialPosition
		currentDirection = 0
		visitedPaths[[3]int{currentPosition[0], currentPosition[1], currentDirection}] = true
		for {
			newPos := newPosition(currentPosition, directions[currentDirection])
			if newPos[0] >= nr || newPos[0] < 0 || newPos[1] >= nc || newPos[1] < 0 {
				break
			} else if (newPos[0] == r && newPos[1] == c) || inputLines[newPos[0]][newPos[1]] == '#' {
				currentDirection += 1
				currentDirection %= 4
			} else {
				currentPosition = newPos
				currentPath := [3]int{currentPosition[0], currentPosition[1], currentDirection}
				_, ok := visitedPaths[currentPath]
				if ok {
					blockCount += 1
					break
				}
				visitedPaths[currentPath] = true
			}
		}
	}
	fmt.Println("Part2:", blockCount)
}
