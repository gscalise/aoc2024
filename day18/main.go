package main

import (
	"fmt"
	"os"
	"strings"

	"bautik.net/advent2024/helpers"
)

type Position = helpers.Position
type Direction = helpers.Direction

type RAMMap struct {
	width, height     int
	corruptedBlocks   []Position
	corruptedBlockMap map[Position]int
}

func asInputCoordinates(p Position) string {
	return fmt.Sprintf("%d,%d", p.X, p.Y)
}

func newRAMMap(input string, width, height int) (ramMap *RAMMap) {
	corruptedBlockMap := map[Position]int{}
	corruptedBlocks := []Position{}
	ramMap = &RAMMap{
		width,
		height,
		corruptedBlocks,
		corruptedBlockMap,
	}
	for i, b := range strings.Split(input, "\n") {
		var x int
		var y int
		fmt.Sscanf(b, "%d,%d", &x, &y)
		blockPosition := Position{X: x, Y: y}
		corruptedBlockMap[blockPosition] = i + 1
		ramMap.corruptedBlocks = append(ramMap.corruptedBlocks, blockPosition)
	}
	return ramMap
}

func (ramMap *RAMMap) isEndPosition(p Position) bool {
	return p.X == ramMap.width-1 && p.Y == ramMap.height-1
}

func (ramMap *RAMMap) isValidPosition(p Position, n int) bool {
	ns, corrupted := ramMap.corruptedBlockMap[p]
	return p.X >= 0 && p.Y >= 0 && p.X <= ramMap.width-1 && p.Y <= ramMap.height-1 && (!corrupted || n < ns)
}

func (ramMap *RAMMap) getMapString(ns int) string {
	var sb strings.Builder
	for r := range ramMap.height {
		for c := range ramMap.width {
			if ramMap.isValidPosition(helpers.Position{X: c, Y: r}, ns) {
				sb.WriteRune('.')
			} else {
				sb.WriteRune('#')
			}
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

type TripNode struct {
	position Position
	path     []Position
}

func (ramMap *RAMMap) getEndPosition() Position {
	return Position{
		X: ramMap.width - 1,
		Y: ramMap.height - 1,
	}
}

func (ramMap *RAMMap) getNeighbors(p Position, time int) (neighbors []Position) {
	neighbors = []Position{}
	for _, d := range helpers.Directions {
		pd := p.Displace(*d)
		if ramMap.isValidPosition(pd, time) {
			neighbors = append(neighbors, pd)
		}
	}
	return neighbors
}

func solve(ramMap *RAMMap, time int) (pathLen int, pathFound bool) {
	endPosition := ramMap.getEndPosition()
	visited := map[Position]bool{}
	pathLen = 0
	queue := []Position{
		helpers.ZERO_ZERO,
	}
	for len(queue) > 0 {
		levelSize := len(queue)

		for range levelSize {
			current := queue[0]
			queue = queue[1:]
			if current == endPosition {
				return pathLen, true
			}
			for _, neighbour := range ramMap.getNeighbors(current, time) {
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

func loadRAMMap(file string, width int, height int) *RAMMap {
	fileBytes, _ := os.ReadFile(file)
	ramMap := newRAMMap(string(fileBytes), width, height)
	return ramMap
}

func main() {
	ramMap := loadRAMMap("input.txt", 71, 71)
	pathLen, found := solve(ramMap, 1024)
	if !found {
		panic("Solution not found!")
	}
	fmt.Println("Part 1:", pathLen)

	lowIndex := 0
	maxIndex := len(ramMap.corruptedBlocks)
	highIndex := maxIndex
	index := 0
	foundBlock := false
	for lowIndex <= highIndex && index+1 < maxIndex && !foundBlock {
		index = (highIndex + lowIndex) / 2

		_, found := solve(ramMap, index)
		if !found { // we overshot
			highIndex = index
		} else {
			_, plusOneFound := solve(ramMap, index+1)
			if !plusOneFound {
				fmt.Println("Part 2:", asInputCoordinates(ramMap.corruptedBlocks[index]), "at index", index)
				foundBlock = true
			} else {
				lowIndex = index
			}
		}
	}
	if !foundBlock {
		fmt.Println("Part 2: No solution found")
	}
}
