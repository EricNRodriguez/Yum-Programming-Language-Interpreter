package eval

import (
	"Yum-Programming-Language-Interpreter/ast"
)

type StackTrace interface {
	Push(*ast.FunctionCallExpression)
	Pop() (*ast.FunctionCallExpression, bool)
}

type stackTrace []*ast.FunctionCallExpression

func NewStackTrace() *stackTrace {
	return &stackTrace{}
}

func (st *stackTrace) Push(fc *ast.FunctionCallExpression) {
	*st = append(*st, fc)
	return
}

func (st *stackTrace) Pop() (fc *ast.FunctionCallExpression, ok bool) {
	if len(*st) > 0 {
		fc = (*st)[len(*st)-1]
		ok = true
		// pop
		*st = (*st)[:len(*st)-1]
	}
	return
}
