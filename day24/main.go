package main

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"
)

type Expression struct {
	operand1, operand2 string
	operation          string
}

func (e Expression) String() string {
	switch e.operation {
	case "NOT":
		return fmt.Sprintf("NOT %s", e.operand1)
	case "EQ":
		return e.operand1
	default:
		return fmt.Sprintf("%s %s %s", e.operand1, e.operation, e.operand2)
	}
}

func isUnary(operation string) bool {
	return operation == "EQ" || operation == "NOT"
}

func EvaluateExpression(e Expression, inputMap map[string]bool) (bool, *Expression) {
	op1value, op1found := inputMap[e.operand1]
	op2value, op2found := inputMap[e.operand2]

	if op1found && (isUnary(e.operation) || op2found) {
		switch e.operation {
		case "AND":
			return op1value && op2value, nil
		case "OR":
			return op1value || op2value, nil
		case "XOR":
			return op1value != op2value, nil
		case "EQ":
			return op1value, nil
		case "NOT":
			return !op1value, nil
		default:
			panic("Couldn't find operation")
		}
	} else if op1found || op2found {
		var foundValue bool
		var missingVariable string

		if op1found {
			missingVariable = e.operand2
			foundValue = op1value
		} else {
			missingVariable = e.operand1
			foundValue = op2value
		}

		switch e.operation {
		case "AND":
			if foundValue {
				return false, &Expression{operand1: missingVariable, operation: "EQ"}
			} else {
				return false, nil
			}
		case "OR":
			if !foundValue {
				return false, &Expression{operand1: missingVariable, operation: "EQ"}
			} else {
				return true, nil
			}
		case "XOR":
			if !foundValue {
				return false, &Expression{operand1: missingVariable, operation: "EQ"}
			} else {
				return false, &Expression{operand1: missingVariable, operation: "NOT"}
			}
		default:
			panic("Couldn't find operation")
		}
	} else {
		return false, &e
	}
}

func swap(k1, k2 string, inputMap map[string]Expression) []string {
	k2val := inputMap[k2]
	inputMap[k2] = inputMap[k1]
	inputMap[k1] = k2val
	return []string{k1, k2}
}

func main() {
	fileBytes, _ := os.ReadFile("input.txt")
	input := strings.Split(string(fileBytes), "\n\n")
	inputs, formulas := strings.Split(input[0], "\n"), strings.Split(input[1], "\n")
	inputMap := make(map[string]bool, len(inputs))
	inputRe := regexp.MustCompile(`^([a-z0-9]+): ([01])`)
	for _, i := range inputs {
		inputMatch := inputRe.FindStringSubmatch(i)
		inputMap[inputMatch[1]] = inputMatch[2] == "1"
	}

	formulaMap := make(map[string]Expression, len(formulas))
	formulaRe := regexp.MustCompile(`^([a-z0-9]+) (AND|OR|XOR) ([a-z0-9]+) -> ([a-z0-9]+)$`)
	for _, f := range formulas {
		formulaMatch := formulaRe.FindStringSubmatch(f)
		formulaMap[formulaMatch[4]] = Expression{
			operand1:  formulaMatch[1],
			operation: formulaMatch[2],
			operand2:  formulaMatch[3],
		}
	}

	p1, _ := evaluate(formulaMap, inputMap)
	fmt.Println("Part 1:", p1)
	idx := 0
	swaps := []string{}
	swaps = append(swaps, swap("z20", "hhh", formulaMap)...)
	swaps = append(swaps, swap("ggk", "rhv", formulaMap)...)
	swaps = append(swaps, swap("htp", "z15", formulaMap)...)
	swaps = append(swaps, swap("z05", "dkr", formulaMap)...)
	for {
		key := fmt.Sprintf("z%02d", idx)
		zFormula, ok := formulaMap[key]
		if !ok {
			break
		}
		op1formula, fop1 := formulaMap[zFormula.operand1]
		op2formula, fop2 := formulaMap[zFormula.operand2]

		op1str := zFormula.operand1
		op2str := zFormula.operand2
		xor := 0
		if fop1 {
			op1str = op1formula.String()
			if op1formula.operation == "XOR" {
				xor = 1
			}
		}
		if fop2 {
			op2str = op2formula.String()
			if op2formula.operation == "XOR" {
				xor = 2
			}
		}
		if xor == 2 {
			fmt.Println(key, "=", zFormula, "\t", zFormula.operand2, "=", op2str, "\t", zFormula.operand1, "=", op1str)
		} else {
			fmt.Println(key, "=", zFormula, "\t", zFormula.operand1, "=", op1str, "\t", zFormula.operand2, "=", op2str)
		}
		idx++
	}

	idx--

	fmt.Println()
	fmt.Println()
	fmt.Println()

	for z := range idx {
		v := uint64(1 << z)
		hasError := false
		input := getCleanInputs(idx)
		input[fmt.Sprintf("x%02d", z)] = true
		x1, resInput := evaluate(formulaMap, input)
		if x1 != v {
			fmt.Println("ERROR x1, correct:", v)
			fmt.Println("x", z, "=", x1)
			hasError = true
			oneKeys := printOnes(resInput)
			for _, k := range oneKeys {
				fmt.Println(k, "=>", formulaMap[k])
			}
		}

		input = getCleanInputs(idx)
		input[fmt.Sprintf("y%02d", z)] = true
		y1, resInput := evaluate(formulaMap, input)
		if y1 != v {
			fmt.Println("ERROR y1, correct:", v)
			fmt.Println("y", z, "=", y1)
			hasError = true
			oneKeys := printOnes(resInput)
			for _, k := range oneKeys {
				fmt.Println(k, "=>", formulaMap[k])
			}
		}

		input = getCleanInputs(idx)
		input[fmt.Sprintf("x%02d", z)] = true
		input[fmt.Sprintf("y%02d", z)] = true

		xy1, resInput := evaluate(formulaMap, input)
		if xy1 != v*2 {
			fmt.Println("ERROR xy1, correct:", v*2)
			fmt.Println("xy", z, "=", xy1)
			hasError = true
			oneKeys := printOnes(resInput)
			for _, k := range oneKeys {
				fmt.Println(k, "=>", formulaMap[k])
			}
		}
		if hasError {
			fmt.Println()
		}
	}
	slices.Sort(swaps)
	fmt.Println(strings.Join(swaps, ","))
}

func getCleanInputs(idx int) map[string]bool {
	input := make(map[string]bool, idx*2)
	for zz := range idx {
		input[fmt.Sprintf("x%02d", zz)] = false
		input[fmt.Sprintf("y%02d", zz)] = false
	}
	return input
}

func printOnes(inputMap map[string]bool) []string {
	keys := []string{}
	for k, v := range inputMap {
		if v {
			keys = append(keys, k)
		}
	}
	return keys
}

func evaluate(formulaMap map[string]Expression, inputMap map[string]bool) (uint64, map[string]bool) {
	formulas := make(map[string]Expression, len(formulaMap))
	for k, v := range formulaMap {
		formulas[k] = v
	}
	inputs := make(map[string]bool, len(inputMap))
	for k, v := range inputMap {
		inputs[k] = v
	}
	for len(formulas) > 0 {

		for variable, expression := range formulas {
			result, retExpression := EvaluateExpression(expression, inputs)
			if retExpression == nil {
				inputs[variable] = result
			} else {
				formulaMap[variable] = expression
			}
		}

		for variable := range inputs {
			delete(formulas, variable)
		}
	}
	return getInteger("z", inputs), inputs
}

func getInteger(prefix string, inputMap map[string]bool) (result uint64) {
	idx := 0
	result = 0
	for {
		key := fmt.Sprintf("%s%02d", prefix, idx)
		if res, ok := inputMap[key]; ok {
			if res {
				result += 1 << idx
			}
			idx++
		} else {
			break
		}
	}
	return result
}
