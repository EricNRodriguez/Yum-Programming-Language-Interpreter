package eval

import (
	"Yum-Programming-Language-Interpreter/ast"
	"Yum-Programming-Language-Interpreter/internal"
	"Yum-Programming-Language-Interpreter/object"
	"Yum-Programming-Language-Interpreter/symbol_table"
	"Yum-Programming-Language-Interpreter/token"
	"fmt"
	"os"
)

type evalMethod func(node ast.Node) object.Object

type Evaluator struct {
	stackTrace   internal.StackTrace
	symbolTable  symbol_table.SymbolTable
	methodRouter map[ast.NodeType]evalMethod
}

func NewEvaluator() (e *Evaluator) {
	e = &Evaluator{
		symbolTable: symbol_table.NewSymbolTable(),
		stackTrace:  internal.NewStackTrace(),
	}

	e.methodRouter = map[ast.NodeType]evalMethod{
		ast.PROGRAM:                        e.evaluateProgram,
		ast.IDENTIFIER:                     e.evaluateIdentifier,
		ast.ARRAY:                          e.evaluateArrayExpression,
		ast.ARRAY_INDEX_EXPRESSION: e.evaluateArrayIndexExpression,
		ast.PREFIX_EXPRESSION:              e.evaluatePrefixExpression,
		ast.INFIX_EXPRESSION:               e.evaluateInfixExpression,
		ast.INTEGER_EXPRESSION:             e.evaluateIntegerExpression,
		ast.FLOATING_POINT_EXPRESSION:      e.evaluateFloatingPointExpression,
		ast.BOOLEAN_EXPRESSION:             e.evaluateBooleanExpression,
		ast.FUNC_CALL_EXPRESSION:           e.evaluateFunctionCallExpression,
		ast.IDENTIFIER_EXPRESSION:          e.evaluateIdentifierExpression,
		ast.VAR_STATEMENT:                  e.evaluateVarStatement,
		ast.RETURN_STATEMENT:               e.evaluateReturnStatement,
		ast.IF_STATEMENT:                   e.evaluateIfStatement,
		ast.FUNCTION_DECLARATION_STATEMENT: e.evaluateFunctionDeclarationStatement,
		ast.FUNCTION_CALL_STATEMENT:        e.evaluateFunctionCallStatement,
		ast.ASSIGNMENT_STATEMENT:           e.evaluateAssignmentStatement,
	}

	return
}

func (e *Evaluator) Evaluate(node ast.Node) (o object.Object) {
	if method, ok := e.methodRouter[node.Type()]; ok {
		return method(node)
	}
	e.panic(internal.NewError(token.NewMetatadata(node.LineNumber(), node.FileName()),
		fmt.Sprintf(internal.UnimplementedType, node.Type()), internal.InternalErr))
	return nil
}

func (e *Evaluator) evaluateProgram(node ast.Node) (o object.Object) {
	prog := node.(*ast.Program)
	return e.evaluateBlockStatement(prog.Statements...)
}

func (e *Evaluator) evaluateIdentifier(node ast.Node) (o object.Object) {
	iden := node.(*ast.Identifier)
	o, _ = e.symbolTable.GetVar(iden.Name)
	return
}

func (e *Evaluator) evaluateIdentifierExpression(node ast.Node) object.Object {
	idenExpr := node.(*ast.IdentifierExpression)
	return e.evaluateIdentifier(idenExpr.Identifier)
}

func (e *Evaluator) evaluatePrefixExpression(node ast.Node) (o object.Object) {
	pExpr := node.(*ast.PrefixExpression)
	rObj := e.Evaluate(pExpr.Expression)

	if rObj.Type() == object.INTEGER {
		rObj := rObj.(*object.Integer)

		switch pExpr.Token.Type() {
		case token.ADD:
			o = rObj
		case token.SUB:
			o = object.NewInteger(-1 * rObj.Value)
		default:
			e.panic(internal.NewError(pExpr.Data(), fmt.Sprintf(internal.TypeErr, rObj.Literal(), object.BOOLEAN),
				internal.RuntimeErr))
		}

	} else if rObj.Type() == object.FLOAT {
		rObj := rObj.(*object.Float)

		switch pExpr.Token.Type() {
		case token.ADD:
			o = rObj
		case token.SUB:
			o = object.NewFloat(-1 * rObj.Value)
		default:
			e.panic(internal.NewError(pExpr.Data(), fmt.Sprintf(internal.TypeErr, rObj.Literal(), object.BOOLEAN),
				internal.RuntimeErr))
		}
	} else if rObj.Type() == object.BOOLEAN {

		rObj := rObj.(*object.Boolean)
		switch pExpr.Token.Type() {
		case token.NEGATE:
			o = object.NewBoolean(!rObj.Value)
		default:
			e.panic(internal.NewError(pExpr.Data(), fmt.Sprintf(internal.TypeErr, rObj.Literal(),
				fmt.Sprintf("%v or %v", object.INTEGER, object.FLOAT)), internal.RuntimeErr))
		}

	} else {
		// null object
		e.panic(internal.NewError(pExpr.Data(), fmt.Sprintf(internal.TypeErr, rObj.Literal(),
			fmt.Sprintf("%v or %v or %v", object.INTEGER, object.FLOAT, object.BOOLEAN)), internal.RuntimeErr))
	}
	return
}

func (e *Evaluator) evaluateInfixExpression(node ast.Node) (o object.Object) {
	iExpr := node.(*ast.InfixExpression)
	lObj := e.unpack(e.Evaluate(iExpr.LeftExpression))
	rObj := e.unpack(e.Evaluate(iExpr.RightExpression))

	if lObj.Type() == object.INTEGER && rObj.Type() == object.INTEGER {
		if lObj.Type() == object.INTEGER {
			lObj := lObj.(*object.Integer)
			rObj := rObj.(*object.Integer)

			switch iExpr.Token.Type() {
			case token.ADD:
				o = object.NewInteger(lObj.Value + rObj.Value)
			case token.SUB:
				o = object.NewInteger(lObj.Value - rObj.Value)
			case token.DIV:
				if rObj.Value == 0 {
					e.panic(internal.NewError(iExpr.Data(), internal.DivisionByZeroErr, internal.RuntimeErr))
				}

				o = object.NewInteger(lObj.Value / rObj.Value)
			case token.MULT:
				o = object.NewInteger(lObj.Value * rObj.Value)
			case token.GTHAN:
				o = object.NewBoolean(lObj.Value > rObj.Value)
			case token.LTHAN:
				o = object.NewBoolean(lObj.Value < rObj.Value)
			case token.GTEQUAL:
				o = object.NewBoolean(lObj.Value >= rObj.Value)
			case token.LTEQUAL:
				o = object.NewBoolean(lObj.Value <= rObj.Value)
			case token.EQUAL:
				o = object.NewBoolean(lObj.Value == rObj.Value)
			case token.NEQUAL:
				o = object.NewBoolean(lObj.Value != rObj.Value)
			default:
				e.panic(internal.NewError(iExpr.Data(), fmt.Sprintf(internal.TypeOperationErr, iExpr.Token.Type(),
					lObj.Type()), internal.RuntimeErr))
				o = object.NewNull()
			}
		} else {
			lObj := lObj.(*object.Boolean)
			rObj := rObj.(*object.Boolean)
			switch iExpr.Token.Type() {
			case token.EQUAL:
				o = object.NewBoolean(lObj.Value == rObj.Value)
			case token.NEQUAL:
				o = object.NewBoolean(lObj.Value != rObj.Value)
			case token.AND:
				o = object.NewBoolean(lObj.Value && rObj.Value)
			case token.OR:
				o = object.NewBoolean(lObj.Value || rObj.Value)
			default:
				e.panic(internal.NewError(iExpr.Data(), fmt.Sprintf(internal.TypeOperationErr, iExpr.Token.Type(),
					lObj.Type()), internal.RuntimeErr))
			}
		}
	} else if (lObj.Type() == object.INTEGER || lObj.Type() == object.FLOAT) &&
		(rObj.Type() == object.INTEGER || rObj.Type() == object.FLOAT) {

		// type cast to floating point numbers
		lObj := e.castToFloat(lObj)
		rObj := e.castToFloat(rObj)

		switch iExpr.Token.Type() {
		case token.ADD:
			o = object.NewFloat(lObj.Value + rObj.Value)
		case token.SUB:
			o = object.NewFloat(lObj.Value - rObj.Value)
		case token.DIV:
			if rObj.Value == 0 {
				e.panic(internal.NewError(iExpr.Data(), internal.DivisionByZeroErr, internal.RuntimeErr))
			}

			o = object.NewFloat(lObj.Value / rObj.Value)
		case token.MULT:
			o = object.NewFloat(lObj.Value * rObj.Value)
		case token.GTHAN:
			o = object.NewBoolean(lObj.Value > rObj.Value)
		case token.LTHAN:
			o = object.NewBoolean(lObj.Value < rObj.Value)
		case token.GTEQUAL:
			o = object.NewBoolean(lObj.Value >= rObj.Value)
		case token.LTEQUAL:
			o = object.NewBoolean(lObj.Value <= rObj.Value)
		case token.EQUAL:
			o = object.NewBoolean(lObj.Value == rObj.Value)
		case token.NEQUAL:
			o = object.NewBoolean(lObj.Value != rObj.Value)
		default:
			e.panic(internal.NewError(iExpr.Data(), fmt.Sprintf(internal.TypeOperationErr, iExpr.Token.Type(),
				lObj.Type()), internal.RuntimeErr))
			o = object.NewNull()
		}

	} else {
		e.panic(internal.NewError(iExpr.Data(), fmt.Sprintf(internal.MismatchedTypeErr, lObj.Type(), rObj.Type()),
			internal.RuntimeErr))
	}

	return
}

// must receive either an integer or a float, will panic otherwise
func (e *Evaluator) castToFloat(obj object.Object) (o *object.Float) {
	var ok bool
	if o, ok = obj.(*object.Float); !ok {
		o = object.NewFloat(float64(obj.(*object.Integer).Value))
	}
	return o
}

func (e *Evaluator) unpack(o object.Object) object.Object {
	if o.Type() != object.RETURN {
		return o
	}
	oV := o.(*object.ReturnValue).Value
	return e.unpack(oV)
}

func (e *Evaluator) evaluateIntegerExpression(node ast.Node) object.Object {
	i := node.(*ast.IntegerExpression)
	o := object.NewInteger(i.Value)
	return o
}

func (e *Evaluator) evaluateArrayExpression(node ast.Node) object.Object {
	a := node.(*ast.ArrayExpression)
	oData := make([]object.Object, a.Length)
	for i, ex := range a.Data {
		oData[i] = e.Evaluate(ex)
	}
	o := object.NewArray(oData)
	return o
}

func (e *Evaluator) evaluateArrayIndexExpression(node ast.Node) (o object.Object) {
	iden := node.(*ast.ArrayIndexExpression)

	arrE, _ := e.symbolTable.GetVar(iden.ArrayName)
	if arrE.Type() != object.ARRAY {
		errMsg := fmt.Sprintf(internal.TypeErr, arrE.Literal(), object.ARRAY)
		e.panic(internal.NewError(iden.Metadata, errMsg, internal.RuntimeErr))
	}
	arr := arrE.(*object.Array)

	indE := e.Evaluate(iden.IndexExpr)
	if indE.Type() != object.INTEGER {
		errMsg := fmt.Sprintf(internal.TypeErr, arrE.Literal(), object.INTEGER)
		e.panic(internal.NewError(iden.Metadata, errMsg, internal.RuntimeErr))
	}

	indI := indE.(*object.Integer)
	if indI.Value > arr.Length - 1 || indI.Value < 0 {
		// index out of bounds error
		e.panic(internal.NewError(iden.Metadata, internal.IndexOutOfBoundsErr, internal.RuntimeErr))
	}

	o = arr.Data[indI.Value]
	return
}

func (e *Evaluator) evaluateFloatingPointExpression(node ast.Node) object.Object {
	i := node.(*ast.FloatingPointExpression)
	o := object.NewFloat(i.Value)
	return o
}

func (e *Evaluator) evaluateBooleanExpression(node ast.Node) object.Object {
	b := node.(*ast.BooleanExpression)
	o := object.NewBoolean(b.Value)
	return o
}

func (e *Evaluator) evaluateFunctionCallExpression(node ast.Node) (o object.Object) {
	fCall := node.(*ast.FunctionCallExpression)
	e.stackTrace.Push(fCall) // record function call

	// native function call
	if f, ok := e.symbolTable.GetUserFunc(fCall.FunctionName); !ok {
		// evaluate parameters
		evalParams := make([]object.Object, len(fCall.Parameters))
		for i, expr := range fCall.Parameters {
			evalParams[i] = e.Evaluate(expr)
		}

		f, _ := e.symbolTable.GetNativeFunc(fCall.FunctionName)
		o = f.Function(evalParams...)

		// user defined function
	} else {

		// evaluate parameters
		paramValues := map[string]object.Object{}
		for i, v := range fCall.Parameters {
			paramValues[f.Parameters[i]] = e.Evaluate(v)
		}

		// new symbol table
		e.symbolTable.EnterFunction()

		for k, v := range paramValues {
			e.symbolTable.SetVar(k, v)
		}

		o = e.evaluateBlockStatement(f.Body...)
		if o == nil || o.Type() != object.RETURN {
			o = object.NewNull()
		}
		e.symbolTable.ExitFunction()

	}

	e.stackTrace.Pop()
	return o
}

func (e *Evaluator) evaluateFunctionCallStatement(node ast.Node) (o object.Object) {
	stmt := node.(*ast.FunctionCallStatement)
	e.evaluateFunctionCallExpression(stmt.FunctionCallExpression)
	return
}

func (e *Evaluator) evaluateVarStatement(node ast.Node) object.Object {
	vStmt := node.(*ast.VarStatement)
	leftObj := e.Evaluate(vStmt.Expression)
	e.symbolTable.SetVar(vStmt.Identifier.String(), leftObj)
	return object.NewNull()
}

func (e *Evaluator) evaluateAssignmentStatement(node ast.Node) object.Object {
	vStmt := node.(*ast.AssignmentStatement)
	leftObj := e.Evaluate(vStmt.Expression)
	e.symbolTable.UpdateVar(vStmt.Identifier.String(), leftObj)
	return object.NewNull()
}

func (e *Evaluator) evaluateReturnStatement(node ast.Node) object.Object {
	n := node.(*ast.ReturnStatement)
	o := e.Evaluate(n.Expression)
	return object.NewReturnValue(o)
}

func (e *Evaluator) evaluateBlockStatement(stmt ...ast.Statement) (o object.Object) {
	for _, s := range stmt {
		if o = e.Evaluate(s); o != nil && o.Type() == object.RETURN {
			return
		}
	}
	return
}

func (e *Evaluator) evaluateIfStatement(node ast.Node) (o object.Object) {
	ifStmt := node.(*ast.IfStatement)
	cond := e.Evaluate(ifStmt.Condition)

	if cond.Type() == object.BOOLEAN {

		cond := cond.(*object.Boolean)
		e.symbolTable.EnterScope() // enter nested scope

		if cond.Value {
			o = e.evaluateBlockStatement(ifStmt.IfBlock...)
		} else {
			o = e.evaluateBlockStatement(ifStmt.ElseBlock...)
		}

		e.symbolTable.ExitScope() // exit nested scope

	} else {
		// record an error
		o = object.NewNull()
	}
	return
}

func (e *Evaluator) evaluateFunctionDeclarationStatement(node ast.Node) object.Object {
	fDec := node.(*ast.FunctionDeclarationStatement)
	paramNames := make([]string, len(fDec.Parameters))
	for i, n := range fDec.Parameters {
		paramNames[i] = n.Name
	}
	o := object.NewUserFunction(fDec.Name, paramNames, fDec.Body)
	e.symbolTable.SetUserFunc(o)
	return object.NewNull()
}

func (e *Evaluator) panic(err error) {
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
