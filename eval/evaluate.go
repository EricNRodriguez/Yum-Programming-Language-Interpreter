package eval

import (
	"Yum-Programming-Language-Interpreter/ast"
	"Yum-Programming-Language-Interpreter/object"
)

var sT = NewSymbolTable()

func Evaluate(node ast.Node) object.Object {
	if method, ok := evalMethodRouter[node.Type()]; ok {
		return method(node)
	}
	return nil
}
