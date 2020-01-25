package parser

import (
	"Yum-Programming-Language-Interpreter/ast"
	"Yum-Programming-Language-Interpreter/internal"
	"Yum-Programming-Language-Interpreter/lexer"
	"Yum-Programming-Language-Interpreter/token"
	"bytes"
	"fmt"
	"strconv"
)

type nudMethod func() ast.Expression                          // prefix
type ledMethod func(expression ast.Expression) ast.Expression //infix

type operatorPrecedence int

const (
	MINPRECEDENCE operatorPrecedence = iota
	OR
	AND
	EQUALS      // ==
	CONDITIONAL // >=, >, <=, <
	SUMSUB      // + -
	PRODDIV     // * /
	EXPONENT    // ^ (NOT IMPLEMENTED)
	PREFIX      // -x !x
	POSTFIX     // x++ (NOT IMPLEMENTED)
	CALL        // fn(a,b)

)

var (
	tokenOperPrecedence = map[token.TokenType]operatorPrecedence{
		token.OR:      OR,
		token.AND:     AND,
		token.ADD:     SUMSUB,
		token.SUB:     SUMSUB,
		token.MULT:    PRODDIV,
		token.DIV:     PRODDIV,
		token.EQUAL:   EQUALS,
		token.NEQUAL:  EQUALS,
		token.LTHAN:   CONDITIONAL,
		token.GTHAN:   CONDITIONAL,
		token.LTEQUAL: CONDITIONAL,
		token.GTEQUAL: CONDITIONAL,
	}
)

type prattParserInterface interface {
	parseExpression(precedence operatorPrecedence) (ast.Expression, error)
	parseParameters(bool) []ast.Expression
	parserDataInterface
}

type prattParser struct {
	nudMethods map[token.TokenType]nudMethod
	ledMethods map[token.TokenType]ledMethod
	parserDataInterface
}

func newPrattParser(l lexer.Lexer) (prattParserInterface, error) {
	var (
		nMs = make(map[token.TokenType]nudMethod)
		lMs = make(map[token.TokenType]ledMethod)
		pd  parserDataInterface
		err error
	)

	if pd, err = newParserData(l); err != nil {
		return nil, err
	}

	pp := &prattParser{
		nudMethods:          nMs,
		ledMethods:          lMs,
		parserDataInterface: pd,
	}

	// initialise nud methods
	nMs[token.ADD] = pp.parsePrefixOperator
	nMs[token.SUB] = pp.parsePrefixOperator
	nMs[token.NEGATE] = pp.parsePrefixOperator
	nMs[token.INT] = pp.parseInteger
	nMs[token.FLOAT] = pp.parseFloatingPointNumber
	nMs[token.IDEN] = pp.parseIdent
	nMs[token.BOOLEAN] = pp.parseBoolean
	nMs[token.QUOTATION_MARK] = pp.parseString
	nMs[token.LPAREN] = pp.parseGroupExpression
	nMs[token.LBRACKET] = pp.parseArrayDeclaration

	// initialise led methods
	lMs[token.ADD] = pp.parseInfixOperator
	lMs[token.SUB] = pp.parseInfixOperator
	lMs[token.MULT] = pp.parseInfixOperator
	lMs[token.DIV] = pp.parseInfixOperator
	lMs[token.GTHAN] = pp.parseInfixOperator
	lMs[token.GTEQUAL] = pp.parseInfixOperator
	lMs[token.LTHAN] = pp.parseInfixOperator
	lMs[token.LTEQUAL] = pp.parseInfixOperator
	lMs[token.EQUAL] = pp.parseInfixOperator
	lMs[token.NEQUAL] = pp.parseInfixOperator
	lMs[token.AND] = pp.parseInfixOperator
	lMs[token.OR] = pp.parseInfixOperator

	return pp, err
}

// if error occurs, parser immediately returns the error
func (pp *prattParser) parseExpression(precedence operatorPrecedence) (leftExpr ast.Expression, err error) {
	prefixParseMethod, ok := pp.nudMethods[pp.currentToken().Type()]
	if !ok {
		errMsg := fmt.Sprintf(internal.InvalidPrefixOperatorErr, pp.currentToken().Literal())
		err = internal.NewError(pp.currentToken().Data(), errMsg, internal.SyntaxErr)
		return
	}

	leftExpr = prefixParseMethod()

	for !(pp.currentToken().Type() == token.SEMICOLON) && precedence < pp.currentPrecedence() {
		ledMethod, ok := pp.ledMethods[pp.currentToken().Type()]
		if !ok {
			errMsg := fmt.Sprintf(internal.InvalidInfixOperatorErr, pp.currentToken().Literal())
			err = internal.NewError(pp.currentToken().Data(), errMsg, internal.SyntaxErr)
			return
		}

		leftExpr = ledMethod(leftExpr)
	}
	return
}

func (pp *prattParser) currentPrecedence() operatorPrecedence {
	if p, ok := tokenOperPrecedence[pp.currentToken().Type()]; ok {
		return p
	}
	return MINPRECEDENCE
}

func (pp *prattParser) parseArrayDeclaration() ast.Expression {
	md := pp.currentToken().Data()

	expressionArray := pp.parseParameters(true)
	if expressionArray == nil {
		pp.consumeStatement()
		return nil
	}

	return ast.NewArray(md, expressionArray)
}

func (pp *prattParser) parsePrefixOperator() (expr ast.Expression) {
	prefixOperatorToken := pp.currentToken()
	pp.consume(1)
	if rightExpr, err := pp.parseExpression(PREFIX); err != nil {
		pp.recordError(err)
		pp.consumeStatement()
		return
	} else if rightExpr != nil {
		expr = ast.NewPrefixExpression(prefixOperatorToken, rightExpr)
	}

	return
}

func (pp *prattParser) parseInteger() (expr ast.Expression) {
	var (
		i   int64
		err error
	)

	// convert string literal to int
	if i, err = strconv.ParseInt(pp.currentToken().Literal(), 10, 64); err != nil {
		errMsg := fmt.Sprintf(internal.TypeErr, pp.currentToken().Literal(), token.INT)
		pp.recordError(internal.NewError(pp.currentToken().Data(), errMsg, internal.SyntaxErr))
		pp.consumeStatement() // move to next statement and continue
		return
	}

	expr = ast.NewIntegerExpression(pp.currentToken(), i)
	pp.consume(1) // consume int

	return
}

func (pp *prattParser) parseFloatingPointNumber() (expr ast.Expression) {
	var (
		i   float64
		err error
	)

	// convert string literal to int
	if i, err = strconv.ParseFloat(pp.currentToken().Literal(), 10); err != nil {
		errMsg := fmt.Sprintf(internal.TypeErr, pp.currentToken().Literal(), token.FLOAT)
		pp.recordError(internal.NewError(pp.currentToken().Data(), errMsg, internal.SyntaxErr))
		pp.consumeStatement() // move to next statement and continue
		return
	}

	expr = ast.NewFloatingPointExpression(pp.currentToken(), i)
	pp.consume(1) // consume float

	return
}

func (pp *prattParser) parseString() (expr ast.Expression) {
	md := pp.currentToken().Data()
	pp.consume(1) // consume left quotation mark

	sBuff := bytes.Buffer{}
	for pp.currentToken().Type() != token.QUOTATION_MARK {
		if pp.currentToken().Type() == token.EOF {
			pp.recordError(internal.NewError(pp.currentToken().Data(), fmt.Sprintf(internal.EndOfFileErr, pp.currentToken().LineNumber()), internal.SyntaxErr))
			pp.consumeStatement()
			return
		}

		sBuff.WriteString(pp.currentToken().Literal())
		pp.consume(1)
	}
	pp.consume(1) // consume right quotation mark
	expr = ast.NewStringExpression(md, sBuff.String())
	return
}

// useBrackets false for (), true for []
func (pp *prattParser) parseParameters(useBrackets bool) (parameters []ast.Expression) {
	parameters = make([]ast.Expression, 0)

	sToken := token.LPAREN
	eToken := token.RPAREN
	if useBrackets {
		sToken = token.LBRACKET
		eToken = token.RBRACKET
	}

	if pp.currentToken().Type() != sToken {
		errMsg := fmt.Sprintf(internal.InvalidTokenErr, sToken, pp.currentToken().Literal())
		pp.recordError(internal.NewError(pp.currentToken().Data(), errMsg, internal.SyntaxErr))
		parameters = nil
		return
	}
	pp.consume(1)

	for pp.currentToken().Type() != eToken && pp.currentToken().Type() != token.EOF {

		if expr, err := pp.parseExpression(MINPRECEDENCE); err != nil {
			pp.recordError(err)
			pp.consumeStatement()
			return
		} else {
			parameters = append(parameters, expr)
		}

		if pp.currentToken().Type() != eToken {

			if pp.currentToken().Type() != token.COMMA {
				errMsg := fmt.Sprintf(internal.InvalidTokenErr, ",", pp.currentToken().Literal())
				pp.recordError(internal.NewError(pp.currentToken().Data(), errMsg, internal.SyntaxErr))
				parameters = nil
				return

			}
			pp.consume(1) // consume comma
		}
	}

	if pp.currentToken().Type() != eToken {
		errMsg := fmt.Sprintf(internal.InvalidTokenErr, eToken, pp.currentToken().Literal())
		pp.recordError(internal.NewError(pp.currentToken().Data(), errMsg, internal.SyntaxErr))
		parameters = nil
	}

	pp.consume(1) // consume right paren
	return
}

func (pp *prattParser) parseIdent() (expr ast.Expression) {

	// function call
	idenToken := pp.currentToken()
	if pp.peekToken().Type() == token.LPAREN {
		pp.consume(1)
		if params := pp.parseParameters(false); params != nil {
			expr = ast.NewFunctionCallExpression(idenToken.Data(), idenToken.Literal(), params...)
		} else {
			expr = ast.NewFunctionCallExpression(idenToken.Data(), idenToken.Literal(), nil)
			pp.consumeStatement()
			return
		}

	} else if pp.peekToken().Type() == token.LBRACKET {
		// array index

		pp.consume(2) // consume iden and left bracket
		//indexExpr := pp.parseExpression(MINPRECEDENCE)
		var (
			indexExpr ast.Expression
			err       error
		)
		if indexExpr, err = pp.parseExpression(MINPRECEDENCE); err != nil {
			pp.recordError(err)
			pp.consumeStatement()
			return
		}

		if pp.currentToken().Type() != token.RBRACKET {
			errMsg := fmt.Sprintf(internal.InvalidTokenErr, token.RBRACKET, pp.currentToken().Type())
			pp.recordError(internal.NewError(pp.currentToken().Data(), errMsg, internal.SyntaxErr))
			pp.consumeStatement()
			return
		}
		expr = ast.NewArrayIndexExpression(idenToken.Data(), idenToken.Literal(), indexExpr)
		pp.consume(1) //  right bracket

	} else {
		expr = ast.NewIdentifierExpression(idenToken)
		pp.consume(1)
	}
	return
}

func (pp *prattParser) parseBoolean() (expr ast.Expression) {
	expr = ast.NewBooleanExpression(pp.currentToken(), pp.currentToken().Literal() == "true")
	pp.consume(1)
	return
}

func (pp *prattParser) parseGroupExpression() (expr ast.Expression) {
	var err error
	pp.consume(1) // consume left paren

	if expr, err = pp.parseExpression(MINPRECEDENCE); err != nil {
		pp.recordError(err)
		pp.consumeStatement()
		return
	}

	if pp.currentToken().Type() != token.RPAREN {
		errMsg := fmt.Sprintf(internal.InvalidTokenErr, token.RPAREN, pp.currentToken().Literal())
		pp.recordError(internal.NewError(pp.currentToken().Data(), errMsg, internal.SyntaxErr))
		pp.consumeStatement()
		return nil
	}
	pp.consume(1)
	return
}

func (pp *prattParser) parseInfixOperator(leftExpr ast.Expression) (expr ast.Expression) {
	t := pp.currentToken()
	pp.consume(1)
	if rightExpr, err := pp.parseExpression(tokenOperPrecedence[t.Type()]); err != nil {
		pp.recordError(err)
		pp.consumeStatement()
		return
	} else {
		expr = ast.NewInfixExpression(t, leftExpr, rightExpr)
	}
	return
}
