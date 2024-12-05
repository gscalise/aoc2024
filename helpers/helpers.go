package helpers

import (
	"regexp"
	"strconv"
	"strings"
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
