package ast

import (
	"Yum-Programming-Language-Interpreter/token"
	"fmt"
	"strings"
)

type VarStatement struct {
	*AssignmentStatement
}

func NewVarStatement(md token.Metadata, i *Identifier, e Expression) *VarStatement {
	return &VarStatement{
		AssignmentStatement: NewAssignmentStatement(md, i, e),
	}
}

func (v *VarStatement) String() string {
	return fmt.Sprintf("var %v", v.AssignmentStatement.String())
}

func (v *VarStatement) Type() NodeType {
	return VAR_STATEMENT
}

type AssignmentStatement struct {
	token.Metadata
	Identifier *Identifier
	Expression Expression
}

func NewAssignmentStatement(md token.Metadata, i *Identifier, e Expression) *AssignmentStatement {
	return &AssignmentStatement{
		Metadata:   md,
		Identifier: i,
		Expression: e,
	}
}

func (as *AssignmentStatement) String() string {
	return fmt.Sprintf(" %v = %v;", as.Identifier.String(), as.Expression.String())
}

func (as *AssignmentStatement) Type() NodeType {
	return ASSIGNMENT_STATEMENT
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
	return fmt.Sprintf("return %v;", r.Expression.String())
}

func (r *ReturnStatement) Type() NodeType {
	return RETURN_STATEMENT
}

func (r *ReturnStatement) statementFunction() {}

type ExpressionStatement struct {
	Expression
}

func NewExpressionStatment(e Expression) Statement {
	return &ExpressionStatement{
		Expression: e,
	}
}

func (es *ExpressionStatement) String() string {
	return fmt.Sprintf("%v;", es.Expression.String())
}

func (es *ExpressionStatement) Type() NodeType {
	return EXPRESSION_STATEMENT
}

func (es *ExpressionStatement) statementFunction() {}

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
		return fmt.Sprintf("if %v { %v } else { %v };", ifs.Condition.String(), statementArrayToString(ifs.IfBlock),
			statementArrayToString(ifs.ElseBlock))
	} else {
		return fmt.Sprintf("if %v { %v };", ifs.Condition.String(), statementArrayToString(ifs.IfBlock))
	}
}

func (ifs *IfStatement) Type() NodeType {
	return IF_STATEMENT
}

func (ifs *IfStatement) statementFunction() {}

type FunctionDeclarationStatement struct {
	token.Metadata
	Name       string
	Parameters []Identifier
	Body       []Statement
}

func NewFuntionDeclarationStatement(t token.Token, n string, b []Statement, ps []Identifier) Statement {
	return &FunctionDeclarationStatement{
		Metadata:   t.Data(),
		Name:       n,
		Parameters: ps,
		Body:       b,
	}
}

func (fds *FunctionDeclarationStatement) String() string {
	var identifierNames = make([]string, len(fds.Parameters))
	for i, p := range fds.Parameters {
		identifierNames[i] = p.String()
	}
	return fmt.Sprintf("func %v(%v) { %v };", fds.Name, strings.Join(identifierNames, ", "),
		statementArrayToString(fds.Body))
}

func (fds *FunctionDeclarationStatement) Type() NodeType {
	return FUNCTION_DECLARATION_STATEMENT
}

func (fds *FunctionDeclarationStatement) statementFunction() {}
//
//type FunctionCallStatement struct {
//	token.Metadata
//	Name       string
//	Parameters []Expression
//}
//
//func NewFunctionCallStatement(md token.Metadata, n string, ps []Expression) *FunctionCallStatement {
//	return &FunctionCallStatement{
//		Metadata:   md,
//		Name:       n,
//		Parameters: ps,
//	}
//}
//
//func (fcs *FunctionCallStatement) String() string {
//	return fmt.Sprintf("%v(%v);", fcs.Name, expressionArrayToString(fcs.Parameters))
//}
//
//func (fcs *FunctionCallStatement) Type() NodeType {
//	return FUNCTION_CALL_STATEMENT
//}
//
//func (fcs *FunctionCallStatement) statementFunction() {}

func statementArrayToString(staArr []Statement) string {
	var strArr = make([]string, len(staArr))
	for i, sta := range staArr {
		strArr[i] = sta.String()
	}
	return strings.Join(strArr, ", ")
}


