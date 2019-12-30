package parser

import (
	"Yum-Programming-Language-Interpreter/ast"
	"Yum-Programming-Language-Interpreter/lexer"
	"Yum-Programming-Language-Interpreter/token"
	"fmt"
	"errors"
	"strconv"
)

type nudMethod func() ast.ExpressionInterface              // prefix
type ledMethod func(ast.ExpressionInterface) ast.ExpressionInterface //infix


type operatorPrecedence int

const (
	MINPRECEDENCE operatorPrecedence = iota
	EQUALS                           // ==
	CONDITIONAL                      // >=, >, <=, <
	SUMSUB                           // + -
	PRODDIV                          // * /
	EXPONENT                         // ^ (NOT IMPLEMENTED)
	PREFIX                           // -x !x
	POSTFIX                          // x++ (NOT IMPLEMENTED)
	CALL                             // fn(a,b)

)

var (
	tokenOperPrecedence = map[token.TokenType]operatorPrecedence{
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
	parseExpression(precedence operatorPrecedence) ast.ExpressionInterface
	parserDataInterface
}

type prattParser struct {
	nudMethods map[token.TokenType]nudMethod
	ledMethods map[token.TokenType]ledMethod
	parserDataInterface
}

func newPrattParser(l lexer.LexerInterface) (pp *prattParser, err error) {
	var (
		nMs = make(map[token.TokenType]nudMethod)
		lMs = make(map[token.TokenType]ledMethod)
		pd parserDataInterface
	)

	if pd, err = newParserData(l); err != nil {
		return
	}

	pp =  &prattParser{
		nudMethods: nMs,
		ledMethods: lMs,
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

	return
}

func (pp *prattParser) parseExpression(precedence operatorPrecedence) (leftExpr ast.ExpressionInterface) {
	prefixParseMethod, ok := pp.nudMethods[pp.currentToken().Type()]

	if !ok {
		err := errors.New(fmt.Sprintf("unable to parse %v | prefix parse function undefined for token type %v",
			pp.currentToken().Literal(), pp.currentToken().Type()))
		pp.recordError(err)
		// iterate over statement and continue
		pp.progressToNextSemicolon()
		return
	}

	leftExpr = prefixParseMethod()

	for !(pp.currentToken().Type() == token.SEMICOLON) && precedence < pp.currentPrecedence() {

		ledMethod, ok := pp.ledMethods[pp.currentToken().Type()]
		if !ok {
			err := errors.New(fmt.Sprintf("no led parse function available for %v", pp.currentToken().Type()))
			pp.recordError(err)
			pp.progressToNextSemicolon()
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


func (pp *prattParser) parsePrefixOperator() (expr ast.ExpressionInterface){
	prefixOperatorToken := pp.currentToken()
	pp.consume(1)
	rightExpr := pp.parseExpression(PREFIX)
	expr = ast.NewPrefixExpression(prefixOperatorToken, rightExpr)
	return
}


func (pp *prattParser) parseInteger() (expr ast.ExpressionInterface){
	var (
		i   int
		err error
	)

	// convert string literal to int
	if i, err = strconv.Atoi(pp.currentToken().Literal()); err != nil {
		err = errors.New(fmt.Sprintf("unable to parse %v as an integer | %v", pp.currentToken().Literal(),
			err.Error()))
		pp.recordError(err)
		pp.progressToNextSemicolon() // move to next statement and continue
		return
	}

	expr = ast.NewIntegerExpression(pp.currentToken(), i)
	pp.consume(1) // consume int

	return
}


func (pp *prattParser) parseIdent() (expr ast.ExpressionInterface){
	expr = ast.NewTokenExpression(pp.currentToken())
	pp.consume(1)
	return
}


func (pp *prattParser) parseBoolean() (expr ast.ExpressionInterface){
	expr = ast.NewBooleanExpression(pp.currentToken(), pp.currentToken().Literal() == "true")
	pp.consume(1)
	return
}


func (pp *prattParser) parseGroupExpression() (expr ast.ExpressionInterface){
	pp.consume(1)
	if expr = pp.parseExpression(MINPRECEDENCE); expr == nil {
		pp.progressToNextSemicolon()
		return
	}

	if pp.currentToken().Type() != token.RPAREN {
		err := errors.New(fmt.Sprintf("invalid expression, expected RPAREN , recieved %v",
			pp.currentToken().Type()))
		pp.recordError(err)
	}
	pp.consume(1)
	return
}

func (pp *prattParser) parseInfixOperator(leftExpr ast.ExpressionInterface) (expr ast.ExpressionInterface) {
	t := pp.currentToken()
	pp.consume(1)
	rightExpr := pp.parseExpression(tokenOperPrecedence[t.Type()])
	expr = ast.NewInfixExpression(t, leftExpr, rightExpr)
	return
}
