package ast

import (
	"Yum-Programming-Language-Interpreter/token"
	"fmt"
	"strconv"
	"strings"
)

type PrefixExpression struct {
	token.Token
	Expression Expression
}

func NewPrefixExpression(t token.Token, e Expression) *PrefixExpression {
	return &PrefixExpression{
		Token:      t,
		Expression: e,
	}
}

func (p *PrefixExpression) String() string {
	return fmt.Sprintf("(%v%v)", p.Literal(), p.Expression.String())
}

func (p *PrefixExpression) Type() NodeType {
	return PREFIX_EXPRESSION
}

func (p *PrefixExpression) expressionFunction() {}

type InfixExpression struct {
	token.Token
	LeftExpression  Expression
	RightExpression Expression
}

func NewInfixExpression(t token.Token, le, re Expression) Expression {
	return &InfixExpression{
		Token:           t,
		LeftExpression:  le,
		RightExpression: re,
	}
}

func (i *InfixExpression) String() string {
	return fmt.Sprintf("(%v %v %v)", i.LeftExpression.String(), i.Token.Literal(), i.RightExpression.String())
}

func (i *InfixExpression) Type() NodeType {
	return INFIX_EXPRESSION
}

func (i *InfixExpression) expressionFunction() {}

type IntegerExpression struct {
	token.Metadata
	Value int
}

func NewIntegerExpression(t token.Token, i int) *IntegerExpression {
	return &IntegerExpression{
		Metadata: t.Data(),
		Value:    i,
	}
}

func (ie *IntegerExpression) String() string {
	return strconv.Itoa(ie.Value)
}

func (ie *IntegerExpression) Type() NodeType {
	return INTEGER_EXPRESSION
}

func (ie *IntegerExpression) expressionFunction() {}

type BooleanExpression struct {
	token.Metadata
	Value bool
}

func NewBooleanExpression(t token.Token, v bool) *BooleanExpression {
	return &BooleanExpression{
		Metadata: t.Data(),
		Value:    v,
	}
}

func (be *BooleanExpression) String() string {
	if be.Value {
		return "true"
	}
	return "false"
}

func (be *BooleanExpression) Type() NodeType {
	return BOOLEAN_EXPRESSION
}

func (be *BooleanExpression) expressionFunction() {}

type FunctionCallExpression struct {
	token.Metadata
	FunctionName string
	Parameters   []Expression
}

func NewFunctionCallExpression(md token.Metadata, fName string, params ...Expression) Expression {
	return &FunctionCallExpression{
		Metadata:     md,
		FunctionName: fName,
		Parameters:   params,
	}
}

func (fc *FunctionCallExpression) String() string {
	return fmt.Sprintf("%v(%v)", fc.FunctionName, expressionArrayToString(fc.Parameters))
}

func (fc *FunctionCallExpression) Type() NodeType {
	return FUNC_CALL_EXPRESSION
}

func (fc *FunctionCallExpression) expressionFunction() {}

type IdentifierExpression struct {
	Node
}

func NewIdentifierExpression(t token.Token) Expression {
	return &IdentifierExpression{
		Node: NewIdentifier(t),
	}
}

func (i *IdentifierExpression) Type() NodeType {
	return IDENTIFIER_EXPRESSION
}

func (i *IdentifierExpression) expressionFunction() {}


func expressionArrayToString(staArr []Expression) string {
	var strArr = make([]string, len(staArr))
	for i, sta := range staArr {
		strArr[i] = sta.String()
	}
	return strings.Join(strArr, ", ")
}