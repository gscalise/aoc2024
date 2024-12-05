package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
)

func Abs[T constraints.Integer](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

func isSafeReportPart1(report []int) bool {

	curDif := report[1] - report[0]

	if curDif == 0 || curDif > 3 || curDif < -3 {
		return false
	}

	for i := range report[2:] {
		newDif := report[i+2] - report[i+1]
		if newDif == 0 || newDif*curDif < 0 || newDif > 3 || newDif < -3 {
			return false
		}
	}
	return true
}

func isSafeReportPart2(report []int, skip int) bool {

	reps := report

	if skip >= 0 {
		reps = append([]int{}, report[:skip]...)
		reps = append(reps, report[skip+1:]...)
	}

	curDif := reps[1] - reps[0]

	if curDif == 0 || Abs(curDif) > 3 {
		return skip != -1 || isSafeReportPart2(report, 0) || isSafeReportPart2(report, 1)
	}

	for i := range reps[2:] {
		newDif := reps[i+2] - reps[i+1]
		if newDif*curDif <= 0 || Abs(newDif) > 3 {
			return skip != -1 ||
				isSafeReportPart2(report, i) ||
				isSafeReportPart2(report, i+1) ||
				isSafeReportPart2(report, i+2)
		}
	}
	return true
}

func main() {

	inputName := "input.txt"

	fileBytes, err := os.ReadFile(inputName)
	if err != nil {
		panic("Failed to read file")
	}
	inputLines := strings.Split(string(fileBytes), "\n")

	rows := [][]int{}

	for i := range inputLines {
		row := []int{}
		lineSplit := strings.Split(inputLines[i], " ")
		for _, x := range lineSplit {
			v, _ := strconv.Atoi(x)
			row = append(row, v)
		}
		rows = append(rows, row)
	}

	count := 0

	for _, report := range rows {
		if isSafeReportPart1(report) {
			count += 1
		}
	}

	fmt.Println(count)

	count = 0

	for _, report := range rows {
		if isSafeReportPart2(report, -1) {
			count += 1
		}
	}

	fmt.Println(count)
}
