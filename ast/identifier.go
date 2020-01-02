package ast

import (
	"Yum-Programming-Language-Interpreter/token"
)

type Identifier struct {
	token.Metadata
	Name string
}

func NewIdentifier(t token.Token) *Identifier {
	return &Identifier{
		Metadata: t.Data(),
		Name:     t.Literal(),
	}
}

func (i *Identifier) String() string {
	return i.Name
}

func (i *Identifier) Type() NodeType {
	return IDENTIFIER
}
