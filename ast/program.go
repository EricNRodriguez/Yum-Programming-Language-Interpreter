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
		Metadata:   m,
		Statements: s,
	}
}

// moves imports, followed by func declarations to the start of the program
func (p *Program) Hoist() {
	var (
		hoistedStatementsDecs = make([]Statement, 0)
		hoistedImportStatements = make([]Statement, 0)
		remainingStatements = make([]Statement, 0)
	)

	for i := range p.Statements {
		switch p.Statements[i].Type() {
		case FUNCTION_DECLARATION_STATEMENT:
			hoistedStatementsDecs = append(hoistedStatementsDecs, p.Statements[i])
		case IMPORT_STATEMENT:
			hoistedImportStatements = append(hoistedImportStatements, p.Statements[i])
		default:
			remainingStatements = append(remainingStatements, p.Statements[i])
		}
	}

	p.Statements = append(append(hoistedImportStatements, hoistedStatementsDecs...), remainingStatements...)
	return
}

func (p *Program) String() string {
	lBuff := bytes.Buffer{}
	for _, s := range p.Statements {
		lBuff.WriteString(s.String())
		lBuff.WriteString("\n")
	}

	return lBuff.String()
}

func (p *Program) Type() NodeType {
	return PROGRAM
}
