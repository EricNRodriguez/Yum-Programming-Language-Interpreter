package ast

import (
	"Yum-Programming-Language-Interpreter/token"
	"bytes"
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
	Value int64
}

func NewIntegerExpression(t token.Token, i int64) *IntegerExpression {
	return &IntegerExpression{
		Metadata: t.Data(),
		Value:    i,
	}
}

func (ie *IntegerExpression) String() string {
	return strconv.FormatInt(ie.Value, 10)
}

func (ie *IntegerExpression) Type() NodeType {
	return INTEGER_EXPRESSION
}

func (ie *IntegerExpression) expressionFunction() {}

type FloatingPointExpression struct {
	token.Metadata
	Value float64
}

func NewFloatingPointExpression(t token.Token, i float64) *FloatingPointExpression {
	return &FloatingPointExpression{
		Metadata: t.Data(),
		Value:    i,
	}
}

func (fpe *FloatingPointExpression) String() string {
	return fmt.Sprintf("%f", fpe.Value)
}

func (fpe *FloatingPointExpression) Type() NodeType {
	return FLOATING_POINT_EXPRESSION
}

func (fpe *FloatingPointExpression) expressionFunction() {}

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
	*Identifier
}

func NewIdentifierExpression(t token.Token) Expression {
	return &IdentifierExpression{
		Identifier: NewIdentifier(t),
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

type ArrayExpression struct {
	token.Metadata
	Data   []Expression
	Length int
}

func NewArray(md token.Metadata, data []Expression) *ArrayExpression {
	return &ArrayExpression{
		Metadata: md,
		Data:     data,
		Length:   len(data),
	}
}

func (a *ArrayExpression) String() string {
	buff := bytes.Buffer{}
	buff.WriteString("[")
	for i, e := range a.Data {
		buff.WriteString(e.String())
		if i != len(a.Data)-1 {
			buff.WriteString(",")
		}
	}
	buff.WriteString("]")
	return buff.String()
}

func (a *ArrayExpression) Type() NodeType {
	return ARRAY
}

func (a *ArrayExpression) expressionFunction() {}

type ArrayIndexExpression struct {
	token.Metadata
	ArrayName string
	IndexExpr Expression
}

func NewArrayIndexExpression(md token.Metadata, n string, e Expression) *ArrayIndexExpression {
	return &ArrayIndexExpression{
		Metadata: md,
		ArrayName: n,
		IndexExpr: e,
	}
}

func (a *ArrayIndexExpression) String() string {
	return fmt.Sprintf("%v[%v];", a.ArrayName, a.IndexExpr)
}

func (a *ArrayIndexExpression) Type() NodeType {
	return ARRAY_INDEX_EXPRESSION
}

func (a *ArrayIndexExpression) expressionFunction() {}
