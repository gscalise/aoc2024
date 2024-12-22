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

type SequenceResult struct {
	secret      int
	sequenceMap map[[4]int]int
}

func processSeller(retValCh chan SequenceResult, seed int, n int) {
	currentValue := seed
	sequenceDelta := []int{}
	sequenceMap := map[[4]int]int{}
	for range n {
		newValue := calculateNext(currentValue)
		currentValueMod := currentValue % 10
		newValueMod := newValue % 10
		delta := newValueMod - currentValueMod
		sequenceDelta = append(sequenceDelta, int(delta))
		if len(sequenceDelta) > 4 {
			sequenceDelta = sequenceDelta[1:]
		}
		if len(sequenceDelta) == 4 {
			sequenceArray := [4]int(sequenceDelta)
			if _, ok := sequenceMap[sequenceArray]; !ok {
				sequenceMap[sequenceArray] = newValueMod
			}
		}
		currentValue = newValue
	}
	retValCh <- SequenceResult{secret: currentValue, sequenceMap: sequenceMap}
}

func main() {
	duration := helpers.MeasureRuntime(func() {

		fileBytes, _ := os.ReadFile("input.txt")
		part1sum := 0
		sequenceMap := map[[4]int]int{}
		sellerSeeds := []int{}
		retValCh := make(chan SequenceResult)
		for _, l := range strings.Split(string(fileBytes), "\n") {
			seed, _ := strconv.Atoi(l)
			sellerSeeds = append(sellerSeeds, seed)
			go processSeller(retValCh, seed, 2000)
		}

		for range sellerSeeds {
			result := <-retValCh
			part1sum += result.secret
			for k, v := range result.sequenceMap {
				sequenceMap[k] += v
			}
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
