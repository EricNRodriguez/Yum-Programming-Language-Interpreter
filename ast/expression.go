package ast

import (
	"Yum-Programming-Language-Interpreter/token"
	"fmt"
)

type Prefix struct {
	token.TokenInterface
	parent   NodeInterface
	children [1]NodeInterface
}

func NewPrefixExpression(t token.TokenInterface, e ExpressionInterface, p NodeInterface) ExpressionInterface {
	return &Prefix{
		TokenInterface: t,
		parent:         p,
		children:       [1]NodeInterface{e},
	}
}

func (p *Prefix) String() string {
	return fmt.Sprintf("(%v%v)", p.Literal(), p.children[0].String())
}

func (p *Prefix) Parent() NodeInterface {
	return p.parent
}

func (p *Prefix) Children() []NodeInterface {
	return p.children[0:len(p.children)]
}

func (p *Prefix) expressionFunction() {}

type Infix struct {
	token.TokenInterface
	parent   NodeInterface // nil if expression is root
	children [2]NodeInterface
}

func NewInfixExpression(t token.TokenInterface, le, re ExpressionInterface, p NodeInterface) ExpressionInterface {
	return &Infix{
		TokenInterface: t,
		parent:         p,
		children:       [2]NodeInterface{le, re},
	}
}

func (i *Infix) String() string {
	return fmt.Sprintf("(%v %v %v)", i.children[0].String(), i.TokenInterface.Literal(), i.children[1].String())
}

func (i *Infix) Parent() NodeInterface {
	return i.parent
}

func (i *Infix) Children() []NodeInterface {
	return i.children[0:len(i.children)]
}

func (i *Infix) expressionFunction() {}

type TokenExpression struct {
	token.TokenInterface
	parent NodeInterface
}

func NewTokenExpression(t token.TokenInterface, p NodeInterface) ExpressionInterface {
	return &TokenExpression{
		TokenInterface: t,
		parent:         p,
	}
}

func (te *TokenExpression) String() string {
	return te.TokenInterface.Literal()
}

func (te *TokenExpression) Parent() NodeInterface {
	return te.parent
}

func (te *TokenExpression) Children() []NodeInterface {
	return nil
}

func (te *TokenExpression) expressionFunction() {}
