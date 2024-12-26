package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

type Lock struct {
	heights [5]int
}

type Key struct {
	heights [5]int
}

func parseLock(lines []string) Lock {
	return Lock{heights: parseHeights(lines)}
}

func parseHeights(lines []string) (heights [5]int) {
	heights = [5]int{}
	for i, l := range lines[1:] {
		for p, c := range l {
			if c == '.' && lines[i][p] == '#' {
				heights[p] = i
			}
		}
	}
	return heights
}

func parseKey(lines []string) Key {
	slices.Reverse(lines)
	return Key{heights: parseHeights(lines)}
}

func readInput(input string) (locks []Lock, keys []Key) {
	locks = []Lock{}
	keys = []Key{}

	for _, block := range strings.Split(input, "\n\n") {
		blockLines := strings.Split(block, "\n")
		if blockLines[0] == "#####" {
			locks = append(locks, parseLock(blockLines))
		} else {
			keys = append(keys, parseKey(blockLines))
		}
	}

	return locks, keys
}

func (lock Lock) matches(key Key) bool {
	for i := range 5 {
		if lock.heights[i]+key.heights[i] > 5 {
			return false
		}
	}
	return true
}

func main() {
	fileBytes, _ := os.ReadFile("input.txt")
	locks, keys := readInput(string(fileBytes))
	count := 0

	for _, l := range locks {
		for _, k := range keys {
			if l.matches(k) {
				count++
			}
		}
	}
	fmt.Println("Part 1:", count)
}
