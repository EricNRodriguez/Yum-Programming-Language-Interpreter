package ast

import (
	"github.com/EricNRodriguez/yum/token"
)

type Node interface {
	String() string // string representation of expression
	Type() NodeType
	token.Metadata // literal and metadata
}

type Statement interface {
	Node
	statementFunction()
}

type Expression interface {
	Node
	expressionFunction()
}
