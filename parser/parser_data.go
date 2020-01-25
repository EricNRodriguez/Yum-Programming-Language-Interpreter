package parser

import (
	"Yum-Programming-Language-Interpreter/internal"
	"Yum-Programming-Language-Interpreter/lexer"
	"Yum-Programming-Language-Interpreter/token"
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
	consumeBlockStatement()
	consumeProgram()
	consumeStatement()
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
		err = internal.NewError(cT.Data(), internal.ErrInitParser, internal.InternalErr)
		return nil, err

	} else if cT.Type() == token.EOF {
		err = internal.NewError(cT.Data(), internal.ErrEmptyFile, internal.SyntaxErr)
		return nil, err

	}

	pd := &parserData{
		tokBuf:       make([]token.Token, 0),
		currTok:      cT,
		syntaxErrors: make([]error, 0),
	}

	for cT.Type() != token.EOF {
		fmt.Println(cT)
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

func (pd *parserData) expectTokenType(e token.TokenType) (b bool) {
	if pd.peekToken().Type() != e {
		errMsg := fmt.Sprintf(internal.InvalidTokenErr, e, pd.peekToken().Type())
		pd.recordError(internal.NewError(pd.peekToken().Data(), errMsg, internal.SyntaxErr))
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

// need to update to account for nested block statement s
func (pd *parserData) consumeBlockStatement() {
	for pd.currentToken().Type() != token.RBRACE {
		if pd.currentToken().Type() == token.EOF {
			errMsg := fmt.Sprintf(internal.InvalidTokenErr, token.RBRACE, pd.currentToken().Literal())
			pd.recordError(internal.NewError(pd.currentToken().Data(), errMsg, internal.SyntaxErr))
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

func (pd *parserData) consumeProgram() {
	pd.currTok = pd.tokBuf[len(pd.tokBuf)-2]
	pd.tokBuf = pd.tokBuf[len(pd.tokBuf)-1:len(pd.tokBuf)]
	return
}

func (pd *parserData) consumeStatement() {
	for pd.currentToken().Type() != token.SEMICOLON && pd.currentToken().Type() != token.LBRACE &&
		pd.currentToken().Type() != token.EOF {
		pd.consume(1)
	}

	if pd.currentToken().Type() == token.LBRACE {
		pd.consumeIfStatement()
	}

}
