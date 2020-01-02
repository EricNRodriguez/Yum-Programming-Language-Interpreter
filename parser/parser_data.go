package parser

import (
	"Yum-Programming-Language-Interpreter/lexer"
	"Yum-Programming-Language-Interpreter/token"
	"errors"
	"fmt"
)

type parserDataInterface interface {
	consume(int)
	currentToken() token.Token
	expectTokenType(token.TokenType) bool
	peekToken() token.Token
	checkNextToken() bool
	recordError(error)
	errors() []error
	progressToNextSemicolon()
	consumeIfStatement()
	consumeBlockStatement()
}

type parserData struct {
	tokBuf       []token.Token
	currTok      token.Token
	syntaxErrors []error
}

func newParserData(l lexer.Lexer) (parserDataInterface, error) {
	var (
		cT  token.Token
		err error
	)

	if cT, err = l.NextToken(); err != nil {
		err = errors.New(fmt.Sprintf("unable to initalise parser, failed to read token | %s", err.Error()))
		return nil, err

	} else if cT.Type() == token.EOF {
		err = errors.New("unable to parse program with no content | EOF detected at start of program")
		return nil, err

	}

	pd := &parserData{
		tokBuf:       make([]token.Token, 0),
		currTok:      cT,
		syntaxErrors: make([]error, 0),
	}

	for cT.Type() != token.EOF {
		cT, _ = l.NextToken()
		pd.addToken(cT)
	}

	return pd, err
}

func (pd *parserData) addToken(t token.Token) {
	if t != nil {
		pd.tokBuf = append(pd.tokBuf, t)
	}
}

func (pd *parserData) consume(i int) {
	if pd.checkNextToken() {
		pd.currTok = pd.tokBuf[i-1]
		pd.tokBuf = pd.tokBuf[i:]
	}
	return
}

func (pd *parserData) currentToken() token.Token {
	return pd.currTok
}

func (pd *parserData) expectTokenType(e token.TokenType) bool {
	if err := pd.peekToken().Type().AssertEqual(e); err != nil {
		err = errors.New(fmt.Sprintf(" error on line %v | %v", pd.peekToken().LineNumber(), err.Error()))
		pd.recordError(err)
		return false
	}
	return true
}

func (pd *parserData) peekToken() (t token.Token) {
	if len(pd.tokBuf) > 0 {
		t = pd.tokBuf[0]
	}
	return
}

func (pd *parserData) checkNextToken() bool {
	return len(pd.tokBuf) > 0
}

func (pd *parserData) recordError(err error) {
	if err != nil {
		pd.syntaxErrors = append(pd.syntaxErrors, err)
	}
}

func (pd *parserData) errors() []error {
	return pd.syntaxErrors
}

func (pd *parserData) progressToNextSemicolon() {
	for pd.currTok.Type() != token.SEMICOLON && pd.currTok.Type() != token.EOF {
		pd.consume(1)
	}
}

func (pd *parserData) consumeBlockStatement() {
	for pd.currentToken().Type() != token.RBRACE {
		if pd.currentToken().Type() == token.EOF {
			err := errors.New(fmt.Sprintf("invalid block statement | expected RBRACE on line %v, recieved %v",
				pd.currentToken().LineNumber(), pd.currentToken().Type()))
			pd.recordError(err)
			return
		}
		pd.consume(1)
	}
	pd.consume(1)
	return
}

func (pd *parserData) consumeIfStatement() {
	// move to next closing parenthesis
	for pd.currentToken().Type() != token.RBRACE && pd.currentToken().Type() != token.EOF {
		pd.consume(1)
	}
	pd.consume(1) // move to token following }

	if pd.currentToken().Type() == token.ELSE {
		for pd.currentToken().Type() != token.RBRACE {
			pd.consume(1)
		}
		pd.consume(1)
	}

	return
}
