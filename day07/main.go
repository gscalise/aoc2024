package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type InputLine struct {
	target int
	values []int
}

func numSolutions(currentResult int, valueIdx int, line InputLine) int {
	if currentResult > line.target {
		return 0
	}

	resultSum := 0
	if currentResult*line.values[valueIdx] == line.target {
		if valueIdx+1 == len(line.values) {
			resultSum += 1
		} else {
			resultSum += numSolutions(currentResult*line.values[valueIdx], valueIdx+1, line)
		}
	}
	if currentResult*line.values[valueIdx] < line.target {
		if valueIdx+1 < len(line.values) {
			resultSum += numSolutions(currentResult*line.values[valueIdx], valueIdx+1, line)
		}
	}

	if currentResult+line.values[valueIdx] == line.target {
		if valueIdx+1 == len(line.values) {
			resultSum += 1
		} else {
			resultSum += numSolutions(currentResult+line.values[valueIdx], valueIdx+1, line)
		}
	}
	if currentResult+line.values[valueIdx] < line.target {
		if valueIdx+1 < len(line.values) {
			resultSum += numSolutions(currentResult+line.values[valueIdx], valueIdx+1, line)
		}
	}

	return resultSum
}

func concatIntValues(v1, v2 int) int {
	return v1*int(math.Pow10(1+int(math.Log10(float64(v2))))) + v2
}

func numSolutionsP2(currentResult int, valueIdx int, line InputLine) int {
	if currentResult > line.target {
		return 0
	}

	resultSum := 0

	concatResult := concatIntValues(currentResult, line.values[valueIdx])

	if concatResult == line.target {
		if valueIdx+1 == len(line.values) {
			resultSum += 1
		} else {
			resultSum += numSolutionsP2(concatResult, valueIdx+1, line)
		}
	}
	if concatResult < line.target {
		if valueIdx+1 < len(line.values) {
			resultSum += numSolutionsP2(concatResult, valueIdx+1, line)
		}
	}

	if currentResult*line.values[valueIdx] == line.target {
		if valueIdx+1 == len(line.values) {
			resultSum += 1
		} else {
			resultSum += numSolutionsP2(currentResult*line.values[valueIdx], valueIdx+1, line)
		}
	}
	if currentResult*line.values[valueIdx] < line.target {
		if valueIdx+1 < len(line.values) {
			resultSum += numSolutionsP2(currentResult*line.values[valueIdx], valueIdx+1, line)
		}
	}

	if currentResult+line.values[valueIdx] == line.target {
		if valueIdx+1 == len(line.values) {
			resultSum += 1
		} else {
			resultSum += numSolutionsP2(currentResult+line.values[valueIdx], valueIdx+1, line)
		}
	}
	if currentResult+line.values[valueIdx] < line.target {
		if valueIdx+1 < len(line.values) {
			resultSum += numSolutionsP2(currentResult+line.values[valueIdx], valueIdx+1, line)
		}
	}

	return resultSum
}

func parseInputLines(lines []string) []InputLine {
	resultLines := []InputLine{}
	for _, l := range lines {
		parts := strings.Split(l, ":")
		target, _ := strconv.Atoi(parts[0])
		values := []int{}
		for _, v := range strings.Split(strings.TrimSpace(parts[1]), " ") {
			vint, _ := strconv.Atoi(v)
			values = append(values, vint)
		}
		result := InputLine{
			target: target,
			values: values,
		}
		resultLines = append(resultLines, result)
	}
	return resultLines
}

func main() {
	input, _ := os.ReadFile("input.txt")
	inputLines := strings.Split(string(input), "\n")
	totalSumP1 := 0
	totalSumP2 := 0
	for _, l := range parseInputLines(inputLines) {
		if numSolutions(l.values[0], 1, l) > 0 {
			totalSumP1 += l.target
		}
		if numSolutionsP2(l.values[0], 1, l) > 0 {
			totalSumP2 += l.target
		}
	}
	fmt.Println("Part1:", totalSumP1)
	fmt.Println("Part2:", totalSumP2)
}
