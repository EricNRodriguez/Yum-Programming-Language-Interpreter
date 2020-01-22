package semantic

import (
	"Yum-Programming-Language-Interpreter/ast"
	"Yum-Programming-Language-Interpreter/internal"
	"Yum-Programming-Language-Interpreter/object"
	"Yum-Programming-Language-Interpreter/symbol_table"
	"fmt"
)

type analysisMethod func(node ast.Node)

type SemanticAnalyser struct {
	symbol_table.SymbolTable
	semanticErrors   []error
	methodRouter     map[ast.NodeType]analysisMethod
	currentStatement ast.NodeType
}

func NewSemanticAnalyser() (sA *SemanticAnalyser) {
	sA = &SemanticAnalyser{
		SymbolTable:    symbol_table.NewSymbolTable(),
		semanticErrors: make([]error, 0),
		methodRouter:   make(map[ast.NodeType]analysisMethod),
	}

	sA.methodRouter = map[ast.NodeType]analysisMethod{
		ast.PROGRAM:                        sA.analyseProgram,
		ast.PREFIX_EXPRESSION:              sA.analysePrefixExpression,
		ast.INFIX_EXPRESSION:               sA.analyseInfixExpression,
		ast.FUNC_CALL_EXPRESSION:           sA.analyseFunctionCallExpression,
		ast.VAR_STATEMENT:                  sA.analyseVarStatement,
		ast.RETURN_STATEMENT:               sA.analyseReturnStatement,
		ast.IF_STATEMENT:                   sA.analyseIfStatement,
		ast.WHILE_STATEMENT: sA.analyseWhileStatement,
		ast.FUNCTION_DECLARATION_STATEMENT: sA.analyseFunctionDeclarationStatement,
		ast.FUNCTION_CALL_STATEMENT:        sA.analyseFunctionCallStatement,
		ast.ASSIGNMENT_STATEMENT:           sA.analyseAssignmentStatement,
		ast.IDENTIFIER_EXPRESSION:          sA.analyseIdentifierExpression,
		ast.IMPORT_STATEMENT: sA.analyseImportStatement,
	}

	return
}

func (sA *SemanticAnalyser) Analyse(node ast.Node) {
	if method, ok := sA.methodRouter[node.Type()]; ok {
		method(node)
	}
	return
}

func (sA *SemanticAnalyser) analyseProgram(node ast.Node) {
	for _, s := range node.(*ast.Program).Statements {
		sA.currentStatement = s.Type()
		sA.Analyse(s)
	}
	return
}

func (sA *SemanticAnalyser) analyseIdentifierExpression(node ast.Node) {
	stmt := node.(*ast.IdentifierExpression)

	if sA.AvailableVar(stmt.Name, true) {
		errMsg := fmt.Sprintf(internal.UndeclaredIdentifierErr, stmt.Name)
		sA.recordError(internal.NewError(stmt.Metadata, errMsg, internal.SemanticErr))
		return
	}

	return
}

func (sA *SemanticAnalyser) analysePrefixExpression(node ast.Node) {
	pExpr := node.(*ast.PrefixExpression)
	sA.Analyse(pExpr.Expression)
	return
}

func (sA *SemanticAnalyser) analyseInfixExpression(node ast.Node) {
	pExpr := node.(*ast.InfixExpression)
	sA.Analyse(pExpr.LeftExpression)
	sA.Analyse(pExpr.RightExpression)
	return
}

func (sA *SemanticAnalyser) analyseFunctionCallExpression(node ast.Node) {
	fCall := node.(*ast.FunctionCallExpression)

	if f, ok := sA.GetUserFunc(fCall.FunctionName); !ok {
		// not a user defined func
		if _, ok := sA.GetNativeFunc(fCall.FunctionName); !ok {
			// not a function
			errMsg := fmt.Sprintf(internal.UndeclaredFunctionErr, fCall.FunctionName)
			sA.recordError(internal.NewError(fCall.Metadata, errMsg, internal.SemanticErr))
			return
		} else {
			// check that native function params align up!
		}
	} else {
		if len(fCall.Parameters) != len(f.Parameters) {
			errMsg := fmt.Sprintf(internal.InvalidFunctionCallParametersErr, fCall.FunctionName, len(f.Parameters), len(fCall.Parameters))
			sA.recordError(internal.NewError(fCall.Metadata, errMsg, internal.SemanticErr))
			return
		}
	}

	sA.analyseBlockExpression(fCall.Parameters...)

	return
}

func (sA *SemanticAnalyser) analyseReturnStatement(node ast.Node) {
	rS := node.(*ast.ReturnStatement)
	if !sA.InFunctionCall() {
		sA.recordError(internal.NewError(rS.Data(), internal.ReturnLocationErr, internal.SemanticErr))
	}

	// analyse return expression
	sA.Analyse(rS.Expression)
	return
}

func (sA *SemanticAnalyser) analyseImportStatement(node ast.Node) {
	iStmt := node.(*ast.ImportStatement)

	if !sA.AvailableFunc(iStmt.ImportFunctionName) {
		errMsg := fmt.Sprintf(internal.DeclaredFunctionErr, iStmt.ImportFunctionName)
		sA.recordError(internal.NewError(iStmt.Metadata, errMsg, internal.SemanticErr))
		return
	}

	// declare func name
	sA.SetUserFunc(object.NewUserFunction(iStmt.ImportFunctionName, make([]string, 0), []ast.Statement{}))

	return
}

func (sA *SemanticAnalyser) analyseIfStatement(node ast.Node) {
	ifStmt := node.(*ast.IfStatement)

	// analyse condition
	sA.Analyse(ifStmt.Condition)

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

func (sA *SemanticAnalyser) analyseWhileStatement(node ast.Node) {
	wStmt := node.(*ast.WhileStatement)

	// analyse condition
	sA.Analyse(wStmt.Condition)

	// analyse true block
	sA.EnterScope()
	sA.analyseBlockStatement(wStmt.Block...)
	sA.ExitScope()

	return
}

// checks that the variable has not been previously declared in the current scope
func (sA *SemanticAnalyser) analyseVarStatement(node ast.Node) {
	stmt := node.(*ast.VarStatement)

	// analyse expression
	sA.Analyse(stmt.Expression)

	if !sA.AvailableVar(stmt.Identifier.Name, false) {
		errMsg := fmt.Sprintf(internal.DeclaredVariableErr, stmt.Identifier.Name)
		sA.recordError(internal.NewError(stmt.Metadata, errMsg,
			internal.SemanticErr))
		return
	}

	// save var
	sA.SetVar(stmt.Identifier.Name, object.NewNull())
	return
}

func (sA *SemanticAnalyser) analyseAssignmentStatement(node ast.Node) {
	stmt := node.(*ast.AssignmentStatement)
	// check that it exists
	if sA.AvailableVar(stmt.Identifier.Name, true) {
		errMsg := fmt.Sprintf(internal.UndeclaredIdentifierErr, stmt.Identifier.Name)
		sA.recordError(internal.NewError(stmt.Metadata, errMsg, internal.SemanticErr))
		return
	}

	// analyse expression
	sA.Analyse(stmt.Expression)

	return
}

func (sA *SemanticAnalyser) analyseFunctionDeclarationStatement(node ast.Node) {
	fDec := node.(*ast.FunctionDeclarationStatement)
	if !sA.AvailableFunc(fDec.Name) {
		errMsg := fmt.Sprintf(internal.DeclaredFunctionErr, fDec.Name)
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

func (sA *SemanticAnalyser) analyseFunctionCallStatement(node ast.Node) {
	fCallStmt := node.(*ast.FunctionCallStatement)
	sA.Analyse(fCallStmt.FunctionCallExpression)
	return
}

func (sA *SemanticAnalyser) recordError(err error) {
	if err != nil {
		sA.semanticErrors = append(sA.semanticErrors, err)
	}
	return
}

func (sA *SemanticAnalyser) SemanticErrors() []error {
	return sA.semanticErrors
}

func (sA *SemanticAnalyser) analyseBlockStatement(stmts ...ast.Statement) {
	for _, s := range stmts {
		sA.Analyse(s)
	}
	return
}

func (sA *SemanticAnalyser) analyseBlockExpression(expr ...ast.Expression) {
	for _, e := range expr {
		sA.Analyse(e)
	}
	return
}
