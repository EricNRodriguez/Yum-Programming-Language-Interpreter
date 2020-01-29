package parser

import (
	"github.com/EricNRodriguez/yum/ast"
	"github.com/EricNRodriguez/yum/internal"
	"github.com/EricNRodriguez/yum/lexer"
	"github.com/EricNRodriguez/yum/token"
	"bytes"
	"fmt"
	"strconv"
)

type nudMethod func() (ast.Expression, error)                          // prefix
type ledMethod func(expression ast.Expression) (ast.Expression, error) //infix

type operatorPrecedence int

const (
	MinPrecedence operatorPrecedence = iota
	OrPrecedence
	AndPrecedence
	EqualsPrecedence
	ConditionalPrecedence
	AddSubPrecedence
	MultDivPrecedence
	PrefixPrecedence
)

var (
	tokenOperPrecedence = map[token.TokenType]operatorPrecedence{
		token.OrToken:         OrPrecedence,
		token.AndToken:        AndPrecedence,
		token.AddToken:        AddSubPrecedence,
		token.SubToken:        AddSubPrecedence,
		token.MultToken:       MultDivPrecedence,
		token.DivToken:        MultDivPrecedence,
		token.EqualToken:      EqualsPrecedence,
		token.NotEqualToken:   EqualsPrecedence,
		token.LThanToken:      ConditionalPrecedence,
		token.GThanToken:      ConditionalPrecedence,
		token.LThanEqualToken: ConditionalPrecedence,
		token.GThanEqualToken: ConditionalPrecedence,
	}
)

type PrattParser interface {
	parseExpression(precedence operatorPrecedence) (ast.Expression, error)
	parseParameters(bool) ([]ast.Expression, error)
	ParserData
}

type prattParser struct {
	nudMethods map[token.TokenType]nudMethod
	ledMethods map[token.TokenType]ledMethod
	ParserData
}

func newPrattParser(l lexer.Lexer) (*prattParser, error) {
	var (
		nMs = make(map[token.TokenType]nudMethod)
		lMs = make(map[token.TokenType]ledMethod)
		pd  ParserData
		err error
	)

	if pd, err = newParserData(l); err != nil {
		return nil, err
	}

	pp := &prattParser{
		nudMethods: nMs,
		ledMethods: lMs,
		ParserData: pd,
	}

	// initialise nud methods
	nMs[token.AddToken] = pp.parsePrefixOperator
	nMs[token.SubToken] = pp.parsePrefixOperator
	nMs[token.NegateToken] = pp.parsePrefixOperator
	nMs[token.IntegerToken] = pp.parseInteger
	nMs[token.FloatingPointToken] = pp.parseFloatingPointNumber
	nMs[token.IdentifierToken] = pp.parseIdent
	nMs[token.BooleanToken] = pp.parseBoolean
	nMs[token.QuotationMarkToken] = pp.parseString
	nMs[token.LeftParenToken] = pp.parseGroupExpression
	nMs[token.LeftBracketToken] = pp.parseArrayNodeDeclaration

	// initialise led methods
	lMs[token.AddToken] = pp.parseInfixOperator
	lMs[token.SubToken] = pp.parseInfixOperator
	lMs[token.MultToken] = pp.parseInfixOperator
	lMs[token.DivToken] = pp.parseInfixOperator
	lMs[token.GThanToken] = pp.parseInfixOperator
	lMs[token.GThanEqualToken] = pp.parseInfixOperator
	lMs[token.LThanToken] = pp.parseInfixOperator
	lMs[token.LThanEqualToken] = pp.parseInfixOperator
	lMs[token.EqualToken] = pp.parseInfixOperator
	lMs[token.NotEqualToken] = pp.parseInfixOperator
	lMs[token.AndToken] = pp.parseInfixOperator
	lMs[token.OrToken] = pp.parseInfixOperator

	return pp, err
}

// if error occurs, parser immediately returns the error
func (pp *prattParser) parseExpression(precedence operatorPrecedence) (leftExpr ast.Expression, err error) {
	prefixParseMethod, ok := pp.nudMethods[pp.currentToken().Type()]
	if !ok {
		errMsg := fmt.Sprintf(internal.ErrInvalidPrefixOperator, pp.currentToken().Literal())
		err = internal.NewError(pp.currentToken().Data(), errMsg, internal.SyntaxErr)
		return
	}

	if leftExpr, err = prefixParseMethod(); err != nil {
		return
	}

	for !(pp.currentToken().Type() == token.SemicolonToken) && precedence < pp.currentPrecedence() {
		ledMethod, ok := pp.ledMethods[pp.currentToken().Type()]
		if !ok {
			errMsg := fmt.Sprintf(internal.ErrInvalidInfixOperator, pp.currentToken().Literal())
			err = internal.NewError(pp.currentToken().Data(), errMsg, internal.SyntaxErr)
			return
		}

		if leftExpr, err = ledMethod(leftExpr); err != nil {
			return
		}
	}
	return
}

func (pp *prattParser) currentPrecedence() operatorPrecedence {
	if p, ok := tokenOperPrecedence[pp.currentToken().Type()]; ok {
		return p
	}
	return MinPrecedence
}

func (pp *prattParser) parseArrayNodeDeclaration() (expr ast.Expression, err error) {
	var (
		arrayExprs []ast.Expression
		md         = pp.currentToken().Data()
	)

	if arrayExprs, err = pp.parseParameters(true); err != nil {
		pp.recordError(err)
		pp.consumeStatement()
		return
	}

	expr = ast.NewArray(md, arrayExprs)
	return
}

func (pp *prattParser) parsePrefixOperator() (expr ast.Expression, err error) {
	prefixOperatorToken := pp.currentToken()
	pp.consume(1)

	var rightExpr ast.Expression
	if rightExpr, err = pp.parseExpression(PrefixPrecedence); err != nil {
		return

	} else if rightExpr != nil {
		expr = ast.NewPrefixExpression(prefixOperatorToken, rightExpr)

	}
	return
}

func (pp *prattParser) parseInteger() (expr ast.Expression, err error) {
	var (
		i int64
	)

	// convert string literal to int
	if i, err = strconv.ParseInt(pp.currentToken().Literal(), 10, 64); err != nil {
		errMsg := fmt.Sprintf(internal.ErrType, pp.currentToken().Literal(), token.IntegerToken)
		err = internal.NewError(pp.currentToken().Data(), errMsg, internal.SyntaxErr)
		return
	}

	expr = ast.NewIntegerExpression(pp.currentToken(), i)
	pp.consume(1) // consume int

	return
}

func (pp *prattParser) parseFloatingPointNumber() (expr ast.Expression, err error) {
	var (
		i float64
	)

	// convert string literal to int
	if i, err = strconv.ParseFloat(pp.currentToken().Literal(), 10); err != nil {
		errMsg := fmt.Sprintf(internal.ErrType, pp.currentToken().Literal(), token.FloatingPointToken)
		err = internal.NewError(pp.currentToken().Data(), errMsg, internal.SyntaxErr)
		return
	}

	expr = ast.NewFloatingPointExpression(pp.currentToken(), i)
	pp.consume(1) // consume float

	return
}

func (pp *prattParser) parseString() (expr ast.Expression, err error) {
	md := pp.currentToken().Data()
	pp.consume(1) // consume left quotation mark

	sBuff := bytes.Buffer{}
	for pp.currentToken().Type() != token.QuotationMarkToken {
		if pp.currentToken().Type() == token.EOFToken {
			err = internal.NewError(pp.currentToken().Data(), fmt.Sprintf(internal.ErrEndOfFile,
				pp.currentToken().LineNumber()), internal.SyntaxErr)
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
func (pp *prattParser) parseParameters(useBrackets bool) (parameters []ast.Expression, err error) {
	parameters = make([]ast.Expression, 0)

	sToken := token.LeftParenToken
	eToken := token.RightParenToken
	if useBrackets {
		sToken = token.LeftBracketToken
		eToken = token.RightBracketToken
	}

	if pp.currentToken().Type() != sToken {
		errMsg := fmt.Sprintf(internal.ErrInvalidToken, sToken, pp.currentToken().Literal())
		pp.recordError(internal.NewError(pp.currentToken().Data(), errMsg, internal.SyntaxErr))
		parameters = nil
		return
	}
	pp.consume(1)

	for pp.currentToken().Type() != eToken && pp.currentToken().Type() != token.EOFToken {
		var (
			expr ast.Expression
		)

		if expr, err = pp.parseExpression(MinPrecedence); err != nil {
			return
		}
		parameters = append(parameters, expr)

		if pp.currentToken().Type() != eToken {

			if pp.currentToken().Type() != token.CommaToken {
				errMsg := fmt.Sprintf(internal.ErrInvalidToken, eToken, pp.currentToken().Literal())
				err = internal.NewError(pp.currentToken().Data(), errMsg, internal.SyntaxErr)
				return

			}

			if pp.peekToken().Type() == eToken && pp.currentToken().Type() == token.CommaToken {
				errMsg := fmt.Sprintf(internal.ErrInvalidToken, eToken, pp.currentToken().Literal())
				err = internal.NewError(pp.currentToken().Data(), errMsg, internal.SyntaxErr)
				return
			}

			pp.consume(1) // consume comma
		}

	}

	if pp.currentToken().Type() != eToken {
		errMsg := fmt.Sprintf(internal.ErrInvalidToken, eToken, pp.currentToken().Literal())
		pp.recordError(internal.NewError(pp.currentToken().Data(), errMsg, internal.SyntaxErr))
		parameters = nil
	}

	pp.consume(1) // consume right paren
	return
}

func (pp *prattParser) parseIdent() (expr ast.Expression, err error) {

	// function call
	idenToken := pp.currentToken()
	if pp.peekToken().Type() == token.LeftParenToken {
		pp.consume(1)

		var params []ast.Expression
		if params, err = pp.parseParameters(false); err != nil {
			return
		}
		expr = ast.NewFunctionCallExpression(idenToken.Data(), idenToken.Literal(), params...)

	} else if pp.peekToken().Type() == token.LeftBracketToken {
		// ArrayNode index
		pp.consume(2) // consume iden and left bracket

		var indexExpr ast.Expression
		if indexExpr, err = pp.parseExpression(MinPrecedence); err != nil {
			return
		}

		if pp.currentToken().Type() != token.RightBracketToken {
			errMsg := fmt.Sprintf(internal.ErrInvalidToken, token.RightBracketToken, pp.currentToken().Type())
			err = internal.NewError(pp.currentToken().Data(), errMsg, internal.SyntaxErr)
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

func (pp *prattParser) parseBoolean() (expr ast.Expression, err error) {
	expr = ast.NewBooleanExpression(pp.currentToken(), pp.currentToken().Literal() == "true")
	pp.consume(1)
	return
}

func (pp *prattParser) parseGroupExpression() (expr ast.Expression, err error) {
	pp.consume(1) // consume left paren

	if expr, err = pp.parseExpression(MinPrecedence); err != nil {
		return
	}

	if pp.currentToken().Type() != token.RightParenToken {
		errMsg := fmt.Sprintf(internal.ErrInvalidToken, token.RightParenToken, pp.currentToken().Literal())
		err = internal.NewError(pp.currentToken().Data(), errMsg, internal.SyntaxErr)
		return
	}
	pp.consume(1)
	return
}

func (pp *prattParser) parseInfixOperator(leftExpr ast.Expression) (expr ast.Expression, err error) {
	var (
		rightExpr ast.Expression
		t         = pp.currentToken()
	)

	pp.consume(1)

	if rightExpr, err = pp.parseExpression(tokenOperPrecedence[t.Type()]); err != nil {
		return
	} else {
		expr = ast.NewInfixExpression(t, leftExpr, rightExpr)
	}
	return
}
