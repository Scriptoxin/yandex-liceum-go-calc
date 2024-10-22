package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func Calc(expression string) (float64, error) {
	expression = strings.ReplaceAll(expression, " ", "")
	result, _, err := eval(expression, 0)
	return result, err
}

func eval(expression string, pos int) (float64, int, error) {
	var numStack []float64
	var opStack []byte

	for pos < len(expression) {
		char := expression[pos]

		if unicode.IsDigit(rune(char)) || char == '.' {
			value, newPos, err := parseNumber(expression, pos)
			if err != nil {
				return 0, pos, err
			}
			numStack = append(numStack, value)
			pos = newPos
		} else if char == '(' {
			value, newPos, err := eval(expression, pos+1)
			if err != nil {
				return 0, pos, err
			}
			numStack = append(numStack, value)
			pos = newPos
		} else if char == ')' {
			break
		} else if isOperator(char) {
			for len(opStack) > 0 && precedence(opStack[len(opStack)-1]) >= precedence(char) {
				result, err := applyOperation(&numStack, &opStack)
				if err != nil {
					return 0, pos, err
				}
				numStack = append(numStack, result)
			}
			opStack = append(opStack, char)
			pos++
		} else {
			return 0, pos, errors.New("invalid character in expression")
		}
	}

	for len(opStack) > 0 {
		result, err := applyOperation(&numStack, &opStack)
		if err != nil {
			return 0, pos, err
		}
		numStack = append(numStack, result)
	}

	if len(numStack) != 1 {
		return 0, pos, errors.New("invalid expression")
	}

	return numStack[0], pos + 1, nil
}

func parseNumber(expression string, pos int) (float64, int, error) {
	startPos := pos
	for pos < len(expression) && (unicode.IsDigit(rune(expression[pos])) || expression[pos] == '.') {
		pos++
	}
	value, err := strconv.ParseFloat(expression[startPos:pos], 64)
	if err != nil {
		return 0, pos, errors.New("invalid number")
	}
	return value, pos, nil
}

func applyOperation(numStack *[]float64, opStack *[]byte) (float64, error) {
	if len(*numStack) < 2 {
		return 0, errors.New("not enough operands")
	}

	b := (*numStack)[len(*numStack)-1]
	a := (*numStack)[len(*numStack)-2]
	*numStack = (*numStack)[:len(*numStack)-2]

	op := (*opStack)[len(*opStack)-1]
	*opStack = (*opStack)[:len(*opStack)-1]

	switch op {
	case '+':
		return a + b, nil
	case '-':
		return a - b, nil
	case '*':
		return a * b, nil
	case '/':
		if b == 0 {
			return 0, errors.New("division by zero")
		}
		return a / b, nil
	default:
		return 0, errors.New("unknown operator")
	}
}

func isOperator(char byte) bool {
	return char == '+' || char == '-' || char == '*' || char == '/'
}

func precedence(op byte) int {
	switch op {
	case '+', '-':
		return 1
	case '*', '/':
		return 2
	}
	return 0
}

func main() {
	expr := "1.5 * 1.5"
	result, err := Calc(expr)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result)
	}
}
