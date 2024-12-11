package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type Stone uint

func (s Stone) Blink() []Stone {
	if s == 0 {
		return []Stone{1}
	} else {
		lenNum := int(math.Floor(math.Log10(float64(s)))) + 1
		if lenNum%2 == 0 {
			highPartFactor := uint(math.Pow10(lenNum / 2))
			secondHalf := uint(s) % highPartFactor
			firstHalf := (uint(s) - secondHalf) / highPartFactor
			return []Stone{Stone(firstHalf), Stone(secondHalf)}
		} else {
			return []Stone{Stone(s * 2024)}
		}
	}
}

type BlinkLenKey struct {
	stone  Stone
	blinks int
}

var blinkLenMap map[BlinkLenKey]int = make(map[BlinkLenKey]int)

func (s Stone) StoneCountAfterBlinking(times int) int {
	if times == 0 {
		return 1
	}
	memoKey := BlinkLenKey{
		stone:  s,
		blinks: times,
	}
	blinks, ok := blinkLenMap[memoKey]
	if !ok {
		blinkStones := s.Blink()
		for _, bs := range blinkStones {
			blinks += bs.StoneCountAfterBlinking(times - 1)
		}
		blinkLenMap[memoKey] = blinks
	}
	return blinks
}

func main() {
	parseStartTime := time.Now()
	input, _ := os.ReadFile("input.txt")
	stoneStrings := strings.Split(string(input), " ")
	fmt.Println("Parsing Took", time.Since(parseStartTime).Microseconds(), "μs")

	p1startTime := time.Now()
	part1LenSum := 0
	for _, s := range stoneStrings {
		v, _ := strconv.Atoi(s)
		part1LenSum += Stone(v).StoneCountAfterBlinking(25)
	}
	fmt.Println("Part 1:", part1LenSum)
	fmt.Println("Part1 Took", time.Since(p1startTime).Microseconds(), "μs")

	p2startTime := time.Now()
	part2LenSum := 0
	for _, s := range stoneStrings {
		v, _ := strconv.Atoi(s)
		part2LenSum += Stone(v).StoneCountAfterBlinking(75)
	}
	fmt.Println("Part 2:", part2LenSum)
	fmt.Println("Part2 Took", time.Since(p2startTime).Microseconds(), "μs")

	fmt.Println(len(blinkLenMap), "elements in the memoized map")
}
