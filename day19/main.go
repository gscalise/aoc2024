package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	fileBytes, _ := os.ReadFile("input.txt")

	originalPatterns := []string{}
	validPatterns := map[string]int{}
	invalidPatterns := map[string]bool{}
	input := strings.Split(string(fileBytes), "\n\n")

	originalPatterns = append(originalPatterns, strings.Split(input[0], ", ")...)

	var possiblePatterns func(string) int

	possiblePatterns = func(p string) int {
		if invalidPatterns[p] {
			return 0
		}
		if validPatterns[p] > 0 {
			return validPatterns[p]
		}
		candidatePatterns := []string{}
		for _, validPattern := range originalPatterns {
			if len(p) >= len(validPattern) && p[0:len(validPattern)] == validPattern {
				candidatePatterns = append(candidatePatterns, validPattern)
			}
		}
		if len(candidatePatterns) == 0 {
			invalidPatterns[p] = true
			return 0
		} else {
			sumPossible := 0
			for _, candidatePattern := range candidatePatterns {
				if p == candidatePattern {
					sumPossible += 1
				} else {
					sumPossible += possiblePatterns(p[len(candidatePattern):])
				}
			}
			if sumPossible > 0 {
				validPatterns[p] = sumPossible
				return sumPossible
			}
		}
		invalidPatterns[p] = true
		return 0
	}

	count := 0
	validTowels := 0

	for _, towel := range strings.Split(input[1], "\n") {
		towelPatternCount := possiblePatterns(towel)
		count += towelPatternCount
		if towelPatternCount > 0 {
			validTowels++
		}
	}
	fmt.Println("Part 1:", validTowels)
	fmt.Println("Part 2:", count)
}
