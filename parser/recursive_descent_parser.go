package parser

import (
	"Yum-Programming-Language-Interpreter/ast"
	"Yum-Programming-Language-Interpreter/lexer"
	"Yum-Programming-Language-Interpreter/token"
	"errors"
	"fmt"
)

type parseMethod func() ast.Statement

type RecursiveDescentParser struct {
	parseMethodRouter map[token.TokenType]parseMethod
	prattParserInterface
}

func NewRecursiveDescentParser(l lexer.LexerInterface) (rdp *RecursiveDescentParser, err error) {
	var (
		pMR         = make(map[token.TokenType]parseMethod)
		prattParser prattParserInterface
	)

	if prattParser, err = newPrattParser(l); err != nil {
		return
	}

	rdp = &RecursiveDescentParser{
		parseMethodRouter:    pMR,
		prattParserInterface: prattParser,
	}

	// initialise pMR
	pMR[token.VAR] = rdp.parseVarStatement
	pMR[token.RETURN] = rdp.parseReturnStatement
	pMR[token.IDEN] = rdp.parseExpressionStatement
	pMR[token.IF] = rdp.parseIfStatement
	pMR[token.FUNC] = rdp.parseFuncDeclarationStatement

	return
}

func (rdp *RecursiveDescentParser) Parse() (prog *ast.Program) {
	prog = ast.NewProgram(rdp.currentToken().Data())

	for rdp.checkNextToken() {
		if stmt, err := rdp.parseStatement(); err != nil {
			rdp.recordError(err)
			rdp.progressToNextSemicolon()
		} else {
			prog.AddStatement(stmt)
		}

		rdp.consume(1) // moving to next statement
	}

	// print all errors - dev purposes only
	for _, e := range rdp.errors() {
		fmt.Println("ERROR : ", e.Error())
	}

	return
}

func (rdp *RecursiveDescentParser) parseStatement() (stmt ast.Statement, err error) {
	if parser, ok := rdp.parseMethodRouter[rdp.currentToken().Type()]; !ok {
		err = errors.New(fmt.Sprintf("unable to parse statement beginning on line %v | no parse statement for "+
			"token type %v", rdp.currentToken().LineNumber(), rdp.currentToken().Type()))
	} else {
		cT := rdp.currentToken()
		stmt = parser()
		if err := rdp.currentToken().Type().AssertEqual(token.SEMICOLON); err != nil {
			err = errors.New(fmt.Sprintf("missing semicolon in statement beginning on line %v | %v",
				cT.LineNumber(), err))
		}
	}
	return
}

func (rdp *RecursiveDescentParser) parseVarStatement() (stmt ast.Statement) {
	var (
		iden     token.Token
		varToken = rdp.currentToken()
	)

	if !rdp.expectTokenType(token.IDEN) {
		rdp.progressToNextSemicolon()
		return
	}
	rdp.consume(1) // consume var

	iden = rdp.currentToken()

	if !rdp.expectTokenType(token.ASSIGN) {
		rdp.progressToNextSemicolon()
		return
	}
	rdp.consume(2) // consume identifier and assign

	// skip stmts with syntax errors
	if expr := rdp.parseExpression(MINPRECEDENCE); expr != nil {
		stmt = ast.NewVarStatement(varToken, iden, expr)
	}

	return
}

func (rdp *RecursiveDescentParser) parseReturnStatement() (stmt ast.Statement) {
	var (
		retToken = rdp.currentToken()
	)

	rdp.consume(1) // consume return
	expr := rdp.parseExpression(MINPRECEDENCE)
	stmt = ast.NewReturnStatment(retToken, expr)
	return
}

func (rdp *RecursiveDescentParser) parseExpressionStatement() (stmt ast.Statement) {
	expr := rdp.parseExpression(MINPRECEDENCE)
	stmt = ast.NewExpressionStatment(expr)
	return
}

func (rdp *RecursiveDescentParser) parseIfStatement() (stmt ast.Statement) {
	var (
		t          = rdp.currentToken()
		condition  ast.Expression
		trueBlock  []ast.Statement
		falseBlock []ast.Statement
	)

	if !rdp.expectTokenType(token.LPAREN) {
		rdp.consumeIfStatement()
		return
	}
	rdp.consume(1)

	if condition = rdp.parseExpression(MINPRECEDENCE); condition == nil {
		rdp.consumeIfStatement()
		return
	}

	trueBlock = rdp.parseBlockStatement()

	// else
	if rdp.currentToken().Type() == token.ELSE {
		rdp.consume(1) // consume ELSE
		falseBlock = rdp.parseBlockStatement()
	}

	stmt = ast.NewIfStatement(t, condition, trueBlock, falseBlock)
	return
}

func (rdp *RecursiveDescentParser) parseBlockStatement() (bStmt []ast.Statement) {
	bStmt = make([]ast.Statement, 0)
	if err := rdp.currentToken().Type().AssertEqual(token.LBRACE); err != nil {
		rdp.recordError(err)
		rdp.consumeIfStatement()
		return
	}

	rdp.consume(1)

	for rdp.currentToken().Type() != token.RBRACE && rdp.currentToken().Type() != token.EOF {
		if stmt, err := rdp.parseStatement(); err != nil {
			rdp.recordError(err)
			rdp.progressToNextSemicolon() // skip the rest of the current statement
		} else {
			bStmt = append(bStmt, stmt)
		}
		rdp.consume(1)
	}

	if err := rdp.currentToken().Type().AssertEqual(token.RBRACE); err != nil {
		rdp.recordError(err)
		rdp.consumeBlockStatement()
	}

	// consume right brace
	rdp.consume(1)
	return
}

func (rdp *RecursiveDescentParser) parseFuncDeclarationStatement() (stmt ast.Statement) {
	var (
		t      = rdp.currentToken()
		iden   string
		params = make([]*ast.Identifier, 0)
		body   []ast.Statement
		err error
	)

	if !rdp.expectTokenType(token.IDEN) {
		rdp.consumeBlockStatement()
		return
	}
	rdp.consume(1) // consume func token

	iden = rdp.currentToken().Literal()

	if !rdp.expectTokenType(token.LPAREN) {
		rdp.consumeBlockStatement()
		return
	}

	rdp.consume(2) // consume function name and left paren

	for {
		if err = rdp.currentToken().Type().AssertEqual(token.IDEN); err != nil {
			err = errors.New(fmt.Sprintf("invalid function declaration on line %v | %v",
				rdp.currentToken().LineNumber(), err))
			rdp.recordError(err)
			rdp.consumeBlockStatement()
			return
		}

		params = append(params, ast.NewIdentifier(rdp.currentToken()))
		rdp.consume(1) // consume ident

		if rdp.currentToken().Type() == token.RPAREN { break }

		if err = rdp.currentToken().Type().AssertEqual(token.COMMA); err != nil {
			err := errors.New(fmt.Sprintf("invalid function declaration on line %v | %v",
				rdp.currentToken().LineNumber(), err))
			rdp.recordError(err)
			rdp.consumeBlockStatement()
			return
		}
		rdp.consume(1) // consume comma

	}
	rdp.consume(1) // consume right paren


	body = rdp.parseBlockStatement()
	stmt = ast.NewFuntionDeclarationStatement(t, iden, body, params...)
	return
}
