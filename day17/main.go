package main

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

func division(numerator int, operand int) int {
	return numerator >> operand
}

func (c *Computer) adv(operand int) {
	c.ra = division(c.ra, operand)
	c.pc += 2
}

func (c *Computer) bxl(operand int) {
	c.rb = c.rb ^ operand
	c.pc += 2
}

func (c *Computer) bst(operand int) {
	c.rb = operand % 8
	c.pc += 2
}

func (c *Computer) jnz(operand int) {
	if c.ra == 0 {
		c.pc += 2
	} else {
		c.pc = operand
	}
}

func (c *Computer) bxc(_ int) {
	c.rb = c.rb ^ c.rc
	c.pc += 2
}

func (c *Computer) out(operand int) {
	c.output = append(c.output, operand%8)
	c.pc += 2
}

func (c *Computer) bdv(operand int) {
	c.rb = division(c.ra, operand)
	c.pc += 2
}

func (c *Computer) cdv(operand int) {
	c.rc = division(c.ra, operand)
	c.pc += 2
}

type Computer struct {
	ra, rb, rc int
	program    []int
	pc         int
	output     []int
}

func (c *Computer) step() (halt bool) {
	opcode, literalOperand := c.program[c.pc], c.program[c.pc+1]
	var operand int
	switch literalOperand {
	case 0:
		operand = 0
	case 1:
		operand = 1
	case 2:
		operand = 2
	case 3:
		operand = 3
	case 4:
		operand = c.ra
	case 5:
		operand = c.rb
	case 6:
		operand = c.rc
	default:
		operand = -1
	}
	switch opcode {
	case 0:
		if operand < 0 {
			return false
		}
		c.adv(operand)
	case 1:
		c.bxl(literalOperand)
	case 2:
		c.bst(operand)
	case 3:
		c.jnz(literalOperand)
	case 4:
		c.bxc(operand)
	case 5:
		c.out(operand)
	case 6:
		if operand < 0 {
			return false
		}
		c.bdv(operand)
	case 7:
		if operand < 0 {
			return false
		}
		c.cdv(operand)
	default:
		panic("invalid opcode")
	}

	return c.pc == len(c.program)
}

func (c *Computer) outputString(ints []int) string {
	out := []string{}
	for _, o := range ints {
		out = append(out, strconv.Itoa(o))
	}
	return strings.Join(out, ",")
}

func (c *Computer) run() {
	for !c.step() {
	}
}

func (c *Computer) String() string {
	var sb strings.Builder
	sb.WriteString("RA: ")
	sb.WriteString(strconv.Itoa(c.ra))
	sb.WriteString("\nRB: ")
	sb.WriteString(strconv.Itoa(c.rb))
	sb.WriteString("\nRC: ")
	sb.WriteString(strconv.Itoa(c.rc))
	sb.WriteString("\nProgram: ")
	sb.WriteString("[")
	sb.WriteString(c.outputString(c.program))
	sb.WriteString("]")
	sb.WriteString("\nPC: ")
	sb.WriteString(strconv.Itoa(c.pc))
	sb.WriteString("\nOut: ")
	sb.WriteString("[")
	sb.WriteString(c.outputString(c.output))
	sb.WriteString("]")

	return sb.String()
}

func (c *Computer) calculateQuine() uint64 {
	candidates := map[uint64]bool{
		0: true,
	}
	reverseProgram := make([]int, len(c.program))
	copy(reverseProgram, c.program)
	slices.Reverse(reverseProgram)
	for _, outputTarget := range reverseProgram {
		toRemove := map[uint64]bool{}
		toAdd := map[uint64]bool{}
		for cand := range candidates {
			for delta := range 8 {
				raCandidate := cand<<3 + uint64(delta)
				m3 := raCandidate % 8
				s5 := m3 ^ 5
				rs5 := (raCandidate >> s5) % 8
				if m3^rs5^3 == uint64(outputTarget) {
					toAdd[raCandidate] = true
				}
			}
			toRemove[cand] = true
		}
		for candidate := range toRemove {
			delete(candidates, candidate)
		}
		for candidate := range toAdd {
			candidates[candidate] = true
		}
	}
	var minimumValue uint64 = math.MaxUint64
	for raCandidate := range candidates {
		if raCandidate < minimumValue {
			minimumValue = raCandidate
		}
	}
	return minimumValue
}

func newComputer(input string) *Computer {
	inputSections := strings.Split(input, "\n\n")
	registerRE := regexp.MustCompile(`^Register ([ABC]): (\d+)$`)
	programRE := regexp.MustCompile(`^Program: ([\d,]+)$`)
	var ra, rb, rc int
	for _, l := range strings.Split(inputSections[0], "\n") {
		match := registerRE.FindStringSubmatch(l)
		v, _ := strconv.Atoi(match[2])
		switch match[1] {
		case "A":
			ra = v
		case "B":
			rb = v
		case "C":
			rc = v
		default:
			panic("BANG!")
		}
	}
	program := []int{}
	programMatch := programRE.FindStringSubmatch(inputSections[1])
	for _, opcode := range strings.Split(programMatch[1], ",") {
		o, _ := strconv.Atoi(opcode)
		program = append(program, o)
	}
	pc := 0
	output := []int{}
	return &Computer{
		ra, rb, rc, program, pc, output,
	}
}

func main() {
	fileBytes, _ := os.ReadFile("input.txt")
	referenceComputer := newComputer(string(fileBytes))
	referenceComputer.run()
	fmt.Println("Part 1:", referenceComputer.outputString(referenceComputer.output))
	fmt.Println("Part 2:", referenceComputer.calculateQuine())
}
