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

	// 检查表达式是否包含非法字符 (允许 s, q, r, t, ^)
	for _, c := range expr {
		if !strings.ContainsRune("0123456789.+-*/%()sqrt^ ", c) {
			return 0, fmt.Errorf("invalid character in expression: %c", c)
		}
	}

	// 处理特殊一元运算
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

	// (expression)^2 已经在 evaluateExpression 中通过 types.Eval 的回退逻辑处理
	// 但为了更明确，可以考虑在这里也添加一个替换，但这可能导致与 types.Eval 冲突
	// 例如，如果 types.Eval 本身就能处理 (val)^2，这里的替换就是多余的

	return expr
}
