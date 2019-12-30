package ast

import (
	"Yum-Programming-Language-Interpreter/token"
	"fmt"
	"strconv"
)

type Prefix struct {
	token.TokenInterface
	expression ExpressionInterface
}

func NewPrefixExpression(t token.TokenInterface, e ExpressionInterface) *Prefix {
	return &Prefix{
		TokenInterface: t,
		expression: e,
	}
}

func (p *Prefix) String() string {
	return fmt.Sprintf("(%v%v)", p.Literal(), p.expression.String())
}

func (p *Prefix) expressionFunction() {}

type Infix struct {
	token.TokenInterface
	leftExpression ExpressionInterface
	rightExpression ExpressionInterface

}

func NewInfixExpression(t token.TokenInterface, le, re ExpressionInterface) ExpressionInterface {
	return &Infix{
		TokenInterface: t,
		leftExpression: le,
		rightExpression: re,
	}
}

func (i *Infix) String() string {
	return fmt.Sprintf("(%v %v %v)", i.leftExpression.String(), i.TokenInterface.Literal(), i.rightExpression.String())
}

func (i *Infix) expressionFunction() {}

type TokenExpression struct {
	token.TokenInterface
}

func NewTokenExpression(t token.TokenInterface) ExpressionInterface {
	return &TokenExpression{
		TokenInterface: t,
	}
}

func (te *TokenExpression) String() string {
	return te.TokenInterface.Literal()
}

func (te *TokenExpression) expressionFunction() {}


type IntegerExpression struct {
	token.MetadataInterface
	value int
}

func NewIntegerExpression(t token.TokenInterface, i int) *IntegerExpression {
	return &IntegerExpression{
		MetadataInterface: t.Metadata(),
		value:          i,
	}
}

func (ie *IntegerExpression) String() string {
	return strconv.Itoa(ie.value)
}

func (te *IntegerExpression) expressionFunction() {}


type BooleanExpression struct {
	token.MetadataInterface
	Value bool
}

func NewBooleanExpression(t token.TokenInterface, v bool) *BooleanExpression {
	return &BooleanExpression{
		MetadataInterface:t.Metadata(),
		Value:v,
	}
}

func (be *BooleanExpression) String() string {
	if be.Value {
		return "true"
	}
	return "false"
}

func (be *BooleanExpression) expressionFunction() {}



