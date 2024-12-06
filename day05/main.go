package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func checkOrder(rules map[string]map[string]bool, pages []string) bool {
	for pIdx, page := range pages {
		pageRules := rules[page]
		if pIdx == len(pages)-1 || len(pageRules) == 0 {
			continue
		}
		for spIdx := range pages[pIdx+1:] {
			if pageRules[pages[pIdx+spIdx+1]] {
				return false
			}
		}
	}
	return true
}

func getFixedMidNum(rules map[string]map[string]bool, pages []string) int {
	for pIdx, page := range pages[:len(pages)/2+1] {
		pageRules := rules[page]
		if len(pageRules) == 0 {
			continue
		}
		for spIdx := range pages[pIdx+1:] {
			if pageRules[pages[pIdx+spIdx+1]] {
				pages[pIdx], pages[pIdx+spIdx+1] = pages[pIdx+spIdx+1], pages[pIdx]
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
	rules := map[string]map[string]bool{}
	for _, l := range inputLines {
		rule := ruleRE.FindStringSubmatch(l)
		if len(rule) == 3 {
			_, ok := rules[rule[2]]
			if !ok {
				rules[rule[2]] = map[string]bool{}
			}
			rules[rule[2]][rule[1]] = true
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
