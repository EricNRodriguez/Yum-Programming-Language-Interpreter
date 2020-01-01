package ast

import "Yum-Programming-Language-Interpreter/token"

type Identifier struct {
	token.Metadata
	name string
}

func NewIdentifier(t token.Token) *Identifier {
	return &Identifier{
		Metadata: t.Data(),
		name: t.Literal(),
	}
}

func (i *Identifier) String() string {
	return i.name
}



