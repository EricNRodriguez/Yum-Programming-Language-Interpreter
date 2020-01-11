package parser

import (
	"Yum-Programming-Language-Interpreter/ast"
	"Yum-Programming-Language-Interpreter/internal"
	"Yum-Programming-Language-Interpreter/lexer"
	"Yum-Programming-Language-Interpreter/token"
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
	parseExpression(precedence operatorPrecedence) ast.Expression
	parseParameters() []ast.Expression
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
	nMs[token.IDEN] = pp.parseIdent
	nMs[token.BOOLEAN] = pp.parseBoolean
	nMs[token.LPAREN] = pp.parseGroupExpression

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

func (pp *prattParser) parseExpression(precedence operatorPrecedence) ( leftExpr ast.Expression) {
	prefixParseMethod, ok := pp.nudMethods[pp.currentToken().Type()]
	if !ok {
		errMsg := fmt.Sprintf(internal.ERR_INVALID_PREFIX_OPERATOR, pp.currentToken().Literal())
		pp.recordError(internal.NewError(pp.currentToken().Data(), errMsg, internal.SyntaxErr))
		//pp.consumeCurrentStatement()
		pp.progressToNextSemicolon()

		return nil
	}

	leftExpr = prefixParseMethod()

	for !(pp.currentToken().Type() == token.SEMICOLON) && precedence < pp.currentPrecedence() {
		ledMethod, ok := pp.ledMethods[pp.currentToken().Type()]
		if !ok {
			errMsg := fmt.Sprintf(internal.ERR_INVALID_INFIX_OPERATOR, pp.currentToken().Literal())
			pp.recordError(internal.NewError(pp.currentToken().Data(), errMsg, internal.SyntaxErr))
			//pp.consumeCurrentStatement()
			pp.progressToNextSemicolon()
			return nil
		}

		leftExpr = ledMethod(leftExpr)
	}
	return leftExpr
}


func (pp *prattParser) currentPrecedence() operatorPrecedence {
	if p, ok := tokenOperPrecedence[pp.currentToken().Type()]; ok {
		return p
	}
	return MINPRECEDENCE
}

func (pp *prattParser) parsePrefixOperator() (expr ast.Expression) {
	prefixOperatorToken := pp.currentToken()
	pp.consume(1)
	if rightExpr := pp.parseExpression(PREFIX); rightExpr != nil {
		expr = ast.NewPrefixExpression(prefixOperatorToken, rightExpr)
	}
	return
}

func (pp *prattParser) parseInteger() (expr ast.Expression) {
	var (
		i   int
		err error
	)

	// convert string literal to int
	if i, err = strconv.Atoi(pp.currentToken().Literal()); err != nil {
		errMsg := fmt.Sprintf(internal.ERR_INVALID_INFIX_OPERATOR, pp.currentToken().Literal())
		pp.recordError(internal.NewError(pp.currentToken().Data(), errMsg, internal.SyntaxErr))
		pp.progressToNextSemicolon() // move to next statement and continue
		return
	}

	expr = ast.NewIntegerExpression(pp.currentToken(), i)
	pp.consume(1) // consume int

	return
}

func (pp *prattParser) parseParameters() (parameters []ast.Expression) {
	parameters = make([]ast.Expression, 0)

	if pp.currentToken().Type() != token.LPAREN {
		errMsg := fmt.Sprintf(internal.ERR_INVALID_TOKEN, token.LPAREN, pp.currentToken().Literal())
		pp.recordError(internal.NewError(pp.currentToken().Data(), errMsg, internal.SyntaxErr))
		parameters = nil
		return
	}
	pp.consume(1)

	for pp.currentToken().Type() != token.RPAREN && pp.currentToken().Type() != token.EOF {
		parameters = append(parameters, pp.parseExpression(MINPRECEDENCE))

		if pp.currentToken().Type() != token.RPAREN {

			if pp.currentToken().Type() != token.COMMA {
				errMsg := fmt.Sprintf(internal.ERR_INVALID_TOKEN, ",", pp.currentToken().Literal())
				pp.recordError(internal.NewError(pp.currentToken().Data(), errMsg, internal.SyntaxErr))
				parameters = nil
				return

			}
			pp.consume(1) // consume comma
		}
	}

	pp.consume(1) // consume right paren
	return
}

func (pp *prattParser) parseIdent() (expr ast.Expression) {

	// function call
	if pp.peekToken().Type() == token.LPAREN {
		idenToken := pp.currentToken()
		pp.consume(1)
		if params := pp.parseParameters(); params != nil {
			expr = ast.NewFunctionCallExpression(idenToken.Data(), idenToken.Literal(), params...)
		} else {
			pp.progressToNextSemicolon()
			return
		}

	} else {
		expr = ast.NewIdentifierExpression(pp.currentToken())
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
	pp.consume(1)
	if expr = pp.parseExpression(MINPRECEDENCE); expr == nil {
		pp.progressToNextSemicolon()
		return nil
	}

	if pp.currentToken().Type() != token.RPAREN {
		errMsg := fmt.Sprintf(internal.ERR_INVALID_TOKEN, token.RPAREN, pp.currentToken().Literal())
		pp.recordError(internal.NewError(pp.currentToken().Data(), errMsg, internal.SyntaxErr))
		//pp.consumeCurrentStatement()
		return nil
	}
	pp.consume(1)
	return
}

func (pp *prattParser) parseInfixOperator(leftExpr ast.Expression) (expr ast.Expression) {
	t := pp.currentToken()
	pp.consume(1)
	rightExpr := pp.parseExpression(tokenOperPrecedence[t.Type()])
	expr = ast.NewInfixExpression(t, leftExpr, rightExpr)
	return
}
