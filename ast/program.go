package ast

import (
	"Yum-Programming-Language-Interpreter/token"
	"bytes"
)

type Program struct {
	token.Metadata
	Statements []Statement
}

func NewProgram(m token.Metadata, s ...Statement) *Program {
	return &Program{
		Metadata: m,
		Statements:        s,
	}
}

func (p *Program) AddStatement(s Statement) {
	if s != nil {
		p.Statements = append(p.Statements, s)
	}
}

func (p *Program) String() string {
	lBuff := bytes.Buffer{}
	for _, s := range p.Statements {
		lBuff.WriteString(s.String())
		lBuff.WriteString(" ")
	}

	return lBuff.String()
}
