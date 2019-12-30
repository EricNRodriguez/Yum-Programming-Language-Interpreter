package parser

import (
	"Yum-Programming-Language-Interpreter/lexer"
	"Yum-Programming-Language-Interpreter/token"
	"fmt"
	"errors"
)

type parserDataInterface interface {
	consume(int)
	currentToken() token.TokenInterface
	expectTokenType(token.TokenType) bool
	peekToken() token.TokenInterface
	checkNextToken() bool
	recordError(error)
	errors() []error
	progressToNextSemicolon()
}

type parserData struct {
	tokBuf []token.TokenInterface
	currTok token.TokenInterface
	syntaxErrors []error
}

func newParserData(l lexer.LexerInterface) (pd *parserData, err error) {
	var (
		cT token.TokenInterface
	)

	if cT, err = l.NextToken(); err != nil {
		err = errors.New(fmt.Sprintf("unable to initalise parser, failed to read token | %s", err.Error()))
		return

	} else if cT.Type() == token.EOF {
		err = errors.New("unable to parse program with no content | EOF detected at start of program")
		return

	}

	pd =  &parserData{
		tokBuf: make([]token.TokenInterface, 0),
		currTok: cT,
		syntaxErrors: make([]error, 0),
	}

	for cT.Type() != token.EOF {
		cT, _ = l.NextToken()
		pd.addToken(cT)
	}

	return
}

func (pd *parserData) addToken(t token.TokenInterface) {
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

func(pd *parserData) currentToken() token.TokenInterface {
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

func (pd *parserData) peekToken() (t token.TokenInterface) {
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



