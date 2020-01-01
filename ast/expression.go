package ast

import (
	"Yum-Programming-Language-Interpreter/token"
	"bytes"
	"fmt"
	"strconv"
)

type prefix struct {
	token.Token
	expression Expression
}

func NewPrefixExpression(t token.Token, e Expression) Expression {
	return &prefix{
		Token: t,
		expression:     e,
	}
}

func (p *prefix) String() string {
	return fmt.Sprintf("(%v%v)", p.Literal(), p.expression.String())
}

func (p *prefix) expressionFunction() {}

type infix struct {
	token.Token
	leftExpression  Expression
	rightExpression Expression
}

func NewInfixExpression(t token.Token, le, re Expression) Expression {
	return &infix{
		Token:  t,
		leftExpression:  le,
		rightExpression: re,
	}
}

func (i *infix) String() string {
	return fmt.Sprintf("(%v %v %v)", i.leftExpression.String(), i.Token.Literal(), i.rightExpression.String())
}

func (i *infix) expressionFunction() {}

type integerExpression struct {
	token.Metadata
	value int
}

func NewIntegerExpression(t token.Token, i int) Expression {
	return &integerExpression{
		Metadata: t.Data(),
		value:             i,
	}
}

func (ie *integerExpression) String() string {
	return strconv.Itoa(ie.value)
}

func (te *integerExpression) expressionFunction() {}

type booleanExpression struct {
	token.Metadata
	Value bool
}

func NewBooleanExpression(t token.Token, v bool) Expression {
	return &booleanExpression{
		Metadata: t.Data(),
		Value:             v,
	}
}

func (be *booleanExpression) String() string {
	if be.Value {
		return "true"
	}
	return "false"
}

func (be *booleanExpression) expressionFunction() {}

type functionCallExpression struct {
	token.Metadata
	FunctionName string
	Parameters   []Expression
}

func NewFunctionCallExpression(md token.Metadata, fName string, params ...Expression) Expression {
	return &functionCallExpression{
		Metadata: md,
		FunctionName:      fName,
		Parameters:        params,
	}
}

func (fc *functionCallExpression) String() string {
	pBuff := bytes.Buffer{}
	if fc.Parameters != nil {
		for i, param := range fc.Parameters {
			pBuff.WriteString(param.String())
			if i != len(fc.Parameters) -1 {
				pBuff.WriteString(", ")
			}
		}
	}

	return fmt.Sprintf("%v(%v)", fc.FunctionName, pBuff.String())
}

func (fc *functionCallExpression) expressionFunction() {}

type identifierExpression struct {
	*Identifier
}

func NewIdentifierExpression(t token.Token) Expression {
	return &identifierExpression{
		Identifier: NewIdentifier(t),
	}
}

func (i *identifierExpression) expressionFunction() {}
