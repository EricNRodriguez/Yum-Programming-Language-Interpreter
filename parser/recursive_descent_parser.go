package parser

import (
	"Yum-Programming-Language-Interpreter/ast"
	"Yum-Programming-Language-Interpreter/internal"
	"Yum-Programming-Language-Interpreter/lexer"
	"Yum-Programming-Language-Interpreter/token"
	"fmt"
)

type parseMethod func() ast.Statement

type RecursiveDescentParser struct {
	parseMethodRouter map[token.TokenType]parseMethod
	prattParserInterface
}

func NewRecursiveDescentParser(l lexer.Lexer) (Parser, error) {
	var (
		pMR         = make(map[token.TokenType]parseMethod) // parse method router
		prattParser prattParserInterface
		err         error
	)

	if prattParser, err = newPrattParser(l); err != nil {
		return nil, err
	}

	rdp := &RecursiveDescentParser{
		parseMethodRouter:    pMR,
		prattParserInterface: prattParser,
	}

	// initialise pMR
	pMR[token.VAR] = rdp.parseVarStatement
	pMR[token.RETURN] = rdp.parseReturnStatement
	pMR[token.IDEN] = rdp.parseIdenStatement
	pMR[token.IF] = rdp.parseIfStatement
	pMR[token.FUNC] = rdp.parseFuncDeclarationStatement

	return rdp, err
}

func (rdp *RecursiveDescentParser) Parse() (prog ast.Node) {

	stmts := make([]ast.Statement, 0)
	for rdp.checkNextToken() {

		if stmt, err := rdp.parseStatement(); err != nil {
			rdp.recordError(err)
			rdp.progressToNextSemicolon()

		} else {

			if stmt != nil {
				stmts = append(stmts, stmt)
			}

		}
		rdp.consume(1) // moving to next statement
	}

	prog = ast.NewProgram(rdp.currentToken().Data(), stmts...)

	// print all errors - dev purposes only
	for _, e := range rdp.errors() {
		fmt.Println(e.Error())
	}

	return prog
}

func (rdp *RecursiveDescentParser) parseStatement() (stmt ast.Statement, err error) {
	if parser, ok := rdp.parseMethodRouter[rdp.currentToken().Type()]; !ok {
		errMsg := fmt.Sprintf(internal.ErrInvalidStatement, rdp.currentToken().Literal())
		rdp.recordError(internal.NewError(rdp.currentToken().Data(), errMsg, internal.SyntaxErr))
		rdp.progressToNextSemicolon()
	} else {
		stmt = parser()
		if rdp.currentToken().Type() != token.SEMICOLON {
			errMsg := fmt.Sprintf(internal.ERR_INVALID_TOKEN, ";", rdp.currentToken().Literal())
			return nil, internal.NewError(rdp.currentToken().Data(), errMsg, internal.SyntaxErr)
		}
	}
	return
}

func (rdp *RecursiveDescentParser) parseVarStatement() (stmt ast.Statement) {
	var (
		iden     *ast.Identifier
		varToken = rdp.currentToken()
	)

	if !rdp.expectTokenType(token.IDEN) {
		rdp.progressToNextSemicolon()
		return
	}
	rdp.consume(1) // consume var

	iden = ast.NewIdentifier(rdp.currentToken())

	if !rdp.expectTokenType(token.ASSIGN) {
		rdp.progressToNextSemicolon()
		return
	}
	rdp.consume(2)

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
	if expr := rdp.parseExpression(MINPRECEDENCE); expr != nil {
		stmt = ast.NewReturnStatment(retToken, expr)

	}

	return
}

func (rdp *RecursiveDescentParser) parseIdenStatement() (stmt ast.Statement) {
	switch rdp.peekToken().Type() {
	case token.ASSIGN:
		iden := ast.NewIdentifier(rdp.currentToken())

		if !rdp.expectTokenType(token.ASSIGN) {
			rdp.progressToNextSemicolon()
			return
		}
		rdp.consume(2) // consume identifier and assign

		if expr := rdp.parseExpression(MINPRECEDENCE); expr != nil {
			stmt = ast.NewAssignmentStatement(iden.Metadata, iden, expr)
		}
	case token.LPAREN:
		stmt = rdp.parseFunctionCallStatement()
	default:
		errMsg := fmt.Sprintf(internal.ErrInvalidStatement, rdp.currentToken().Literal())
		rdp.recordError(internal.NewError(rdp.currentToken().Data(), errMsg, internal.SyntaxErr))
		rdp.progressToNextSemicolon()
	}

	return
}

func (rdp *RecursiveDescentParser) parseFunctionCallStatement() (stmt ast.Statement) {
	expr := rdp.parseExpression(MINPRECEDENCE).(*ast.FunctionCallExpression)
	stmt = ast.NewFunctionCallStatement(expr)
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
		//rdp.consumeCurrentStatement()
		rdp.consumeIfStatement()
		return
	}
	rdp.consume(1) // consume left paren

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
	if rdp.currentToken().Type() != token.LBRACE {
		errMsg := fmt.Sprintf(internal.ERR_INVALID_TOKEN, token.LBRACE, rdp.currentToken().Literal())
		rdp.recordError(internal.NewError(rdp.currentToken().Data(), errMsg, internal.SyntaxErr))
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

	if rdp.currentToken().Type() != token.RBRACE {
		errMsg := fmt.Sprintf(internal.ERR_INVALID_TOKEN, token.RBRACE, rdp.currentToken().Literal())
		rdp.recordError(internal.NewError(rdp.currentToken().Data(), errMsg, internal.SyntaxErr))
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
		params []ast.Identifier
		body   []ast.Statement
	)

	if !rdp.expectTokenType(token.IDEN) {
		rdp.consumeBlockStatement()
		return
	}
	rdp.consume(1) // consume func token

	iden = rdp.currentToken().Literal()
	rdp.consume(1) // consume function name

	ps := rdp.parseParameters()
	params = make([]ast.Identifier, len(ps))
	for i, p := range ps {
		if p.Type() != ast.IDENTIFIER_EXPRESSION {
			errMsg := fmt.Sprintf(internal.ERR_INVALID_TOKEN, token.IDEN, p.String())
			rdp.recordError(internal.NewError(token.NewMetatadata(p.LineNumber(), p.FileName()), errMsg, internal.SyntaxErr))
			rdp.consumeBlockStatement()
			return
		}
		params[i] = *p.(*ast.IdentifierExpression).Identifier
	}

	body = rdp.parseBlockStatement()
	stmt = ast.NewFuntionDeclarationStatement(t, iden, body, params)
	return
}
