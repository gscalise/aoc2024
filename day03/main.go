package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func mainPart1() {
	input := `xmul(2,4)%&mul[3,7]!@^do_not_mul(5,5)+mul(32,64]then(mul(11,8)mul(8,5))`
	total := getMulTotal(input)
	fmt.Println(total)
	fileBytes, err := os.ReadFile("input.txt")
	if err != nil {
		panic("Failed to read file")
	}
	input = string(fileBytes)
	total = getMulTotal(input)
	fmt.Println(total)
}

func main() {
	input := `xmul(2,4)&mul[3,7]!^don't()_mul(5,5)+mul(32,64](mul(11,8)undo()?mul(8,5))`
	fmt.Println(getMulTotal(filterDonts(input)))
	fileBytes, err := os.ReadFile("input.txt")
	if err != nil {
		panic("Failed to read file")
	}
	input = string(fileBytes)
	fmt.Println(getMulTotal(filterDonts(input)))
}

func filterDonts(input string) string {
	out := ""
	doSubstrings := strings.Split(input, "do()")
	for _, s := range doSubstrings {
		dontIndex := strings.Index(s, "don't()")
		if dontIndex > -1 {
			out += s[0:dontIndex]
		} else {
			out += s
		}
	}
	return out
}

func getMulTotal(input string) int {
	mulRegex := regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)`)
	matches := mulRegex.FindAllStringSubmatch(input, -1)
	total := 0
	for _, m := range matches {
		v0, _ := strconv.Atoi(m[1])
		v1, _ := strconv.Atoi(m[2])
		total += v0 * v1
	}
	return total
}
