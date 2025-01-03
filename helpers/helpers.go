package helpers

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/constraints"
)

func ParseInputDay01(input string) ([]int, []int) {
	lines := strings.Split(input, "\n")
	left := []int{}
	right := []int{}
	regex := regexp.MustCompile(`^(\d+)\s+(\d+)$`)
	for _, line := range lines {
		matches := regex.FindStringSubmatch(strings.TrimSpace(line))
		if len(matches) == 3 {
			v0, _ := strconv.Atoi(matches[1])
			v1, _ := strconv.Atoi(matches[2])
			left = append(left, v0)
			right = append(right, v1)
		}
	}
	return left, right
}

type Direction struct {
	DX, DY int
}

var UP Direction = Direction{0, -1}
var DOWN Direction = Direction{0, 1}
var LEFT Direction = Direction{-1, 0}
var RIGHT Direction = Direction{1, 0}

var Directions [4]*Direction = [4]*Direction{&DOWN, &LEFT, &UP, &RIGHT}

type Position struct {
	X, Y int
}

var ZERO_ZERO = Position{X: 0, Y: 0}

func (p Position) Displace(m Direction) Position {
	return Position{
		X: p.X + m.DX,
		Y: p.Y + m.DY,
	}
}

func (d Direction) Opposite() Direction {
	switch d {
	case UP:
		return DOWN
	case DOWN:
		return UP
	case RIGHT:
		return LEFT
	case LEFT:
		return RIGHT
	default:
		panic("invalid")
	}
}

func Abs[T constraints.Integer](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

func (p Position) GetDirectionTo(o Position) Direction {
	if p.X == o.Y {
		if p.Y < o.Y {
			return DOWN
		} else {
			return UP
		}
	} else {
		if p.X < o.X {
			return RIGHT
		} else {
			return LEFT
		}
	}
}

func MeasureRuntime(f func()) time.Duration {
	startTime := time.Now()
	f()
	return time.Since(startTime)
}

func Intersection(a, b []int) []int {
	m := make(map[int]bool)
	for _, v := range b {
		m[v] = true
	}
	var result []int
	for _, v := range a {
		if m[v] {
			result = append(result, v)
		}
	}
	return result
}

func Union(a, b []int) []int {
	m := make(map[int]bool)
	for _, v := range a {
		m[v] = true
	}
	for _, v := range b {
		m[v] = true
	}
	var result []int
	for k := range m {
		result = append(result, k)
	}
	return result
}

func Difference(a, b []int) []int {
	m := make(map[int]bool)
	for _, v := range b {
		m[v] = true
	}
	var result []int
	for _, v := range a {
		if !m[v] {
			result = append(result, v)
		}
	}
	return result
}
