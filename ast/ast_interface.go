package ast

import "Yum-Programming-Language-Interpreter/token"

type NodeInterface interface {
	String() string            // string representation of expression
	token.MetadataInterface       // literal and metadata
}

type StatementInterface interface {
	NodeInterface
	statementFunction()
}

type ExpressionInterface interface {
	NodeInterface
	expressionFunction()
}
