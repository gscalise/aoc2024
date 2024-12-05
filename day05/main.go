package main

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

func checkOrder(rules map[string][]string, pages []string) bool {
	for pIdx, page := range pages {
		pageRules := rules[page]
		if pIdx == len(pages)-1 || len(pageRules) == 0 {
			continue
		}
		for _, must := range pageRules {
			if slices.Contains(pages[pIdx+1:], must) {
				return false
			}
		}
	}
	return true
}

func getFixedMidNum(rules map[string][]string, pages []string) int {
	for pIdx, page := range pages[:len(pages)/2+1] {
		pageRules := rules[page]
		if pIdx == len(pages)-1 || len(pageRules) == 0 {
			continue
		}
		for _, must := range pageRules {
			idxFound := slices.Index(pages[pIdx+1:], must)
			if idxFound != -1 {
				pages[pIdx] = must
				pages[pIdx+idxFound+1] = page
				return getFixedMidNum(rules, pages)
			}
		}
	}
	midNum, _ := strconv.Atoi(pages[len(pages)/2])
	return midNum
}

func main() {
	fileBytes, _ := os.ReadFile("input.txt")
	fileText := string(fileBytes)
	inputLines := strings.Split(fileText, "\n")
	ruleRE := regexp.MustCompile(`^(\d+)\|(\d+)$`)
	sumP1 := 0
	sumP2 := 0
	rules := map[string][]string{}
	for _, l := range inputLines {
		rule := ruleRE.FindStringSubmatch(l)
		if len(rule) == 3 {
			rules[rule[2]] = append(rules[rule[2]], rule[1])
		} else if strings.Contains(l, ",") {
			pages := strings.Split(l, ",")
			if checkOrder(rules, pages) {
				midNum, _ := strconv.Atoi(pages[len(pages)/2])
				sumP1 += midNum
			} else {
				sumP2 += getFixedMidNum(rules, pages)
			}
		}
	}
	fmt.Println("Part1:", sumP1)
	fmt.Println("Part2:", sumP2)
}
