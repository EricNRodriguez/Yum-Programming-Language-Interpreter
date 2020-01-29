package parser

import (
	"github.com/EricNRodriguez/yum/ast"
	"github.com/EricNRodriguez/yum/internal"
	"github.com/EricNRodriguez/yum/lexer"
	"github.com/EricNRodriguez/yum/token"
	"fmt"
)

type parseMethod func() ast.Statement

type RecursiveDescentParser struct {
	parseMethodRouter map[token.TokenType]parseMethod
	PrattParser
}

func NewRecursiveDescentParser(l lexer.Lexer) (Parser, error) {
	var (
		pMR         = make(map[token.TokenType]parseMethod) // parse method router
		prattParser PrattParser
		err         error
	)

	if prattParser, err = newPrattParser(l); err != nil {
		return nil, err
	}

	rdp := &RecursiveDescentParser{
		parseMethodRouter: pMR,
		PrattParser:       prattParser,
	}

	// initialise pMR
	pMR[token.VarToken] = rdp.parseVarStatement
	pMR[token.ReturnToken] = rdp.parseReturnStatement
	pMR[token.IdentifierToken] = rdp.parseIdenStatement
	pMR[token.IfToken] = rdp.parseIfStatement
	pMR[token.FuncToken] = rdp.parseFuncDeclarationStatement
	pMR[token.WhileToken] = rdp.parseWhileStatement

	return rdp, err
}

func (rdp *RecursiveDescentParser) Parse() (prog *ast.Program, errs []error) {

	stmts := make([]ast.Statement, 0)
	for rdp.checkNextToken() {

		if stmt, err := rdp.parseStatement(); err != nil {
			rdp.recordError(err)
			rdp.consumeStatement()

		} else if stmt != nil {
			stmts = append(stmts, stmt)
		}

		rdp.consume(1) // moving to next statement
	}

	prog = ast.NewProgram(rdp.currentToken().Data(), stmts...)
	prog.Hoist() // moves declarations to the top of the file

	return prog, rdp.errors()
}

func (rdp *RecursiveDescentParser) parseStatement() (stmt ast.Statement, err error) {
	var (
		pM parseMethod
		ok bool
	)

	if pM, ok = rdp.parseMethodRouter[rdp.currentToken().Type()]; !ok {
		errMsg := fmt.Sprintf(internal.ErrInvalidStatement, rdp.currentToken().Literal())
		err = internal.NewError(rdp.currentToken().Data(), errMsg, internal.SyntaxErr)
		return

	}

	stmt = pM()
	if rdp.currentToken().Type() != token.SemicolonToken {
		errMsg := fmt.Sprintf(internal.ErrInvalidToken, ";", rdp.currentToken().Literal())
		return nil, internal.NewError(token.NewMetatadata(rdp.currentToken().LineNumber()-1,
			rdp.currentToken().FileName()), errMsg, internal.SyntaxErr)
	}

	return
}

func (rdp *RecursiveDescentParser) parseVarStatement() (stmt ast.Statement) {
	var (
		iden     *ast.IdentifierExpression
		expr     ast.Expression
		varToken = rdp.currentToken()
		err      error
	)

	if !rdp.expectTokenType(token.IdentifierToken) {
		rdp.consumeStatement()
		return
	}
	rdp.consume(1) // consume var

	iden = ast.NewIdentifierExpression(rdp.currentToken())

	if !rdp.expectTokenType(token.AssignToken) {
		rdp.consumeStatement()
		return
	}
	rdp.consume(2)

	// skip stmts with syntax errors
	if expr, err = rdp.parseExpression(MinPrecedence); err != nil {
		rdp.recordError(err)
		rdp.consumeStatement()
		return
	}

	stmt = ast.NewVarStatement(varToken, iden, expr)

	return
}

func (rdp *RecursiveDescentParser) parseReturnStatement() (stmt ast.Statement) {
	retToken := rdp.currentToken()
	rdp.consume(1) // consume return

	// return nothing
	if rdp.currentToken().Type() == token.SemicolonToken {
		stmt = ast.NewReturnStatment(retToken, nil)

	} else if expr, err := rdp.parseExpression(MinPrecedence); err != nil {
		rdp.recordError(err)
		rdp.consumeStatement()

	} else {
		stmt = ast.NewReturnStatment(retToken, expr)

	}

	return
}

func (rdp *RecursiveDescentParser) parseIdenStatement() (stmt ast.Statement) {
	switch rdp.peekToken().Type() {

	case token.AssignToken:
		iden := ast.NewIdentifierExpression(rdp.currentToken())

		if !rdp.expectTokenType(token.AssignToken) {
			rdp.consumeStatement()
			return
		}
		rdp.consume(2) // consume IdentifierNode and assign

		var (
			expr ast.Expression
			err  error
		)

		if expr, err = rdp.parseExpression(MinPrecedence); err != nil {
			rdp.recordError(err)
			rdp.consumeStatement()
			return
		}

		stmt = ast.NewAssignmentStatement(iden.Metadata, iden, expr)

	case token.LeftParenToken:
		stmt = rdp.parseFunctionCallStatement()

	default:
		errMsg := fmt.Sprintf(internal.ErrInvalidStatement, rdp.currentToken().Literal())
		err := internal.NewError(rdp.currentToken().Data(), errMsg, internal.SyntaxErr)
		rdp.recordError(err)
		rdp.consumeStatement()
	}

	return
}

func (rdp *RecursiveDescentParser) parseFunctionCallStatement() (stmt ast.Statement) {
	var (
		expr ast.Expression
		err  error
	)

	if expr, err = rdp.parseExpression(MinPrecedence); err != nil {
		rdp.recordError(err)
		rdp.consumeStatement()
		return
	}

	stmt = ast.NewFunctionCallStatement(expr.(*ast.FunctionCallExpression))
	return
}

func (rdp *RecursiveDescentParser) parseIfStatement() (stmt ast.Statement) {
	var (
		t          = rdp.currentToken()
		trueBlock  []ast.Statement
		falseBlock []ast.Statement
		cond       ast.Expression
		err        error
	)

	if !rdp.expectTokenType(token.LeftParenToken) {
		rdp.consumeStatement()
		return
	}
	rdp.consume(1) // consume left paren

	if cond, err = rdp.parseExpression(MinPrecedence); err != nil {
		rdp.recordError(err)
		rdp.consumeStatement()
		return
	}

	if trueBlock, err = rdp.parseBlockStatement(); err != nil {
		rdp.recordError(err)
		rdp.consumeIfStatement()
		return
	}

	// else
	if rdp.currentToken().Type() == token.ElseToken {
		rdp.consume(1) // consume ELSE
		if falseBlock, err = rdp.parseBlockStatement(); err != nil {
			rdp.recordError(err)
			rdp.consumeBlockStatement()
			return
		}
	}

	stmt = ast.NewIfStatement(t, cond, trueBlock, falseBlock)
	return
}

func (rdp *RecursiveDescentParser) parseWhileStatement() (stmt ast.Statement) {
	var (
		md    = rdp.currentToken().Data()
		cond  ast.Expression
		block []ast.Statement
		err   error
	)

	if !rdp.expectTokenType(token.LeftParenToken) {
		rdp.consumeBlockStatement()
		return
	}
	rdp.consume(2) // consume while and left parenthesis

	if cond, err = rdp.parseExpression(MinPrecedence); err != nil {
		rdp.recordError(err)
		rdp.consumeBlockStatement()
		return
	}

	if rdp.currentToken().Type() != token.RightParenToken {
		rdp.recordError(internal.NewError(rdp.currentToken().Data(), fmt.Sprintf(internal.ErrInvalidToken, token.RightParenToken,
			rdp.currentToken().Type()), internal.SyntaxErr))
		rdp.consumeBlockStatement()
		return
	}
	rdp.consume(1) // consume right parenthesis

	if block, err = rdp.parseBlockStatement(); err != nil {
		rdp.recordError(err)
		rdp.consumeBlockStatement()
		return
	}
	stmt = ast.NewWhileStatement(md, cond, block)
	return
}

func (rdp *RecursiveDescentParser) parseBlockStatement() (bStmt []ast.Statement, err error) {
	bStmt = make([]ast.Statement, 0)
	if rdp.currentToken().Type() != token.LeftBraceToken {
		errMsg := fmt.Sprintf(internal.ErrInvalidToken, token.LeftBraceToken, rdp.currentToken().Literal())
		err = internal.NewError(rdp.currentToken().Data(), errMsg, internal.SyntaxErr)
		return

	}
	rdp.consume(1)

	for rdp.currentToken().Type() != token.RightBraceToken && rdp.currentToken().Type() != token.EOFToken {
		var (
			stmt ast.Statement
		)
		if stmt, err = rdp.parseStatement(); err != nil {
			return
		}

		bStmt = append(bStmt, stmt)
		rdp.consume(1)
	}

	if rdp.currentToken().Type() != token.RightBraceToken {
		errMsg := fmt.Sprintf(internal.ErrInvalidToken, token.RightBraceToken, rdp.currentToken().Literal())
		err = internal.NewError(rdp.currentToken().Data(), errMsg, internal.SyntaxErr)
		return

	}

	// consume right brace
	rdp.consume(1)
	return
}

func (rdp *RecursiveDescentParser) parseFuncDeclarationStatement() (stmt ast.Statement) {
	var (
		t      = rdp.currentToken()
		iden   string
		params []ast.IdentifierExpression
		body   []ast.Statement
		pExprs []ast.Expression // function parameter expressions
		err    error
	)

	if !rdp.expectTokenType(token.IdentifierToken) {
		rdp.consumeStatement()
		return
	}

	rdp.consume(1) // consume func token

	iden = rdp.currentToken().Literal()
	rdp.consume(1) // consume function name

	if pExprs, err = rdp.parseParameters(false); err != nil {
		rdp.recordError(err)
		rdp.consumeStatement()
		return
	}

	params = make([]ast.IdentifierExpression, len(pExprs))

	for i, pExpr := range pExprs {
		if pExpr.Type() != ast.IdentifierExpressionNode {
			errMsg := fmt.Sprintf(internal.ErrInvalidToken, token.IdentifierToken, pExpr.String())
			rdp.recordError(internal.NewError(token.NewMetatadata(pExpr.LineNumber(), pExpr.FileName()), errMsg,
				internal.SyntaxErr))
			rdp.consumeStatement()
			return
		}
		params[i] = *pExpr.(*ast.IdentifierExpression)
	}

	if body, err = rdp.parseBlockStatement(); err != nil {
		rdp.recordError(err)
		rdp.consumeBlockStatement()
		return
	}
	stmt = ast.NewFuntionDeclarationStatement(t, iden, body, params)
	return
}
