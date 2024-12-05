package main

import (
	"fmt"
	"os"
	"slices"

	"bautik.net/advent2024/helpers"
	"golang.org/x/exp/constraints"
)

func Abs[T constraints.Integer](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	input := `3   4
	4   3
	2   5
	1   3
	3   9
	3   3`

	if true {
		fileBytes, err := os.ReadFile("input.txt")
		if err != nil {
			panic("Failed to read file")
		}
		input = string(fileBytes)
	}
	left, right := helpers.ParseInputDay01(input)

	fmt.Println("Part 1: ", day01processPart1(left, right))
	fmt.Println("Part 2: ", day01processPart2(left, right))
}

func day01processPart1(left []int, right []int) int {
	slices.Sort(left)
	slices.Sort(right)

	sum := 0

	for i := range left {
		sum += Abs(left[i] - right[i])
	}

	return sum
}

func day01processPart2(left []int, right []int) int {

	rightAppearances := map[int]int{}
	for i := range right {
		r := right[i]
		val := rightAppearances[r]
		rightAppearances[r] = val + 1
	}
	sum := 0

	for i := range left {
		sum += Abs(left[i] * rightAppearances[left[i]])
	}

	return sum
}

// func parseInput(input string) ([]int, []int) {
// 	lines := strings.Split(input, "\n")
// 	left := []int{}
// 	right := []int{}
// 	regex := regexp.MustCompile(`^(\d+)\s+(\d+)$`)
// 	for line := range lines {
// 		matches := regex.FindStringSubmatch(strings.TrimSpace(lines[line]))
// 		if len(matches) == 3 {
// 			v0, _ := strconv.Atoi(matches[1])
// 			v1, _ := strconv.Atoi(matches[2])
// 			left = append(left, v0)
// 			right = append(right, v1)
// 		}
// 	}
// 	return left, right
// }
