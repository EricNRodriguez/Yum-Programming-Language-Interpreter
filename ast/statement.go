package ast

import (
	"github.com/EricNRodriguez/yum/token"
	"bytes"
	"fmt"
	"strings"
)

type VarStatement struct {
	*AssignmentStatement
}

func NewVarStatement(md token.Metadata, i *IdentifierExpression, e Expression) *VarStatement {
	return &VarStatement{
		AssignmentStatement: NewAssignmentStatement(md, i, e),
	}
}

func (v *VarStatement) String() string {
	return fmt.Sprintf("var %v", v.AssignmentStatement.String())
}

func (v *VarStatement) Type() NodeType {
	return VarStatementNode
}

type AssignmentStatement struct {
	token.Metadata
	IdentifierNode *IdentifierExpression
	Expression     Expression
}

func NewAssignmentStatement(md token.Metadata, i *IdentifierExpression, e Expression) *AssignmentStatement {
	return &AssignmentStatement{
		Metadata:       md,
		IdentifierNode: i,
		Expression:     e,
	}
}

func (as *AssignmentStatement) String() string {
	return fmt.Sprintf("%v = %v;", as.IdentifierNode.String(), as.Expression.String())
}

func (as *AssignmentStatement) Type() NodeType {
	return AssignmentStatementNode
}

func (as *AssignmentStatement) statementFunction() {}

type ReturnStatement struct {
	token.Token
	Expression Expression
}

func NewReturnStatment(t token.Token, e Expression) *ReturnStatement {
	return &ReturnStatement{
		Token:      t,
		Expression: e,
	}
}

func (r *ReturnStatement) String() string {
	if r.Expression == nil {
		return fmt.Sprintf("return;")
	} else {
		return fmt.Sprintf("return %v;", r.Expression.String())
	}
}

func (r *ReturnStatement) Type() NodeType {
	return ReturnStatementNode
}

func (r *ReturnStatement) statementFunction() {}

type FunctionCallStatement struct {
	*FunctionCallExpression
}

func NewFunctionCallStatement(e *FunctionCallExpression) *FunctionCallStatement {
	return &FunctionCallStatement{
		FunctionCallExpression: e,
	}
}

func (fc *FunctionCallStatement) String() string {
	return fmt.Sprintf("%v;", fc.FunctionCallExpression.String())
}

func (fc *FunctionCallStatement) Type() NodeType {
	return FunctionCallStatementNode
}

func (fc *FunctionCallStatement) statementFunction() {}

type IfStatement struct {
	token.Metadata
	Condition Expression // Should make a boolean expression type to classify expressions with conditionals
	IfBlock   []Statement
	ElseBlock []Statement
}

func NewIfStatement(t token.Token, c Expression, tb, fb []Statement) Statement {
	return &IfStatement{
		Metadata:  t.Data(),
		Condition: c,
		IfBlock:   tb,
		ElseBlock: fb,
	}
}

func (ifs *IfStatement) String() string {
	if ifs.ElseBlock != nil {
		return fmt.Sprintf("if (%v) { %v } else { %v };", ifs.Condition.String(), statementArrayNodeToString(ifs.IfBlock),
			statementArrayNodeToString(ifs.ElseBlock))
	} else {
		return fmt.Sprintf("if (%v) { %v };", ifs.Condition.String(), statementArrayNodeToString(ifs.IfBlock))
	}
}

func (ifs *IfStatement) Type() NodeType {
	return IfStatementNode
}

func (ifs *IfStatement) statementFunction() {}

type WhileStatement struct {
	token.Metadata
	Condition Expression // Should make a boolean expression type to classify expressions with conditionals
	Block     []Statement
}

func NewWhileStatement(md token.Metadata, c Expression, b []Statement) *WhileStatement {
	return &WhileStatement{
		Metadata:  md,
		Condition: c,
		Block:     b,
	}
}

func (w *WhileStatement) String() string {
	return fmt.Sprintf("while (%v) { %v };", w.Condition.String(), statementArrayNodeToString(w.Block))
}

func (w *WhileStatement) Type() NodeType {
	return WhileStatementNode
}

func (w *WhileStatement) statementFunction() {}

type FunctionDeclarationStatement struct {
	token.Metadata
	Name       string
	Parameters []IdentifierExpression
	Body       []Statement
}

func NewFuntionDeclarationStatement(t token.Token, n string, b []Statement, ps []IdentifierExpression) Statement {
	return &FunctionDeclarationStatement{
		Metadata:   t.Data(),
		Name:       n,
		Parameters: ps,
		Body:       b,
	}
}

func (fds *FunctionDeclarationStatement) String() string {
	var IdentifierNodeNames = make([]string, len(fds.Parameters))
	for i, p := range fds.Parameters {
		IdentifierNodeNames[i] = p.String()
	}
	return fmt.Sprintf("func %v(%v) { %v };", fds.Name, strings.Join(IdentifierNodeNames, ", "),
		statementArrayNodeToString(fds.Body))
}

func (fds *FunctionDeclarationStatement) Type() NodeType {
	return FunctionDeclarationStatementNode
}

func (fds *FunctionDeclarationStatement) statementFunction() {}

func statementArrayNodeToString(staArr []Statement) string {
	strBuff := bytes.Buffer{}

	for _, sta := range staArr {
		strBuff.WriteString(sta.String())
	}

	return strBuff.String()
}
