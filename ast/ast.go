package ast

import "Yum-Programming-Language-Interpreter/token"

type NodeInterface interface {
	Parent() NodeInterface     // parent node in ast, nil if root
	Children() []NodeInterface // children in ast, nil if leaf node (token)
	String() string            // string representation of expression
	token.TokenInterface       // literal and metadata
}

type StatementInterface interface {
	NodeInterface
	statementFunction()
}

type ExpressionInterface interface {
	NodeInterface
	expressionFunction()
}
