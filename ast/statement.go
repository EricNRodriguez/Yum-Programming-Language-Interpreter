package ast

import (
	"Yum-Programming-Language-Interpreter/token"
	"strings"
	"fmt"
)

type varStatement struct {
	token.Token
	identifier token.Token
	expression Expression
}

func NewVarStatement(t, i token.Token, e Expression) Statement {
	return &varStatement{
		Token: t,
		identifier:     i,
		expression:     e,
	}
}

func (v *varStatement) String() string {
	return fmt.Sprintf("var %v = %v;", v.identifier.Literal(), v.expression.String())
}

func (v *varStatement) statementFunction() {}

type returnStatement struct {
	token.Token
	expression Expression
}

func NewReturnStatment(t token.Token, e Expression) Statement {
	return &returnStatement{
		Token: t,
		expression:     e,
	}
}

func (r *returnStatement) String() string {
	return fmt.Sprintf("return %v;", r.expression.String())
}

func (r *returnStatement) statementFunction() {}

type expressionStatement struct {
	Expression
}

func NewExpressionStatment(e Expression) Statement {
	return &expressionStatement{
		Expression: e,
	}
}

func (es *expressionStatement) String() string {
	return fmt.Sprintf("%v;", es.Expression.String())
}

func (es *expressionStatement) statementFunction() {}

type ifStatement struct {
	token.Metadata
	Condition Expression // Should make a boolean expression type to classify expressions with conditionals
	IfBlock   []Statement
	ElseBlock []Statement
}

func NewIfStatement(t token.Token, c Expression, tb, fb []Statement) Statement {
	return &ifStatement{
		Metadata: t.Data(),
		Condition:         c,
		IfBlock:           tb,
		ElseBlock:         fb,
	}
}

func (ifs *ifStatement) String() string {
	if ifs.ElseBlock != nil {
		return fmt.Sprintf("if %v { %v } else { %v };", ifs.Condition.String(), statementArrayToString(ifs.IfBlock),
			statementArrayToString(ifs.ElseBlock),)
	} else {
		return fmt.Sprintf("if %v { %v };", ifs.Condition.String(), statementArrayToString(ifs.IfBlock),)
	}
}

func (ifs *ifStatement) statementFunction() {}

type functionDeclarationStatement struct {
	token.Metadata
	Name       string
	Parameters []*Identifier
	Body       []Statement
}

func NewFuntionDeclarationStatement(t token.Token, n string, b []Statement, ps ...*Identifier) Statement {
	return &functionDeclarationStatement{
		Metadata: t.Data(),
		Name:              n,
		Parameters:        ps,
		Body:              b,
	}
}

func (fds *functionDeclarationStatement) String() string {
	var identifierNames = make([]string, len(fds.Parameters))
	for i, p := range fds.Parameters {
		identifierNames[i] = p.name
	}
	return fmt.Sprintf("func %v(%v) { %v };", fds.Name, strings.Join(identifierNames, ", "),
		statementArrayToString(fds.Body))
}

func (fds *functionDeclarationStatement) statementFunction() {}


type IdentifierStatement struct {
	*Identifier
}

func NewIdentifierStatement(t token.Token) Statement {
	return &IdentifierStatement{
		Identifier: NewIdentifier(t),
	}
}

func (i *IdentifierStatement) statementFunction() {}

func statementArrayToString(staArr []Statement) string {
	var strArr = make([]string, len(staArr))
	for i, sta := range staArr {
		strArr[i] = sta.String()
	}
	return strings.Join(strArr, ", ")
}

