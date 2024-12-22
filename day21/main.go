package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"bautik.net/advent2024/helpers"
)

var numericMoveMap = map[rune]map[rune]string{
	'0': {
		'A': ">",
		'1': "^<",
		'2': "^",
		'3': "^>",
		'4': "^^<",
		'5': "^^",
		'6': "^^>",
		'7': "^^^<",
		'8': "^^^",
		'9': "^^^>",
	},
	'A': {
		'0': "<",
		'1': "^<<",
		'2': "^<",
		'3': "^",
		'4': "^^<<",
		'5': "<^^",
		'6': "^^",
		'7': "^^^<<",
		'8': "^^^<",
		'9': "^^^",
	},
	'1': {
		'A': ">>v",
		'0': ">v",
		'2': ">",
		'3': ">>",
		'4': "^",
		'5': "^>",
		'6': "^>>",
		'7': "^^",
		'8': "^^>",
		'9': ">>^^",
	},
	'2': {
		'A': "v>",
		'0': "v",
		'1': "<",
		'3': ">",
		'4': "^<",
		'5': "^",
		'6': "^>",
		'7': "^^<",
		'8': "^^",
		'9': ">^^",
	},
	'3': {
		'A': "v",
		'0': "v<",
		'1': "<<",
		'2': "<",
		'4': "^<<",
		'5': "^<",
		'6': "^",
		'7': "<<^^",
		'8': "<^^",
		'9': "^^",
	},
	'4': {
		'A': ">>vv",
		'0': ">vv",
		'1': "v",
		'2': "v>",
		'3': "v>>",
		'5': ">",
		'6': ">>",
		'7': "^",
		'8': "^>",
		'9': "^>>",
	},
	'5': {
		'A': "vv>",
		'0': "vv",
		'1': "v<",
		'2': "v",
		'3': "v>",
		'4': "<",
		'6': ">",
		'7': "^<",
		'8': "^",
		'9': "^>",
	},
	'6': {
		'A': "vv",
		'0': "<vv",
		'1': "v<<",
		'2': "v<",
		'3': "v",
		'4': "<<",
		'5': "<",
		'7': "<<^",
		'8': "<^",
		'9': "^",
	},
	'7': {
		'A': ">>vvv",
		'0': ">vvv",
		'1': "vv",
		'2': "vv>",
		'3': "vv>>",
		'4': "v",
		'5': "v>",
		'6': "v>>",
		'8': ">",
		'9': ">>",
	},
	'8': {
		'A': ">vvv",
		'0': "vvv",
		'1': "vv<",
		'2': "vv",
		'3': "vv>",
		'4': "v<",
		'5': "v",
		'6': "v>",
		'7': "<",
		'9': ">",
	},
	'9': {
		'A': "vvv",
		'0': "<vvv",
		'1': "vv<<",
		'2': "vv<",
		'3': "vv",
		'4': "v<<",
		'5': "v<",
		'6': "v",
		'7': "<<",
		'8': "<",
	},
}

/*
+---+---+---+
| 7 | 8 | 9 |
+---+---+---+
| 4 | 5 | 6 |
+---+---+---+
| 1 | 2 | 3 |
+---+---+---+
    | 0 | A |
    +---+---+

	+---+---+
    | ^ | A |
+---+---+---+
| < | v | > |
+---+---+---+
*/

var directionalMoveMap = map[rune]map[rune]string{
	'A': {
		'>': "v",
		'v': "<v",
		'<': "v<<",
		'^': "<",
	},
	'^': {
		'A': ">",
		'>': "v>",
		'v': "v",
		'<': "v<",
	},
	'>': {
		'A': "^",
		'^': "<^",
		'v': "<",
		'<': "<<",
	},
	'v': {
		'A': "^>",
		'^': "^",
		'>': ">",
		'<': "<",
	},
	'<': {
		'A': ">>^",
		'^': ">^",
		'>': ">>",
		'v': ">",
	},
}

type ExpandKey struct {
	fromKey rune
	code    string
	steps   int
}

var expandMap map[ExpandKey]int = make(map[ExpandKey]int)

func KeyCountAfterExpansion(currentKey rune, code string, times int) int {
	if times == 0 {
		return len(code)
	}
	memoKey := ExpandKey{
		fromKey: currentKey,
		code:    code,
		steps:   times,
	}
	count, ok := expandMap[memoKey]
	if !ok {
		expandedCode := ExpandDirectional(currentKey, code)
		for n, key := range expandedCode {
			if n == 0 {
				count += KeyCountAfterExpansion('A', string(key), times-1)
			} else {
				count += KeyCountAfterExpansion(rune(expandedCode[n-1]), string(key), times-1)
			}
		}
		expandMap[memoKey] = count
	}
	return count
}

func ExpandDirectional(currentKey rune, code string) string {
	var inputBuilder strings.Builder
	for _, c := range code {
		inputBuilder.WriteString(directionalMoveMap[currentKey][c])
		inputBuilder.WriteRune('A')
		currentKey = c
	}
	return inputBuilder.String()
}

func ExpandNumeric(code string) string {
	currentKey := 'A'
	var doorKeyPadBuilder strings.Builder
	for _, c := range code {
		doorKeyPadBuilder.WriteString(numericMoveMap[currentKey][c])
		doorKeyPadBuilder.WriteRune('A')
		currentKey = c
	}

	return doorKeyPadBuilder.String()
}

func solve(code string, numPads int) int {
	iCode, _ := strconv.Atoi(code[:len(code)-1])

	return iCode * KeyCountAfterExpansion('A', ExpandNumeric(code), numPads)
}

func main() {
	fileBytes, _ := os.ReadFile("input.txt")
	part1duration := helpers.MeasureRuntime(func() {
		complexitySum := 0
		for _, s := range strings.Split(string(fileBytes), "\n") {
			complexity := solve(s, 2)
			complexitySum += complexity
		}
		fmt.Println("Part 1:", complexitySum)
	})

	part2duration := helpers.MeasureRuntime(func() {
		p2complexitySum := 0
		for _, s := range strings.Split(string(fileBytes), "\n") {
			p2complexity := solve(s, 25)
			p2complexitySum += p2complexity
		}
		fmt.Println("Part 2:", p2complexitySum)
	})

	fmt.Println("Part 1 took", part1duration.Microseconds(), "μs")
	fmt.Println("Part 2 took", part2duration.Microseconds(), "μs")
}
