package main

import (
	"fmt"
	"os"
	"strings"
)

type Position [2]int
type Direction [2]int

var directions = [4]Direction{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

func newPosition(currentPosition Position, direction Direction) Position {
	return Position{
		currentPosition[0] + direction[0],
		currentPosition[1] + direction[1],
	}
}

func main() {
	fileBytes, _ := os.ReadFile("input.txt")
	fileText := string(fileBytes)
	inputLines := strings.Split(fileText, "\n")
	currentDirection := 0
	var currentPosition Position
	var initialPosition Position
	visitedPositions := map[Position]bool{}
	nr := len(inputLines)
	nc := len(inputLines[0])

	for r, line := range inputLines {
		c := strings.Index(line, "^")
		if c != -1 {
			currentPosition = Position{r, c}
			initialPosition = currentPosition
			visitedPositions[currentPosition] = true
			break
		}
	}

	for {
		newPos := newPosition(currentPosition, directions[currentDirection])
		newPosR, newPosC := newPos[0], newPos[1]

		if newPosR >= nr || newPosR < 0 || newPosC >= nc || newPosC < 0 {
			break
		} else if inputLines[newPosR][newPosC] == '#' {
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
		r, c := pos[0], pos[1]
		visitedPaths := map[[3]int]bool{}
		currentPosition = initialPosition
		currentDirection = 0
		visitedPaths[[3]int{currentPosition[0], currentPosition[1], currentDirection}] = true
		for {
			newPos := newPosition(currentPosition, directions[currentDirection])
			newPosR, newPosC := newPos[0], newPos[1]

			if newPosR >= nr || newPosR < 0 || newPosC >= nc || newPosC < 0 {
				break
			} else if (newPosR == r && newPosC == c) || inputLines[newPosR][newPosC] == '#' {
				currentDirection += 1
				currentDirection %= 4
			} else {
				currentPath := [3]int{newPosR, newPosC, currentDirection}
				_, ok := visitedPaths[currentPath]
				if ok {
					blockCount += 1
					break
				}
				currentPosition = newPos
				visitedPaths[currentPath] = true
			}
		}
	}
	fmt.Println("Part2:", blockCount)
}
