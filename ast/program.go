package ast

type Program struct {
	Statements []StatementInterface
}

func NewProgram(s ...StatementInterface) *Program {
	return &Program{
		Statements: s,
	}
}

func (p *Program) AddStatment(s StatementInterface) {
	if s != nil {
		p.Statements = append(p.Statements, s)
	}
}

func ()
