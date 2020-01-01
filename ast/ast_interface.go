package ast

import "Yum-Programming-Language-Interpreter/token"

type Node interface {
	String() string         // string representation of expression
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
