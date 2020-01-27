package eval

import (
	"Yum-Programming-Language-Interpreter/ast"
	"Yum-Programming-Language-Interpreter/internal"
	"Yum-Programming-Language-Interpreter/object"
	"Yum-Programming-Language-Interpreter/symbol_table"
	"Yum-Programming-Language-Interpreter/token"

	//"Yum-ProgramNodeming-Language-Interpreter/ast"
	//"Yum-ProgramNodeming-Language-Interpreter/internal"
	//"Yum-ProgramNodeming-Language-Interpreter/object"
	//"Yum-ProgramNodeming-Language-Interpreter/symbol_table"
	//"Yum-ProgramNodeming-Language-Interpreter/token"
	"fmt"
	"os"
)

type Evaluator interface {
	Evaluate(ast.Node)
}

type evalMethod func(node ast.Node) object.Object

type evaluator struct {
	stackTrace   internal.StackTrace
	symbolTable  symbol_table.SymbolTable
	methodRouter map[ast.NodeType]evalMethod
}

func NewEvaluator() (e *evaluator) {
	e = &evaluator{
		symbolTable: symbol_table.NewSymbolTable(),
		stackTrace:  internal.NewStackTrace(),
	}

	e.methodRouter = map[ast.NodeType]evalMethod{
		ast.ProgramNode:                      e.evaluateProgramNode,
		ast.ArrayExpressionNode:              e.evaluateArrayExpression,
		ast.ArrayIndexExpressionNode:         e.evaluateArrayIndexExpression,
		ast.PrefixExpressionNode:             e.evaluatePrefixExpression,
		ast.InfixExpressionNode:              e.evaluateInfixExpression,
		ast.IntegerExpressionNode:            e.evaluateIntegerExpression,
		ast.FloatingPointExpressionNode:      e.evaluateFloatingPointExpression,
		ast.StringExpressionNode:             e.evaluateStringExpression,
		ast.BooleanExpressionNode:            e.evaluateBooleanExpression,
		ast.FunctionCallExpressionNode:       e.evaluateFunctionCallExpression,
		ast.IdentifierExpressionNode:         e.evaluateIdentifierExpression,
		ast.VarStatementNode:                 e.evaluateVarStatement,
		ast.ReturnStatementNode:              e.evaluateReturnStatement,
		ast.IfStatementNode:                  e.evaluateIfStatement,
		ast.WhileStatementNode:               e.evaluateWhileStatement,
		ast.FunctionDeclarationStatementNode: e.evaluateFunctionDeclarationStatement,
		ast.FunctionCallStatementNode:        e.evaluateFunctionCallStatement,
		ast.AssignmentStatementNode:          e.evaluateAssignmentStatement,
	}

	return
}

func (e *evaluator) Evaluate(node ast.Node) {
	e.evaluate(node)
	return
}

func (e *evaluator) evaluate(node ast.Node) (o object.Object) {
	if method, ok := e.methodRouter[node.Type()]; ok {
		return method(node)
	}

	e.quit(internal.NewError(token.NewMetatadata(node.LineNumber(), node.FileName()),
		fmt.Sprintf(internal.ErrUnimplementedType, node.Type()), internal.InternalErr))
	return nil
}

func (e *evaluator) evaluateProgramNode(node ast.Node) (o object.Object) {
	prog := node.(*ast.Program)
	return e.evaluateBlockStatement(prog.Statements...)
}

func (e *evaluator) evaluateIdentifierExpression(node ast.Node) (o object.Object) {
	iden := node.(*ast.IdentifierExpression)
	o, _ = e.symbolTable.GetVar(iden.Name)
	return
}

func (e *evaluator) evaluatePrefixExpression(node ast.Node) (o object.Object) {
	pExpr := node.(*ast.PrefixExpression)
	rObj := e.evaluate(pExpr.Expression)

	if rObj.Type() == object.IntegerObject {
		rObj := rObj.(*object.Integer)

		switch pExpr.Token.Type() {
		case token.AddToken:
			o = rObj
		case token.SubToken:
			o = object.NewInteger(-1 * rObj.Value)
		default:
			e.quit(internal.NewError(pExpr.Data(), fmt.Sprintf(internal.ErrType, rObj.Literal(), object.BooleanObject),
				internal.RuntimeErr))
		}

	} else if rObj.Type() == object.FloatingPointObject {
		rObj := rObj.(*object.Float)

		switch pExpr.Token.Type() {
		case token.AddToken:
			o = rObj
		case token.SubToken:
			o = object.NewFloat(-1 * rObj.Value)
		default:
			e.quit(internal.NewError(pExpr.Data(), fmt.Sprintf(internal.ErrType, rObj.Literal(), object.BooleanObject),
				internal.RuntimeErr))
		}
	} else if rObj.Type() == object.BooleanObject {

		rObj := rObj.(*object.Boolean)
		switch pExpr.Token.Type() {
		case token.NegateToken:
			o = object.NewBoolean(!rObj.Value)
		default:
			e.quit(internal.NewError(pExpr.Data(), fmt.Sprintf(internal.ErrType, rObj.Literal(),
				fmt.Sprintf("%v or %v", object.IntegerObject, object.FloatingPointObject)), internal.RuntimeErr))
		}

	} else {
		// null object
		e.quit(internal.NewError(pExpr.Data(), fmt.Sprintf(internal.ErrType, rObj.Literal(),
			fmt.Sprintf("%v or %v or %v", object.IntegerObject, object.FloatingPointObject, object.BooleanObject)), internal.RuntimeErr))
	}
	return
}

func (e *evaluator) evaluateInfixExpression(node ast.Node) (o object.Object) {
	iExpr := node.(*ast.InfixExpression)
	lObj := e.unpack(e.evaluate(iExpr.LeftExpression))
	rObj := e.unpack(e.evaluate(iExpr.RightExpression))

	if lObj.Type() == object.IntegerObject && rObj.Type() == object.IntegerObject {
		lObj := lObj.(*object.Integer)
		rObj := rObj.(*object.Integer)

		switch iExpr.Token.Type() {
		case token.AddToken:
			o = object.NewInteger(lObj.Value + rObj.Value)
		case token.SubToken:
			o = object.NewInteger(lObj.Value - rObj.Value)
		case token.DivToken:
			if rObj.Value == 0 {
				e.quit(internal.NewError(iExpr.Data(), internal.ErrDivisionByZero, internal.RuntimeErr))
			}

			o = object.NewInteger(lObj.Value / rObj.Value)
		case token.MultToken:
			o = object.NewInteger(lObj.Value * rObj.Value)
		case token.GThanToken:
			o = object.NewBoolean(lObj.Value > rObj.Value)
		case token.LThanToken:
			o = object.NewBoolean(lObj.Value < rObj.Value)
		case token.GThanEqualToken:
			o = object.NewBoolean(lObj.Value >= rObj.Value)
		case token.LThanEqualToken:
			o = object.NewBoolean(lObj.Value <= rObj.Value)
		case token.EqualToken:
			o = object.NewBoolean(lObj.Value == rObj.Value)
		case token.NotEqualToken:
			o = object.NewBoolean(lObj.Value != rObj.Value)
		default:
			e.quit(internal.NewError(iExpr.Data(), fmt.Sprintf(internal.ErrTypeOperation, iExpr.Token.Type(),
				lObj.Type()), internal.RuntimeErr))
			o = object.NewNull()
		}
	} else if (lObj.Type() == object.IntegerObject || lObj.Type() == object.FloatingPointObject) &&
		(rObj.Type() == object.IntegerObject || rObj.Type() == object.FloatingPointObject) {

		// type cast to floating point numbers
		lObj := e.castToFloat(lObj)
		rObj := e.castToFloat(rObj)

		switch iExpr.Token.Type() {
		case token.AddToken:
			o = object.NewFloat(lObj.Value + rObj.Value)
		case token.SubToken:
			o = object.NewFloat(lObj.Value - rObj.Value)
		case token.DivToken:
			if rObj.Value == 0 {
				e.quit(internal.NewError(iExpr.Data(), internal.ErrDivisionByZero, internal.RuntimeErr))
			}

			o = object.NewFloat(lObj.Value / rObj.Value)
		case token.MultToken:
			o = object.NewFloat(lObj.Value * rObj.Value)
		case token.GThanToken:
			o = object.NewBoolean(lObj.Value > rObj.Value)
		case token.LThanToken:
			o = object.NewBoolean(lObj.Value < rObj.Value)
		case token.GThanEqualToken:
			o = object.NewBoolean(lObj.Value >= rObj.Value)
		case token.LThanEqualToken:
			o = object.NewBoolean(lObj.Value <= rObj.Value)
		case token.EqualToken:
			o = object.NewBoolean(lObj.Value == rObj.Value)
		case token.NotEqualToken:
			o = object.NewBoolean(lObj.Value != rObj.Value)
		default:
			e.quit(internal.NewError(iExpr.Data(), fmt.Sprintf(internal.ErrTypeOperation, iExpr.Token.Type(),
				lObj.Type()), internal.RuntimeErr))
			o = object.NewNull()
		}

	} else if lObj.Type() == object.BooleanObject && rObj.Type() == object.BooleanObject {
		lObj := lObj.(*object.Boolean)
		rObj := rObj.(*object.Boolean)
		switch iExpr.Token.Type() {
		case token.EqualToken:
			o = object.NewBoolean(lObj.Value == rObj.Value)
		case token.NotEqualToken:
			o = object.NewBoolean(lObj.Value != rObj.Value)
		case token.AndToken:
			o = object.NewBoolean(lObj.Value && rObj.Value)
		case token.OrToken:
			o = object.NewBoolean(lObj.Value || rObj.Value)
		default:
			e.quit(internal.NewError(iExpr.Data(), fmt.Sprintf(internal.ErrTypeOperation, iExpr.Token.Type(),
				lObj.Type()), internal.RuntimeErr))
		}
	} else if lObj.Type() == object.StringObject && rObj.Type() == object.StringObject {
		lObj := lObj.(*object.String)
		rObj := rObj.(*object.String)
		switch iExpr.Token.Type() {
		case token.EqualToken:
			o = object.NewBoolean(lObj.Literal() == rObj.Literal())
		case token.AddToken:
			o = object.NewString(lObj.Lit + rObj.Lit)
		default:
			e.quit(internal.NewError(iExpr.Data(), fmt.Sprintf(internal.ErrTypeOperation, iExpr.Token.Type(),
				lObj.Type()), internal.RuntimeErr))
		}

	} else {
		e.quit(internal.NewError(iExpr.Data(), fmt.Sprintf(internal.ErrTypeOperation, iExpr.Literal(),
			fmt.Sprintf("%v and %v", lObj.Type(), rObj.Type())),
			internal.RuntimeErr))
	}

	return
}

// must receive either an integer or a float, will panic otherwise
func (e *evaluator) castToFloat(obj object.Object) (o *object.Float) {
	var ok bool
	if o, ok = obj.(*object.Float); !ok {
		o = object.NewFloat(float64(obj.(*object.Integer).Value))
	}
	return o
}

func (e *evaluator) unpack(o object.Object) object.Object {
	if o.Type() != object.ReturnObject {
		return o
	}
	oV := o.(*object.ReturnValue).Value
	return e.unpack(oV)
}

func (e *evaluator) evaluateIntegerExpression(node ast.Node) object.Object {
	i := node.(*ast.IntegerExpression)
	o := object.NewInteger(i.Value)
	return o
}

func (e *evaluator) evaluateStringExpression(node ast.Node) object.Object {
	s := node.(*ast.StringExpression)
	o := object.NewString(s.Literal)
	return o
}

func (e *evaluator) evaluateArrayExpression(node ast.Node) object.Object {
	a := node.(*ast.ArrayExpression)
	oData := make([]object.Object, a.Length)
	for i, ex := range a.Data {
		oData[i] = e.evaluate(ex)
	}
	o := object.NewArrayNode(oData)
	return o
}

func (e *evaluator) evaluateArrayIndexExpression(node ast.Node) (o object.Object) {
	iden := node.(*ast.ArrayIndexExpression)

	arrE, _ := e.symbolTable.GetVar(iden.ArrayName)
	if arrE.Type() != object.ArrayObject {
		errMsg := fmt.Sprintf(internal.ErrType, arrE.Literal(), object.ArrayObject)
		e.quit(internal.NewError(iden.Metadata, errMsg, internal.RuntimeErr))
	}
	arr := arrE.(*object.ArrayNode)

	indE := e.evaluate(iden.IndexExpr)
	if indE.Type() != object.IntegerObject {
		errMsg := fmt.Sprintf(internal.ErrType, arrE.Literal(), object.IntegerObject)
		e.quit(internal.NewError(iden.Metadata, errMsg, internal.RuntimeErr))
	}

	indI := indE.(*object.Integer)
	if indI.Value > arr.Length-1 || indI.Value < 0 {
		// index out of bounds error
		e.quit(internal.NewError(iden.Metadata, internal.ErrIndexOutOfBounds, internal.RuntimeErr))
	}

	o = arr.Data[indI.Value]
	return
}

func (e *evaluator) evaluateFloatingPointExpression(node ast.Node) object.Object {
	i := node.(*ast.FloatingPointExpression)
	o := object.NewFloat(i.Value)
	return o
}

func (e *evaluator) evaluateBooleanExpression(node ast.Node) object.Object {
	b := node.(*ast.BooleanExpression)
	o := object.NewBoolean(b.Value)
	return o
}

func (e *evaluator) evaluateFunctionCallExpression(node ast.Node) (o object.Object) {
	var (
		fCall = node.(*ast.FunctionCallExpression)
		err   error
	)

	e.stackTrace.Push(fCall) // record function call

	// native function call
	if f, ok := e.symbolTable.GetUserFunc(fCall.FunctionName); !ok {
		// evaluate parameters
		evalParams := make([]object.Object, len(fCall.Parameters))
		for i, expr := range fCall.Parameters {
			evalParams[i] = e.evaluate(expr)
		}

		f, _ := e.symbolTable.GetNativeFunc(fCall.FunctionName)
		if o, err = f.Function(evalParams...); err != nil {
			e.quit(internal.NewError(fCall.Metadata, err.Error(), internal.RuntimeErr))
		}

		// user defined function
	} else {

		// evaluate parameters
		paramValues := map[string]object.Object{}
		for i, v := range fCall.Parameters {
			paramValues[f.Parameters[i]] = e.evaluate(v)
		}

		// new symbol table
		e.symbolTable.EnterFunction()

		for k, v := range paramValues {
			e.symbolTable.SetVar(k, v)
		}

		o = e.evaluateBlockStatement(f.Body...)
		if o == nil || o.Type() != object.ReturnObject {
			o = object.NewNull()
		}
		e.symbolTable.ExitFunction()

	}

	e.stackTrace.Pop()
	if o.Type() == object.ReturnObject {
		o = o.(*object.ReturnValue).Value
	}
	return o
}

func (e *evaluator) evaluateFunctionCallStatement(node ast.Node) (o object.Object) {
	stmt := node.(*ast.FunctionCallStatement)
	e.evaluateFunctionCallExpression(stmt.FunctionCallExpression)
	return
}

func (e *evaluator) evaluateVarStatement(node ast.Node) object.Object {
	vStmt := node.(*ast.VarStatement)
	value := e.evaluate(vStmt.Expression)
	if value.Type() == object.ReturnObject {
		value = value.(*object.ReturnValue).Value
	}
	e.symbolTable.SetVar(vStmt.IdentifierNode.String(), value)
	return object.NewNull()
}

func (e *evaluator) evaluateAssignmentStatement(node ast.Node) object.Object {
	vStmt := node.(*ast.AssignmentStatement)
	leftObj := e.evaluate(vStmt.Expression)
	e.symbolTable.UpdateVar(vStmt.IdentifierNode.String(), leftObj)
	return object.NewNull()
}

func (e *evaluator) evaluateReturnStatement(node ast.Node) object.Object {
	var (
		o object.Object
	)
	n := node.(*ast.ReturnStatement)
	if n.Expression != nil {
		o = e.evaluate(n.Expression)
	} else {
		o = object.NewNull()
	}
	return object.NewReturnValue(o)
}

func (e *evaluator) evaluateBlockStatement(stmt ...ast.Statement) (o object.Object) {
	for _, s := range stmt {
		if o = e.evaluate(s); o != nil && o.Type() == object.ReturnObject {
			return
		}
	}
	return
}

func (e *evaluator) evaluateIfStatement(node ast.Node) (o object.Object) {
	ifStmt := node.(*ast.IfStatement)
	cond := e.evaluate(ifStmt.Condition)

	if cond.Type() == object.BooleanObject {

		cond := cond.(*object.Boolean)
		e.symbolTable.EnterScope() // enter nested scope

		if cond.Value {
			o = e.evaluateBlockStatement(ifStmt.IfBlock...)
		} else {
			o = e.evaluateBlockStatement(ifStmt.ElseBlock...)
		}

		e.symbolTable.ExitScope() // exit nested scope

	} else {
		e.quit(internal.NewError(ifStmt.Metadata, internal.ErrConditionType, internal.RuntimeErr))
	}
	return
}

func (e *evaluator) evaluateWhileStatement(node ast.Node) (o object.Object) {
	wStmt := node.(*ast.WhileStatement)
	cond := e.evaluate(wStmt.Condition)

	if cond.Type() == object.BooleanObject {
		state := cond.(*object.Boolean).Value
		for state {

			e.symbolTable.EnterScope() // enter nested scope
			e.evaluateBlockStatement(wStmt.Block...)
			e.symbolTable.ExitScope() // exit nested scope

			state = e.evaluate(wStmt.Condition).(*object.Boolean).Value
		}

	} else {
		e.quit(internal.NewError(wStmt.Metadata, internal.ErrConditionType, internal.RuntimeErr))
	}
	return
}

func (e *evaluator) evaluateFunctionDeclarationStatement(node ast.Node) object.Object {
	fDec := node.(*ast.FunctionDeclarationStatement)
	paramNames := make([]string, len(fDec.Parameters))
	for i, n := range fDec.Parameters {
		paramNames[i] = n.Name
	}
	o := object.NewUserFunction(fDec.Name, paramNames, fDec.Body)
	e.symbolTable.SetUserFunc(o)
	return object.NewNull()
}

func (e *evaluator) quit(err error) {
	fmt.Println(err)

	fmt.Println("\nstack trace ---------- ")
	fCall, ok := e.stackTrace.Pop()

	for ok == true {
		fmt.Println(fmt.Sprintf("FUNCTION CALL %v %v - %v", fCall.FileName(), fCall.LineNumber(),
			fCall.String()))
		fCall, ok = e.stackTrace.Pop()
	}
	os.Exit(0)
}
