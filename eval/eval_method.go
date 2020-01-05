package eval

import (
	"Yum-Programming-Language-Interpreter/ast"
	"Yum-Programming-Language-Interpreter/object"
	"Yum-Programming-Language-Interpreter/token"
)

var evalMethodRouter map[ast.NodeType]evalMethod

func init() {
	evalMethodRouter = map[ast.NodeType]evalMethod{
		ast.PROGRAM:                        evaluateProgram,
		ast.IDENTIFIER:                     evaluateIdentifier,
		ast.PREFIX_EXPRESSION:              evaluatePrefixExpression,
		ast.INFIX_EXPRESSION:               evaluateInfixExpression,
		ast.INTEGER_EXPRESSION:             evaluateIntegerExpression,
		ast.BOOLEAN_EXPRESSION:             evaluateBooleanExpression,
		ast.FUNC_CALL_EXPRESSION:           evaluateFunctionCallExpression,
		ast.IDENTIFIER_EXPRESSION:          evaluateIdentifierExpression,
		ast.VAR_STATEMENT:                  evaluateVarStatement,
		ast.RETURN_STATEMENT:               evaluateReturnStatement,
		ast.EXPRESSION_STATEMENT:           evaluateExpressionStatement,
		ast.IF_STATEMENT:                   evaluateIfStatement,
		ast.FUNCTION_DECLARATION_STATEMENT: evaluateFunctionDeclarationStatement,
		ast.ASSIGNMENT_STATEMENT:           evaluateAssignmentStatement,
	}
}

type evalMethod func(node ast.Node) object.Object

func evaluateProgram(node ast.Node) (o object.Object) {
	prog := node.(*ast.Program)
	return evaluateBlockStatement(prog.Statements...)
}

func evaluateIdentifier(node ast.Node) (o object.Object) {
	iden := node.(*ast.Identifier)
	o, _ = sT.GetVar(iden.Name)
	return
}

func evaluateIdentifierExpression(node ast.Node) object.Object {
	idenExpr := node.(*ast.IdentifierExpression)
	return evaluateIdentifier(idenExpr.Node)
}

func evaluatePrefixExpression(node ast.Node) (o object.Object) {
	pExpr := node.(*ast.PrefixExpression)
	rObj := Evaluate(pExpr.Expression)

	if rObj.Type() == object.INTEGER {

		rObj := rObj.(*object.Integer)
		switch pExpr.Token.Type() {
		case token.ADD:
			o = rObj
		case token.SUB:
			o = object.NewInteger(-1 * rObj.Value)
		case token.NEGATE:
			o = object.NewBoolean(rObj.Value == 0)
		default:
			o = object.NewNull()
		}

	} else if rObj.Type() == object.BOOLEAN {

		rObj := rObj.(*object.Boolean)
		switch pExpr.Token.Type() {
		case token.ADD:
			o = object.NewNull()
		case token.SUB:
			o = object.NewNull()
		case token.NEGATE:
			o = object.NewBoolean(!rObj.Value)
		default:
			o = object.NewNull()
		}

	} else {
		// null object
		o = object.NewNull()
	}
	return
}

func evaluateInfixExpression(node ast.Node) (o object.Object) {
	iExpr := node.(*ast.InfixExpression)
	lObj := unpack(Evaluate(iExpr.LeftExpression))
	rObj := unpack(Evaluate(iExpr.RightExpression))

	if lObj.Type() == rObj.Type(){
		if lObj.Type() == object.INTEGER {
			lObj := lObj.(*object.Integer)
			rObj := rObj.(*object.Integer)

			switch iExpr.Token.Type() {
			case token.ADD:
				o = object.NewInteger(lObj.Value + rObj.Value)
			case token.SUB:
				o = object.NewInteger(lObj.Value - rObj.Value)
			case token.DIV:
				// check for zero division !! -----------------------------------------------
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
			case token.AND:
				o = object.NewBoolean(lObj.Value != 0 && rObj.Value != 0)
			case token.OR:
				o = object.NewBoolean(lObj.Value != 0 || rObj.Value != 0)
			default:
				// raise an error here!
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
				// raise an error here! -----------------------------------------------------
				o = object.NewNull()
			}
		}
	} else {
		o = object.NewNull()
	}

	return
}

func unpack(o object.Object) object.Object {
	if o.Type() != object.RETURN {
		return o
	}
	oV := o.(*object.ReturnValue).Value
	return unpack(oV)
}

func evaluateIntegerExpression(node ast.Node) object.Object {
	i := node.(*ast.IntegerExpression)
	o := object.NewInteger(i.Value)
	return o
}

func evaluateBooleanExpression(node ast.Node) object.Object {
	b := node.(*ast.BooleanExpression)
	o := object.NewBoolean(b.Value)
	return o
}

func evaluateFunctionCallExpression(node ast.Node) (o object.Object) {
	fCall := node.(*ast.FunctionCallExpression)
	if f, ok := sT.GetUserFunc(fCall.FunctionName); !ok {
		// native function call

		// evaluate parameters
		evalParams := make([]object.Object, len(fCall.Parameters))
		for i, e := range fCall.Parameters {
			evalParams[i] = Evaluate(e)
		}

		f, _ := sT.GetNativeFunc(fCall.FunctionName)
		o = f.Function(evalParams...)

	} else {
		// user defined function

		// evaluate parameters
		paramValues := map[string]object.Object{}
		for i := range fCall.Parameters {
			paramValues[f.Parameters[i]] = Evaluate(fCall.Parameters[i])
		}

		// execute function call
		sT.EnterScope()
		defer sT.ExitScope()

		for k, v := range paramValues {
			sT.SetVar(k, v)
		}

		o = evaluateBlockStatement(f.Body...)
		if o.Type() != object.RETURN {
			o = object.NewNull()
		}

	}
	return o
}

func evaluateVarStatement(node ast.Node) object.Object {
	vStmt := node.(*ast.VarStatement)
	leftObj := Evaluate(vStmt.Expression)
	sT.SetVar(vStmt.Identifier.String(), leftObj)
	return object.NewNull()
}

func evaluateAssignmentStatement(node ast.Node) object.Object {
	vStmt := node.(*ast.AssignmentStatement)
	leftObj := Evaluate(vStmt.Expression)
	sT.UpdateVar(vStmt.Identifier.String(), leftObj)
	return object.NewNull()
}

func evaluateReturnStatement(node ast.Node) object.Object {
	n := node.(*ast.ReturnStatement)
	o := Evaluate(n.Expression)
	return object.NewReturnValue(o)
}

func evaluateExpressionStatement(node ast.Node) object.Object {
	expr := node.(*ast.ExpressionStatement)
	o := Evaluate(expr.Expression)
	return o
}

func evaluateBlockStatement(stmt ...ast.Statement) (o object.Object) {
	for _, s := range stmt {
		if o = Evaluate(s); o.Type() == object.RETURN {
			return
		}
	}
	return
}

func evaluateIfStatement(node ast.Node) (o object.Object) {
	ifStmt := node.(*ast.IfStatement)
	cond := Evaluate(ifStmt.Condition)
	if cond.Type() == object.BOOLEAN {
		cond := cond.(*object.Boolean)
		if cond.Value {
			sT.EnterScope()
			o = evaluateBlockStatement(ifStmt.IfBlock...)
			sT.ExitScope()
		} else if ifStmt.ElseBlock != nil {
			sT.EnterScope()
			o = evaluateBlockStatement(ifStmt.ElseBlock...)
			sT.ExitScope()
		} else {
			// potentially implement a quick method to just skip to the next statement --------------------------------
			o = object.NewNull()
		}
	} else {
		// record an error
		o = object.NewNull()
	}
	return
}

func evaluateFunctionDeclarationStatement(node ast.Node) object.Object {
	fDec := node.(*ast.FunctionDeclarationStatement)
	paramNames := make([]string, len(fDec.Parameters))
	for i, n := range fDec.Parameters {
		paramNames[i] = n.Name
	}
	o := object.NewUserFunction(fDec.Name, paramNames, fDec.Body)
	sT.SetUserFunc(*o)
	return object.NewNull()
}
