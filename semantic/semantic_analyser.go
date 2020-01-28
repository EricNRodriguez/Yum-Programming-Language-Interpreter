package semantic

import (
	"Yum-Programming-Language-Interpreter/ast"
	"Yum-Programming-Language-Interpreter/internal"
	"Yum-Programming-Language-Interpreter/object"
	"Yum-Programming-Language-Interpreter/symbol_table"
	"fmt"
)

type analysisMethod func(node ast.Node)

type SemanticAnalyser interface {
	Analyse(node ast.Node) []error
}

type semanticAnalyser struct {
	symbol_table.SymbolTable
	semanticErrors   []error
	methodRouter     map[ast.NodeType]analysisMethod
	currentStatement ast.NodeType
}

func NewSemanticAnalyser() (sA *semanticAnalyser) {
	sA = &semanticAnalyser{
		SymbolTable:    symbol_table.NewSymbolTable(),
		semanticErrors: make([]error, 0),
		methodRouter:   make(map[ast.NodeType]analysisMethod),
	}

	sA.methodRouter = map[ast.NodeType]analysisMethod{
		ast.ProgramNode:                      sA.analyseProgram,
		ast.PrefixExpressionNode:             sA.analysePrefixExpression,
		ast.InfixExpressionNode:              sA.analyseInfixExpression,
		ast.FunctionCallExpressionNode:       sA.analyseFunctionCallExpression,
		ast.VarStatementNode:                 sA.analyseVarStatement,
		ast.ReturnStatementNode:              sA.analyseReturnStatement,
		ast.IfStatementNode:                  sA.analyseIfStatement,
		ast.WhileStatementNode:               sA.analyseWhileStatement,
		ast.FunctionDeclarationStatementNode: sA.analyseFunctionDeclarationStatement,
		ast.FunctionCallStatementNode:        sA.analyseFunctionCallStatement,
		ast.AssignmentStatementNode:          sA.analyseAssignmentStatement,
		ast.IdentifierExpressionNode:         sA.analyseIdentifierExpression,
		ast.ArrayIndexExpressionNode:         sA.analyseArrayIndexExpression,
		ast.ArrayExpressionNode:              sA.analyseArrayExpression,
	}

	return
}

func (sA *semanticAnalyser) Analyse(node ast.Node) []error {
	sA.analyse(node)
	return sA.semanticErrors
}

func (sA *semanticAnalyser) analyse(node ast.Node) {
	if method, ok := sA.methodRouter[node.Type()]; ok {
		method(node)
	}
	return
}

func (sA *semanticAnalyser) analyseProgram(node ast.Node) {
	for _, s := range node.(*ast.Program).Statements {
		sA.currentStatement = s.Type()
		sA.analyse(s)
	}
	return
}

func (sA *semanticAnalyser) analyseIdentifierExpression(node ast.Node) {
	stmt := node.(*ast.IdentifierExpression)

	if sA.AvailableVar(stmt.Name, true) {
		errMsg := fmt.Sprintf(internal.ErrUndeclaredIdentifierNode, stmt.Name)
		sA.recordError(internal.NewError(stmt.Metadata, errMsg, internal.SemanticErr))
		return
	}

	return
}

func (sA *semanticAnalyser) analysePrefixExpression(node ast.Node) {
	pExpr := node.(*ast.PrefixExpression)
	sA.analyse(pExpr.Expression)
	return
}

func (sA *semanticAnalyser) analyseInfixExpression(node ast.Node) {
	pExpr := node.(*ast.InfixExpression)
	sA.analyse(pExpr.LeftExpression)
	sA.analyse(pExpr.RightExpression)
	return
}

func (sA *semanticAnalyser) analyseArrayIndexExpression(node ast.Node) {
	aIExpr := node.(*ast.ArrayIndexExpression)
	if sA.AvailableVar(aIExpr.ArrayName, true) {
		errMsg := fmt.Sprintf(internal.ErrUndeclaredIdentifierNode, aIExpr.ArrayName)
		sA.recordError(internal.NewError(aIExpr.Metadata, errMsg, internal.SemanticErr))
		return
	}

	if aIExpr.IndexExpr.Type() == ast.ArrayExpressionNode|| aIExpr.IndexExpr.Type() == ast.FloatingPointExpressionNode ||
		aIExpr.IndexExpr.Type() == ast.StringExpressionNode|| aIExpr.IndexExpr.Type() == ast.BooleanExpressionNode {
		errMsg := fmt.Sprintf(internal.ErrInvalidIndexType, aIExpr.IndexExpr.Type())
		sA.recordError(internal.NewError(aIExpr.Metadata, errMsg, internal.SemanticErr))
		return
	}

	sA.Analyse(aIExpr.IndexExpr)

}

func (sA *semanticAnalyser) analyseArrayExpression(node ast.Node) {
	arrExpr := node.(*ast.ArrayExpression)
	for _, expr := range arrExpr.Data {
		if expr.Type() == ast.IdentifierExpressionNode {
			expr := expr.(*ast.IdentifierExpression)
			if sA.AvailableVar(expr.Name, true) {
				errMsg := fmt.Sprintf(internal.ErrUndeclaredFunction, expr.Name)
				sA.recordError(internal.NewError(expr.Metadata, errMsg,
					internal.SemanticErr))
			}
		}
	}
}

func (sA *semanticAnalyser) analyseFunctionCallExpression(node ast.Node) {
	fCall := node.(*ast.FunctionCallExpression)

	if uf, ok := sA.GetUserFunc(fCall.FunctionName); !ok {
		// not a user defined func
		if nf, ok := sA.GetNativeFunc(fCall.FunctionName); !ok {
			// not a function
			errMsg := fmt.Sprintf(internal.ErrUndeclaredFunction, fCall.FunctionName)
			sA.recordError(internal.NewError(fCall.Metadata, errMsg, internal.SemanticErr))
			return

		} else if nf.NumParams != -1 && len(fCall.Parameters) != nf.NumParams {
			// check that native function params align up
			errMsg := fmt.Sprintf(internal.ErrInvalidFunctionCallParameters, fCall.FunctionName, nf.NumParams, len(fCall.Parameters))
			sA.recordError(internal.NewError(fCall.Metadata, errMsg, internal.SemanticErr))
			return

		}

	} else if len(fCall.Parameters) != len(uf.Parameters) {

		// check valid number of params
		errMsg := fmt.Sprintf(internal.ErrInvalidFunctionCallParameters, fCall.FunctionName, len(uf.Parameters), len(fCall.Parameters))
		sA.recordError(internal.NewError(fCall.Metadata, errMsg, internal.SemanticErr))
		return

	}

	sA.analyseBlockExpression(fCall.Parameters...)

	return
}

func (sA *semanticAnalyser) analyseReturnStatement(node ast.Node) {
	rS := node.(*ast.ReturnStatement)
	if !sA.InFunctionCall() {
		sA.recordError(internal.NewError(rS.Data(), internal.ErrReturnLocation, internal.SemanticErr))
	}

	if rS.Expression != nil {
		// analyse return expression
		sA.analyse(rS.Expression)
	}

	return
}

func (sA *semanticAnalyser) analyseIfStatement(node ast.Node) {
	ifStmt := node.(*ast.IfStatement)

	// analyse condition
	sA.analyse(ifStmt.Condition)

	// analyse true block
	sA.EnterScope()
	sA.analyseBlockStatement(ifStmt.IfBlock...)
	sA.ExitScope()

	// analyse false block
	sA.EnterScope()
	sA.analyseBlockStatement(ifStmt.ElseBlock...)
	sA.ExitScope()

	return
}

func (sA *semanticAnalyser) analyseWhileStatement(node ast.Node) {
	wStmt := node.(*ast.WhileStatement)

	// analyse condition
	sA.analyse(wStmt.Condition)

	// analyse true block
	sA.EnterScope()
	sA.analyseBlockStatement(wStmt.Block...)
	sA.ExitScope()

	return
}

// checks that the variable has not been previously declared in the current scope
func (sA *semanticAnalyser) analyseVarStatement(node ast.Node) {
	stmt := node.(*ast.VarStatement)

	// analyse expression
	sA.analyse(stmt.Expression)

	if !sA.AvailableVar(stmt.IdentifierNode.Name, false) {
		errMsg := fmt.Sprintf(internal.ErrDeclaredVariable, stmt.IdentifierNode.Name)
		sA.recordError(internal.NewError(stmt.Metadata, errMsg,
			internal.SemanticErr))
		return
	}

	// save var
	sA.SetVar(stmt.IdentifierNode.Name, object.NewNull())
	return
}

func (sA *semanticAnalyser) analyseAssignmentStatement(node ast.Node) {
	stmt := node.(*ast.AssignmentStatement)
	// check that it exists
	if sA.AvailableVar(stmt.IdentifierNode.Name, true) {
		errMsg := fmt.Sprintf(internal.ErrUndeclaredIdentifierNode, stmt.IdentifierNode.Name)
		sA.recordError(internal.NewError(stmt.Metadata, errMsg, internal.SemanticErr))
		return
	}

	// analyse expression
	sA.analyse(stmt.Expression)

	return
}

func (sA *semanticAnalyser) analyseFunctionDeclarationStatement(node ast.Node) {
	fDec := node.(*ast.FunctionDeclarationStatement)
	if !sA.AvailableFunc(fDec.Name) {
		errMsg := fmt.Sprintf(internal.ErrDeclaredFunction, fDec.Name)
		sA.recordError(internal.NewError(fDec.Metadata, errMsg, internal.SemanticErr))
		return
	}

	// declare func
	sA.SetUserFunc(object.NewUserFunction(fDec.Name, make([]string, len(fDec.Parameters)), []ast.Statement{}))

	// analyse function body
	sA.EnterFunction()
	// declare params in new scope
	for _, p := range fDec.Parameters {
		sA.SetVar(p.Name, object.NewNull())
	}

	sA.analyseBlockStatement(fDec.Body...)
	sA.ExitFunction()

	return
}

func (sA *semanticAnalyser) analyseFunctionCallStatement(node ast.Node) {
	fCallStmt := node.(*ast.FunctionCallStatement)
	sA.analyse(fCallStmt.FunctionCallExpression)
	return
}

func (sA *semanticAnalyser) recordError(err error) {
	if err != nil {
		sA.semanticErrors = append(sA.semanticErrors, err)
	}
	return
}

func (sA *semanticAnalyser) analyseBlockStatement(stmts ...ast.Statement) {
	for _, s := range stmts {
		sA.analyse(s)
	}
	return
}

func (sA *semanticAnalyser) analyseBlockExpression(expr ...ast.Expression) {
	for _, e := range expr {
		sA.analyse(e)
	}
	return
}
