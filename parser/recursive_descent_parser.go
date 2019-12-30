package parser

import (
	"Yum-Programming-Language-Interpreter/ast"
	"Yum-Programming-Language-Interpreter/lexer"
	"Yum-Programming-Language-Interpreter/token"
	"errors"
	"fmt"
)

type parseMethod func() ast.StatementInterface

type RecursiveDescentParser struct {
	parseMethodRouter map[token.TokenType]parseMethod
	prattParserInterface
}

func NewRecursiveDescentParser(l lexer.LexerInterface) (rdp *RecursiveDescentParser, err error) {
	var (
		pMR = make(map[token.TokenType]parseMethod)
		prattParser prattParserInterface
	)

	if prattParser, err = newPrattParser(l); err != nil {
		return
	}

	rdp = &RecursiveDescentParser{
		parseMethodRouter: pMR,
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
	prog = ast.NewProgram(rdp.currentToken().Metadata())

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

func (rdp *RecursiveDescentParser) parseStatement() (stmt ast.StatementInterface, err error) {
	if parser, ok := rdp.parseMethodRouter[rdp.currentToken().Type()]; !ok {
		err = errors.New(fmt.Sprintf("unable to parse statement beginning on line %v | no parse statement for " +
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

func (rdp *RecursiveDescentParser) parseVarStatement() (stmt ast.StatementInterface) {
	var (
		iden     token.TokenInterface
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

func (rdp *RecursiveDescentParser) parseReturnStatement() (stmt ast.StatementInterface) {
	var (
		retToken = rdp.currentToken()
	)

	rdp.consume(1) // consume return
	expr := rdp.parseExpression(MINPRECEDENCE)
	stmt = ast.NewReturnStatment(retToken, expr)
	return
}

func (rdp *RecursiveDescentParser) parseExpressionStatement() (stmt ast.StatementInterface) {
	expr := rdp.parseExpression(MINPRECEDENCE)
	stmt = ast.NewExpressionStatment(expr)
	return
}

func (rdp *RecursiveDescentParser) parseIfStatement() (stmt ast.StatementInterface) {
	var (
		t          = rdp.currentToken()
		condition  ast.ExpressionInterface
		trueBlock  *ast.BlockStatment
		falseBlock *ast.BlockStatment
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


func (rdp *RecursiveDescentParser) parseBlockStatement() (bStmt *ast.BlockStatment) {
	bStmt = ast.NewBlockStatement(rdp.currentToken())


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
			bStmt.AddStatement(stmt)
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


func (rdp *RecursiveDescentParser) parseFuncDeclarationStatement() (stmt ast.StatementInterface) {
	var (
		t      = rdp.currentToken()
		iden   token.TokenInterface
		params = make([]token.TokenInterface, 0)
		body   *ast.BlockStatment
		errs   []error
	)

	if !rdp.expectTokenType(token.IDEN) {
		rdp.consumeBlockStatement()
		return
	}
	rdp.consume(1) // consume func token

	iden = rdp.currentToken()

	if !rdp.expectTokenType(token.LPAREN) {
		rdp.consumeBlockStatement()
		return
	}
	rdp.consume(1)

	if !rdp.expectTokenType(token.IDEN) {
		rdp.consumeBlockStatement()
		return
	}
	rdp.consume(1) // consume left paren

	// parse parameters
	for rdp.currentToken().Type() != token.RPAREN && rdp.currentToken().Type() != token.EOF {

		// defensive check
		if err := rdp.currentToken().Type().AssertEqual(token.IDEN); err != nil {
			rdp.recordError(err)
			rdp.consumeBlockStatement()
			return
		}

		params = append(params, rdp.currentToken())
		rdp.consume(1)

		if rdp.currentToken().Type() != token.RPAREN && rdp.currentToken().Type() != token.EOF {
			if err := rdp.currentToken().Type().AssertEqual(token.COMMA); err != nil {
				rdp.recordError(err)
				rdp.consumeBlockStatement()
				return
			}
			if !rdp.expectTokenType(token.RPAREN) {
				rdp.consumeBlockStatement()
				return
			}
			rdp.consume(1)
		}
	}

	rdp.consume(1) // consume right paren

	if body = rdp.parseBlockStatement(); len(errs) != 0 {
		rdp.consumeBlockStatement()
		return
	}

	stmt = ast.NewFuntionDeclarationStatement(t, iden.Literal(), body, params...)
	return
}


func (rdp *RecursiveDescentParser) consumeIfStatement() {
	// move to next closing parenthesis
	for rdp.currentToken().Type() != token.RBRACE {
		rdp.consume(1)
	}
	rdp.consume(1) // move to token following }

	if rdp.currentToken().Type() == token.ELSE {
		for rdp.currentToken().Type() != token.RBRACE {
			rdp.consume(1)
		}
		rdp.consume(1)
	}

	return
}

func (rdp *RecursiveDescentParser) consumeBlockStatement() {
	for rdp.currentToken().Type() != token.RBRACE {
		if rdp.currentToken().Type() == token.EOF {
			err := errors.New(fmt.Sprintf("invalid block statement | expected RBRACE on line %v, recieved %v",
				rdp.currentToken().LineNumber(), rdp.currentToken().Type()))
			rdp.recordError(err)
			return
		}
		rdp.consume(1)
	}
	rdp.consume(1)
	return
}