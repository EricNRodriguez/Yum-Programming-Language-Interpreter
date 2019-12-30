package ast

import (
	"Yum-Programming-Language-Interpreter/token"
	"bytes"
)

type Program struct {
	token.MetadataInterface
	Statements []StatementInterface
}

func NewProgram(m token.MetadataInterface, s ...StatementInterface) *Program {
	return &Program{
		MetadataInterface: m,
		Statements: s,
	}
}

func (p *Program) AddStatement(s StatementInterface) {
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
