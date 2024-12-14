package main

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ChallengeTarget struct {
	x, y int
}

type ButtonAction struct {
	dx, dy int
}

type PrizeChallenge struct {
	aButtonAction ButtonAction
	bButtonAction ButtonAction
	target        ChallengeTarget
}

func getButtonAction(str string, re *regexp.Regexp) ButtonAction {
	matches := re.FindStringSubmatch(str)
	dx, _ := strconv.Atoi(matches[1])
	dy, _ := strconv.Atoi(matches[2])
	return ButtonAction{
		dx,
		dy,
	}
}

func getChallengeTarget(str string) ChallengeTarget {
	matches := prizeRe.FindStringSubmatch(str)
	x, _ := strconv.Atoi(matches[1])
	y, _ := strconv.Atoi(matches[2])
	return ChallengeTarget{
		x,
		y,
	}
}

func (loc *ChallengeTarget) LargerThan(target ChallengeTarget) bool {
	return loc.x > target.x || loc.y > target.y
}

func (target *ChallengeTarget) Press(n int, action ButtonAction) {
	target.x += action.dx * n
	target.y += action.dy * n
}

func (challenge PrizeChallenge) minTokensNaive() int {
	pressCounts := [][2]int{}
	for aPresses := range 100 {
		game := ChallengeTarget{}
		game.Press(aPresses, challenge.aButtonAction)
		if game.LargerThan(challenge.target) {
			break
		}
		if game == challenge.target {
			pressCounts = append(pressCounts, [2]int{aPresses, 0})
		}
		for bPresses := range 100 {
			gameCopy := game
			gameCopy.Press(bPresses, challenge.bButtonAction)
			if gameCopy.LargerThan(challenge.target) {
				break
			}
			if gameCopy == challenge.target {
				pressCounts = append(pressCounts, [2]int{aPresses, bPresses})
			}
		}
	}
	if len(pressCounts) == 0 {
		return -1
	}
	minTokens := 404
	for _, pressCount := range pressCounts {
		attemptCost := pressCount[0]*3 + pressCount[1]
		minTokens = min(minTokens, attemptCost)
	}
	return minTokens
}

func (challenge PrizeChallenge) minTokens(delta int) int {
	bAx := challenge.aButtonAction.dx
	bAy := challenge.aButtonAction.dy

	bBx := challenge.bButtonAction.dx
	bBy := challenge.bButtonAction.dy

	tX := challenge.target.x + delta
	tY := challenge.target.y + delta

	ftX := float64(tX)
	ftY := float64(tY)

	iDeterminant := bAx*bBy - bBx*bAy

	if iDeterminant == 0 {
		return -1
	}
	determinant := float64(iDeterminant)
	invMatrix := [2][2]float64{
		{float64(bBy) / determinant, -float64(bBx) / determinant},
		{-float64(bAy) / determinant, float64(bAx) / determinant},
	}
	bApresses := int(math.Round(invMatrix[0][0]*ftX + invMatrix[0][1]*ftY))
	bBpresses := int(math.Round(invMatrix[1][0]*ftX + invMatrix[1][1]*ftY))
	if (bApresses*bAx+bBpresses*bBx) == tX && (bApresses*bAy+bBpresses*bBy) == tY {
		return bApresses*3 + bBpresses
	}
	return -1
}

var buttonAre = regexp.MustCompile(`^Button A: X\+(\d+), Y\+(\d+)$`)
var buttonBre = regexp.MustCompile(`^Button B: X\+(\d+), Y\+(\d+)$`)
var prizeRe = regexp.MustCompile(`^Prize: X=(\d+), Y=(\d+)$`)

func parseInput(input string) []PrizeChallenge {
	output := []PrizeChallenge{}
	splitInput := strings.Split(input, "\n\n")
	for _, inputStr := range splitInput {
		strLines := strings.Split(inputStr, "\n")
		aButtonAction := getButtonAction(strLines[0], buttonAre)
		bButtonAction := getButtonAction(strLines[1], buttonBre)
		target := getChallengeTarget(strLines[2])
		output = append(output, PrizeChallenge{
			aButtonAction,
			bButtonAction,
			target,
		})
	}
	return output
}

func processFile(f string) {
	startTime := time.Now()
	input, _ := os.ReadFile(f)
	challenges := parseInput(string(input))
	fmt.Println("Found", len(challenges))
	totalCostA := 0
	totalCostB := 0

	for _, c := range challenges {
		tokensA := c.minTokens(0)
		tokensB := c.minTokens(10000000000000)
		if tokensA != -1 {
			totalCostA += tokensA
		}
		if tokensB != -1 {
			totalCostB += tokensB
		}
	}
	fmt.Println("Tokens Part 1: ", totalCostA)
	fmt.Println("Tokens Part 2: ", totalCostB)
	fmt.Println("Took", time.Since(startTime).Microseconds(), "μs")
}

func main() {
	processFile("input.txt")
}
