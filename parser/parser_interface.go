package parser

import "Yum-Programming-Language-Interpreter/ast"

type Parser interface {
	Parse() *ast.Program
}