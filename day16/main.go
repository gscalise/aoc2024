package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

type Direction struct {
	dx, dy int
}

var UP Direction = Direction{0, -1}
var DOWN Direction = Direction{0, 1}
var LEFT Direction = Direction{-1, 0}
var RIGHT Direction = Direction{1, 0}

var directions [4]*Direction = [4]*Direction{&UP, &DOWN, &LEFT, &RIGHT}

type Position struct {
	x, y int
}

func (p Position) displace(m Direction) Position {
	return Position{
		x: p.x + m.dx,
		y: p.y + m.dy,
	}
}

func (d *Direction) opposite() *Direction {
	switch *d {
	case UP:
		return &DOWN
	case DOWN:
		return &UP
	case RIGHT:
		return &LEFT
	case LEFT:
		return &RIGHT
	default:
		panic("BANG")
	}
}

type Maze struct {
	width, height int
	wallMap       map[Position]bool
	turnMap       map[Position]bool
	tJunctionMap  map[Position]bool
	startTile     Position
	endTile       Position
}

func (m *Maze) addWall(p Position) {
	m.wallMap[p] = true
}

func (m *Maze) String() string {
	var sb strings.Builder

	for r := range m.height {
		for c := range m.width {
			position := Position{x: c, y: r}
			if m.wallMap[position] {
				sb.WriteRune('#')
			} else if position == m.startTile {
				sb.WriteRune('S')
			} else if position == m.endTile {
				sb.WriteRune('E')
			} else if m.turnMap[position] {
				sb.WriteRune('L')
			} else if m.tJunctionMap[position] {
				sb.WriteRune('T')
			} else {
				sb.WriteRune(' ')
			}
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

type ReindeerNode struct {
	position  Position
	direction Direction
}

type Reindeer struct {
	node  ReindeerNode
	path  []Position
	score int
}

func (node ReindeerNode) getNeighbourNodes() []ReindeerNode {
	neighbours := []ReindeerNode{}
	oppositeDirection := node.direction.opposite()

	for _, d := range directions {
		if d == oppositeDirection {
			continue
		}
		position := node.position.displace(*d)
		neighbours = append(neighbours, ReindeerNode{
			position:  position,
			direction: *d,
		})
	}
	return neighbours
}

func newMaze(input string) *Maze {
	lines := strings.Split(input, "\n")
	height := len(lines)
	width := len(lines[0])

	maze := Maze{
		width:   width,
		height:  height,
		wallMap: map[Position]bool{},
	}

	for r, l := range lines {
		for c, o := range l {
			if o == '#' {
				maze.addWall(Position{x: c, y: r})
			}
			if o == 'S' {
				maze.startTile = Position{x: c, y: r}
			}
			if o == 'E' {
				maze.endTile = Position{x: c, y: r}
			}
		}
	}

	return &maze
}

func (m *Maze) findMinimumCost() (minScore int, bestSeatCount int) {
	minScore = math.MaxInt
	queue := []Reindeer{
		{
			node: ReindeerNode{
				position:  m.startTile,
				direction: RIGHT,
			},
			path:  []Position{m.startTile},
			score: 0,
		},
	}

	visitedNodes := map[ReindeerNode]int{}
	bestSeatsMap := map[int][]Position{}

	for len(queue) > 0 {
		currentReindeer := queue[0]
		queue = queue[1:]

		if currentReindeer.score > minScore {
			continue
		}

		if currentReindeer.node.position == m.endTile {
			if currentReindeer.score <= minScore {
				minScore = currentReindeer.score
				bestSeatsMap[minScore] = append(bestSeatsMap[minScore], currentReindeer.path...)
			}
			continue
		}

		for _, n := range currentReindeer.node.getNeighbourNodes() {
			if m.wallMap[n.position] {
				continue
			}
			score := currentReindeer.score + 1

			if currentReindeer.node.direction != n.direction {
				score += 1000
			}
			if visitedNodeScore, ok := visitedNodes[n]; ok {
				if visitedNodeScore < score {
					continue
				}
			}
			visitedNodes[n] = score
			newPath := []Position{}
			newPath = append(newPath, currentReindeer.path...)
			queue = append(queue, Reindeer{
				node:  n,
				path:  append(newPath, n.position),
				score: score,
			})
		}
	}
	bestSeatCountSet := map[Position]bool{}
	for _, position := range bestSeatsMap[minScore] {
		bestSeatCountSet[position] = true
	}
	return minScore, len(bestSeatCountSet)
}

func (p Position) GetDirectionTo(o Position) Direction {
	if p.x == o.x {
		if p.y < o.y {
			return DOWN
		} else {
			return UP
		}
	} else {
		if p.x < o.x {
			return RIGHT
		} else {
			return LEFT
		}
	}
}

func main() {
	fileBytes, _ := os.ReadFile("input.txt")
	maze := newMaze(string(fileBytes))
	part1, part2 := maze.findMinimumCost()
	fmt.Println("Part 1", part1)
	fmt.Println("Part 2", part2)
}
