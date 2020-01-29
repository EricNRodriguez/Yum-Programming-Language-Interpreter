package parser

import "github.com/EricNRodriguez/yum/ast"

type Parser interface {
	Parse() (*ast.Program, []error)
}
