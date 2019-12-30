package ast

import (
	"Yum-Programming-Language-Interpreter/token"
	"fmt"
)

type Var struct {
	token.TokenInterface
	parent NodeInterface
	children [2]NodeInterface
}

func NewVarStatement(t token.TokenInterface, p NodeInterface, id, v NodeInterface) StatementInterface {
	return &Var{
		TokenInterface: t,
		parent: p,
		children: [2]NodeInterface{id,v},
	}
}

func (v *Var) String() string {
	return fmt.Sprintf("var %v = %v;", v.children[0].Name, v.children[1].String())
}

func (v *Var) Parent() NodeInterface {
	return v.parent
}

func (v *Var) Children() []NodeInterface {
	return v.children[0:len(v.children)]
}

func (v *Var) statementFunction() {}

type Return struct {
	token.TokenInterface
	parent NodeInterface
	children [1]NodeInterface
}

func NewReturnStatment(t token.TokenInterface, p NodeInterface, e ExpressionInterface) StatementInterface {
	return &Return{
		TokenInterface: t,
		parent: p,
		children: [1]NodeInterface{e},
	}
}

func (r *Return) String() string {
	return fmt.Sprintf("return %v;", r.children[0].Literal())
}

func (r *Return) Parent() NodeInterface {
	return r.parent
}

func (r *Return) Children() []NodeInterface {
	return r.children[0:len(r.children)]
}

func (r *Return) statementFunction() {}

type ExpressionStatement struct {
	ExpressionInterface
}

func NewExpressionStatment(e ExpressionInterface) StatementInterface {
	return &ExpressionStatement{
		ExpressionInterface: e,
	}
}

func (es *ExpressionStatement) statementFunction() {}