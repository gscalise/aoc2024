package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"bautik.net/advent2024/helpers"
)

var PRUNE_BITMASK = (2 << 23) - 1

func calculateNext(seed int) int {
	seed = ((seed << 6) ^ seed) & PRUNE_BITMASK
	seed = (seed ^ (seed >> 5)) & PRUNE_BITMASK
	return (seed ^ (seed << 11)) & PRUNE_BITMASK
}

func fillSequenceMap(seed int, n int, sequenceMap map[[4]int]int) int {
	currentValue := seed
	sequenceDelta := []int{}
	seenSequences := map[[4]int]bool{}
	for range n {
		newValue := calculateNext(currentValue)
		currentValueMod := currentValue % 10
		newValueMod := newValue % 10
		delta := int(newValueMod - currentValueMod)
		sequenceDelta = append(sequenceDelta, int(delta))
		if len(sequenceDelta) > 4 {
			sequenceDelta = sequenceDelta[1:]
		}
		if len(sequenceDelta) == 4 {
			sequenceArray := [4]int(sequenceDelta)
			if _, ok := seenSequences[sequenceArray]; !ok {
				sequenceMap[sequenceArray] += int(newValueMod)
				seenSequences[sequenceArray] = true
			}
		}
		currentValue = newValue
	}
	return currentValue
}

func main() {
	duration := helpers.MeasureRuntime(func() {

		fileBytes, _ := os.ReadFile("input.txt")
		part1sum := 0
		sequenceMap := map[[4]int]int{}
		for _, l := range strings.Split(string(fileBytes), "\n") {
			number, _ := strconv.Atoi(l)
			part1sum += fillSequenceMap(number, 2000, sequenceMap)
		}
		fmt.Println("Part 1:", part1sum)

		maxPurchase := 0
		for _, v := range sequenceMap {
			if v > maxPurchase {
				maxPurchase = v
			}
		}

		fmt.Println("Part 2:", maxPurchase)
	})

	fmt.Println("Took", duration.Microseconds(), "μs")
}
