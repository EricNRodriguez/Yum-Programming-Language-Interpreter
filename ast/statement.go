package ast

import (
	"Yum-Programming-Language-Interpreter/token"
	"bytes"
	"fmt"
)

type Var struct {
	token.TokenInterface
	identifier token.TokenInterface
	expression ExpressionInterface
}

func NewVarStatement(t, i token.TokenInterface, e ExpressionInterface) *Var {
	return &Var{
		TokenInterface: t,
		identifier: i,
		expression: e,
	}
}

func (v *Var) String() string {
	return fmt.Sprintf("var %v = %v;", v.identifier.Literal(), v.expression.String())
}

func (v *Var) statementFunction() {}

type Return struct {
	token.TokenInterface
	expression ExpressionInterface
}

func NewReturnStatment(t token.TokenInterface, e ExpressionInterface) *Return {
	return &Return{
		TokenInterface: t,
		expression: e,
	}
}

func (r *Return) String() string {
	return fmt.Sprintf("return %v;", r.expression.String())
}

func (r *Return) statementFunction() {}

type ExpressionStatement struct {
	ExpressionInterface
}

func NewExpressionStatment(e ExpressionInterface) *ExpressionStatement {
	return &ExpressionStatement{
		ExpressionInterface: e,
	}
}

func (es *ExpressionStatement) String() string {
	return fmt.Sprintf("%v;", es.ExpressionInterface.String())
}

func (es *ExpressionStatement) statementFunction() {}

type BlockStatment struct {
	token.MetadataInterface
	Program
}

func NewBlockStatement(t token.TokenInterface) *BlockStatment {
	return &BlockStatment{
		MetadataInterface: t.Metadata(),
		Program: *NewProgram(t.Metadata()),
	}
}

func (bs *BlockStatment) statementFunction() {}


type IfStatement struct {
	token.MetadataInterface
	Condition  ExpressionInterface // Should make a boolean expression type to classify expressions with conditionals
	IfBlock  *BlockStatment
	ElseBlock *BlockStatment
}

func NewIfStatement(t token.TokenInterface, c ExpressionInterface, tb, fb *BlockStatment) *IfStatement {
	return &IfStatement{
		MetadataInterface:      t.Metadata(),
		Condition:  c,
		IfBlock:  tb,
		ElseBlock: fb,
	}
}

func (ifs *IfStatement) String() string {
	if ifs.ElseBlock != nil {
		return fmt.Sprintf("if %v { %v } else { %v };", ifs.Condition.String(), ifs.IfBlock.String(),
			ifs.ElseBlock.String())
	} else {
		return fmt.Sprintf("if %v { %v };", ifs.Condition.String(), ifs.IfBlock.String())
	}
}

func (ifs *IfStatement) statementFunction() {}


type FunctionDeclarationStatement struct {
	token.MetadataInterface
	Name       string
	Parameters []token.TokenInterface
	Body       *BlockStatment
}

func NewFuntionDeclarationStatement(t token.TokenInterface, n string, b *BlockStatment, ps ...token.TokenInterface) *FunctionDeclarationStatement {
	return &FunctionDeclarationStatement{
		MetadataInterface: t.Metadata(),
		Name:       n,
		Parameters: ps,
		Body:       b,
	}
}

func (fds *FunctionDeclarationStatement) String() string {
	pBuff := bytes.Buffer{}
	for i := range fds.Parameters {
		pBuff.WriteString(fds.Parameters[i].Literal())
		if i != len(fds.Parameters)-1 {
			pBuff.WriteString(", ")
		}
	}
	return fmt.Sprintf("func %v(%v) { %v };", fds.Name, pBuff.String(), fds.Body.String())
}

func (fds *FunctionDeclarationStatement) statementFunction() {}
