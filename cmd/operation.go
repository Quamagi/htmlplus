package cmd

import (
	"fmt"
	"strconv"
)

// OperationCommand handles arithmetic and logical operations
func OperationCommand(attributes string, operands []string) {
	attrMap := ParseAttributes(attributes)
	resultVar := attrMap["result"]

	if len(operands) != 2 {
		Variables[resultVar] = "Error: exactly two operands are required"
		return
	}

	op1Str := operands[0]
	op2Str := operands[1]

	op1, err1 := strconv.ParseFloat(op1Str, 64)
	if err1 != nil {
		op1Var, exists := Variables[op1Str]
		if !exists {
			Variables[resultVar] = "Error: operand 1 not found"
			return
		}
		op1, err1 = strconv.ParseFloat(fmt.Sprintf("%v", op1Var), 64)
		if err1 != nil {
			Variables[resultVar] = "Error: operand 1 is not a number"
			return
		}
	}

	op2, err2 := strconv.ParseFloat(op2Str, 64)
	if err2 != nil {
		op2Var, exists := Variables[op2Str]
		if !exists {
			Variables[resultVar] = "Error: operand 2 not found"
			return
		}
		op2, err2 = strconv.ParseFloat(fmt.Sprintf("%v", op2Var), 64)
		if err2 != nil {
			Variables[resultVar] = "Error: operand 2 is not a number"
			return
		}
	}

	var result float64

	switch attrMap["type"] {
	case "addition":
		result = op1 + op2
	case "subtraction":
		result = op1 - op2
	case "multiplication":
		result = op1 * op2
	case "division":
		if op2 == 0 {
			Variables[resultVar] = "Error: division by zero"
			return
		}
		result = op1 / op2
	case "greaterThan":
		Variables[resultVar] = op1 > op2
		return
	case "equalTo":
		Variables[resultVar] = op1 == op2
		return
	case "and":
		Variables[resultVar] = op1 != 0 && op2 != 0
		return
	case "or":
		Variables[resultVar] = op1 != 0 || op2 != 0
		return
	default:
		Variables[resultVar] = "Error: unknown operation type"
		return
	}

	Variables[resultVar] = result
}
