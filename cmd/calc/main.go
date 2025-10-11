package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/dxasu/pure/stdin"
	_ "github.com/dxasu/pure/version"
)

// 运算符优先级
var precedence = map[string]int{
	"+": 1, "-": 1,
	"*": 2, "/": 2,
	"^": 3,
}

// 将中缀表达式转为后缀表达式（逆波兰式）
func infixToPostfix(expr string) []string {
	var output []string
	var stack []string

	tokens := parseTokens(expr)
	for _, token := range tokens {
		if token[len(token)-1] == '%' {
			t, err := strconv.ParseFloat(token[:len(token)-1], 64)
			if err != nil {
				panic(fmt.Sprintf("invalid percentage token: %s", token))
			}
			output = append(output, strconv.FormatFloat(t/100, 'f', -1, 64))
		} else if isNumber(token) {
			output = append(output, token)
		} else if token == "(" {
			stack = append(stack, token)
		} else if token == ")" {
			// 弹出栈顶元素直到遇到 "("
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = stack[:len(stack)-1] // 弹出 "("
		} else {
			// 处理运算符优先级
			for len(stack) > 0 && stack[len(stack)-1] != "(" &&
				precedence[token] <= precedence[stack[len(stack)-1]] {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		}
	}

	// 弹出剩余运算符
	for len(stack) > 0 {
		output = append(output, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return output
}

// 解析表达式为Token列表
func parseTokens(expr string) []string {
	var tokens []string
	var currentToken strings.Builder

	for _, ch := range expr {
		if unicode.IsSpace(ch) {
			continue
		}

		if isOperatorOrBracket(ch) {
			if currentToken.Len() > 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}
			tokens = append(tokens, string(ch))
		} else {
			currentToken.WriteRune(ch)
		}
	}

	if currentToken.Len() > 0 {
		tokens = append(tokens, currentToken.String())
	}

	return tokens
}

// 判断是否为运算符或括号
func isOperatorOrBracket(ch rune) bool {
	return ch == '+' || ch == '-' || ch == '*' || ch == 'x' || ch == 'X' || ch == '/' ||
		ch == '^' || ch == '(' || ch == ')'
}

// 判断是否为数字（含小数和负数）
func isNumber(token string) bool {
	_, err := strconv.ParseFloat(token, 64)
	return err == nil
}

// 计算后缀表达式
func evaluatePostfix(tokens []string) float64 {
	var stack []float64

	for _, token := range tokens {
		if isNumber(token) {
			num, _ := strconv.ParseFloat(token, 64)
			stack = append(stack, num)
		} else {
			// 弹出右操作数和左操作数
			right := stack[len(stack)-1]
			left := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			var result float64
			switch token {
			case "+":
				result = left + right
			case "-":
				result = left - right
			case "*", "x", "X":
				result = left * right
			case "/":
				result = left / right
			case "^":
				result = math.Pow(left, right)
			}
			stack = append(stack, result)
		}
	}

	return stack[0]
}

func main() {
	if len(os.Args) == 2 && os.Args[1] == "-h" {
		cmd := os.Args[0]
		fmt.Println(cmd, ` xxx # calculate expression, xxx from params,clipboard,stdin`)
		return
	}
	expr := stdin.GetInput()
	postfix := infixToPostfix(expr)
	result := evaluatePostfix(postfix)
	fmt.Println(strconv.FormatFloat(result, 'g', 12, 64))
}
