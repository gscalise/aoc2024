package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var PRUNE_BITMASK = (2 << 23) - 1

func calculateNext(seed int) int {
	seed = ((seed << 6) ^ seed) & PRUNE_BITMASK
	seed = (seed ^ (seed >> 5)) & PRUNE_BITMASK
	return (seed ^ (seed << 11)) & PRUNE_BITMASK
}

func fillSellerSequenceMap(seed int, n int, sequenceMap map[int]map[[4]int]int) int {
	currentValue := seed
	sequenceDelta := []int{}
	sequenceMap[seed] = map[[4]int]int{}
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
			if _, ok := sequenceMap[seed][sequenceArray]; !ok {
				sequenceMap[seed][sequenceArray] =
					int(newValueMod)
			}
		}
		currentValue = newValue
	}
	return currentValue
}

func main() {
	fileBytes, _ := os.ReadFile("input.txt")
	part1sum := 0
	sellerSequenceMap := map[int]map[[4]int]int{}
	for _, l := range strings.Split(string(fileBytes), "\n") {
		number, _ := strconv.Atoi(l)
		part1sum += fillSellerSequenceMap(number, 2000, sellerSequenceMap)
	}
	fmt.Println("Part 1:", part1sum)

	sequenceSet := map[[4]int]bool{}
	for _, m := range sellerSequenceMap {
		for k := range m {
			sequenceSet[k] = true
		}
	}

	maxPurchase := 0
	for sequence := range sequenceSet {
		sequenceSum := 0
		for _, sellerSequences := range sellerSequenceMap {
			sequenceSum += sellerSequences[sequence]
		}
		if sequenceSum > maxPurchase {
			maxPurchase = sequenceSum
		}
	}

	fmt.Println("Part 2:", maxPurchase)

}
