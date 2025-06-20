package main

import (
	"encoding/json"
	"fmt"
	"go/token"
	"go/types"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type CalculateRequest struct {
	Expression string `json:"expression"`
}

type CalculateResponse struct {
	Result float64 `json:"result,omitempty"`
	Error  string  `json:"error,omitempty"`
}

func main() {
	http.HandleFunc("/calculate", handleCalculate)
	port := 8000
	log.Printf("Starting Go backend server on port %d...", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handleCalculate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(CalculateResponse{Error: "Only POST method is allowed"})
		return
	}

	var req CalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CalculateResponse{Error: "Invalid request format"})
		return
	}

	result, err := evaluateExpression(req.Expression)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CalculateResponse{Error: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(CalculateResponse{Result: result})
}

func evaluateExpression(expr string) (float64, error) {
	// 移除所有空格
	expr = strings.ReplaceAll(expr, " ", "")

	// 检查表达式是否为空
	if expr == "" {
		return 0, fmt.Errorf("empty expression")
	}

	// 检查表达式是否包含非法字符 (允许 s, q, r, t, ^, %)
	for _, c := range expr {
		if !strings.ContainsRune("0123456789.+-*/()sqrt^% ", c) {
			return 0, fmt.Errorf("invalid character in expression: %c", c)
		}
	}

	// 优先处理百分比，将其作为一元运算符
	expr = preprocessPercent(expr)

	// 预处理表达式，确保所有数字都是浮点数，以强制执行浮点除法
	expr = ensureFloatLiterals(expr)

	// 处理其他一元运算
	expr = preprocessUnaryOperations(expr)

	// 使用 go/types 包来评估表达式
	fset := token.NewFileSet()
	tv, err := types.Eval(fset, nil, token.NoPos, expr)
	if err != nil {
		// 尝试移除括号后再次评估，以处理 (number)^2 的情况
		if strings.HasSuffix(expr, ")^2") && strings.HasPrefix(expr, "(") {
			innerExpr := expr[1 : len(expr)-3]
			tv, err = types.Eval(fset, nil, token.NoPos, innerExpr)
			if err == nil {
				val, _ := strconv.ParseFloat(tv.Value.String(), 64)
				return math.Pow(val, 2), nil
			}
		}
		return 0, fmt.Errorf("invalid expression: %v", err)
	}

	// 将结果转换为 float64
	result, err := strconv.ParseFloat(tv.Value.String(), 64)
	if err != nil {
		return 0, fmt.Errorf("error converting result: %v", err)
	}

	return result, nil
}

// preprocessPercent 预处理百分比表达式，将 X% 转换为 (X/100)
func preprocessPercent(expr string) string {
	// 此正则表达式查找后跟“%”的数字（整数或浮点数）
	re := regexp.MustCompile(`(\d+\.?\d*)%`)
	return re.ReplaceAllString(expr, "($1/100)")
}

// preprocessUnaryOperations 预处理一元运算表达式
func preprocessUnaryOperations(expr string) string {
	// 处理 sqrt(expression)
	sqrtRegex := regexp.MustCompile(`sqrt\(([^)]+)\)`)
	expr = sqrtRegex.ReplaceAllStringFunc(expr, func(match string) string {
		innerExpr := sqrtRegex.FindStringSubmatch(match)[1]
		val, err := evaluateExpression(innerExpr) // 递归处理内部表达式
		if err != nil {
			return match // 如果内部表达式无效，则不替换
		}
		return strconv.FormatFloat(math.Sqrt(val), 'f', -1, 64)
	})

	// 处理 1/(expression)
	invRegex := regexp.MustCompile(`1/\(([^)]+)\)`)
	expr = invRegex.ReplaceAllStringFunc(expr, func(match string) string {
		innerExpr := invRegex.FindStringSubmatch(match)[1]
		val, err := evaluateExpression(innerExpr)
		if err != nil || val == 0 {
			return match // 如果内部表达式无效或为0，则不替换
		}
		return strconv.FormatFloat(1/val, 'f', -1, 64)
	})

	return expr
}

// ensureFloatLiterals 将表达式中的所有整数字面量转换为浮点数字面量
func ensureFloatLiterals(expr string) string {
	// 这个正则表达式查找数字，并使用 ReplaceAllStringFunc 来处理它们
	re := regexp.MustCompile(`-?\d+(\.\d+)?`)

	return re.ReplaceAllStringFunc(expr, func(match string) string {
		// 如果匹配到的数字已经包含小数点，则不进行任何操作
		if strings.Contains(match, ".") {
			return match
		}
		// 否则，在数字后面添加 ".0"，将其转换为浮点数
		return match + ".0"
	})
}
