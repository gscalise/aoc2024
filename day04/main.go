package main

import (
	"fmt"
	"os"
	"strings"
)

func main1() {
	fileBytes, _ := os.ReadFile("input.txt")
	xmasString := "XMAS"
	searchDirections := [][]int{
		{-1, -1},
		{-1, 0},
		{-1, 1},
		{0, -1},
		{0, 1},
		{1, -1},
		{1, 0},
		{1, 1},
	}
	fileText := string(fileBytes)
	wordSearch := strings.Split(fileText, "\n")
	count := 0
	for ri, r := range wordSearch {
		for ci, c := range r {
			if c == rune(xmasString[0]) {
			direction:
				for _, d := range searchDirections {
					for xi, xc := range xmasString[1:] {
						nr := ri + (xi+1)*d[0]
						nc := ci + (xi+1)*d[1]
						if nr < 0 || nr >= len(wordSearch) || nc < 0 || nc >= len(wordSearch[0]) || rune(wordSearch[nr][nc]) != xc {
							continue direction
						}
					}
					count += 1
				}
			}
		}
	}
	fmt.Println(count)
}

func main() {
	fileBytes, _ := os.ReadFile("input.txt")
	fileText := string(fileBytes)
	wordSearch := strings.Split(fileText, "\n")
	count := 0
	for ri, r := range wordSearch {
		if ri == 0 || ri == len(wordSearch)-1 {
			continue
		}
		for ci, c := range r {
			if ci == 0 || ci == len(wordSearch[0])-1 {
				continue
			}
			if c == 'A' {
				mc := 0
				if wordSearch[ri-1][ci-1] == 'M' && wordSearch[ri+1][ci+1] == 'S' {
					mc += 1
				}
				if wordSearch[ri-1][ci-1] == 'S' && wordSearch[ri+1][ci+1] == 'M' {
					mc += 1
				}
				if wordSearch[ri-1][ci+1] == 'M' && wordSearch[ri+1][ci-1] == 'S' {
					mc += 1
				}
				if wordSearch[ri-1][ci+1] == 'S' && wordSearch[ri+1][ci-1] == 'M' {
					mc += 1
				}
				if mc == 2 {
					count += 1
				}
			}
		}
	}
	fmt.Println(count)
}
