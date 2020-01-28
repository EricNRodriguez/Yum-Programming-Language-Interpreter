package lexer

import (
	"Yum-Programming-Language-Interpreter/internal"
	"Yum-Programming-Language-Interpreter/token"
	"bufio"
	"bytes"
	"fmt"
	"github.com/spf13/afero"
	"io"
)

type Lexer interface {
	NextToken() (token.Token, error)
	Close() error
}

type lexer struct {
	*bufio.Reader
	io.Closer
	currentLine       []byte
	currentLineNumber int
	currentLineIndex  int
	fileName          string
	ignoreSpace       bool
}

func NewLexer(f afero.File) (l Lexer, err error) {
	var (
		r    *bufio.Reader
		line []byte
	)

	r = bufio.NewReader(f)

	if line, _, err = r.ReadLine(); err != nil {
		err = internal.NewError(token.NewMetatadata(0, f.Name()), internal.ErrEmptyFile, internal.SyntaxErr)
		return
	}

	l = &lexer{
		Reader:            r,
		Closer:            f,
		currentLineNumber: 1,
		currentLine:       line,
		currentLineIndex:  0,
		fileName:          f.Name(),
		ignoreSpace:       true,
	}

	return

}

func (l *lexer) readChars(n int) (chars []byte, err error) {
	chars = make([]byte, n)
	_, err = l.Reader.Read(chars)
	return
}

func (l *lexer) validVariableNameStartCharacter(b byte) bool {
	return (b >= 65 && b <= 90) || (b >= 97 && b <= 122)
}

func (l *lexer) validVariableNameCharacter(b byte) bool {
	return (b >= 65 && b <= 90) || (b >= 97 && b <= 122) || (b >= 48 && b <= 57)
}

func (l *lexer) readIdentifierNode() (idt []byte) {
	idt = make([]byte, 0)

	// ascii characters
	for l.currentLineIndex < len(l.currentLine) && l.validVariableNameCharacter(l.currentLine[l.currentLineIndex]) {
		idt = append(idt, l.currentLine[l.currentLineIndex])
		l.currentLineIndex++
	}
	return
}

func (l *lexer) readInt() string {
	// [0-9]
	numBuff := bytes.Buffer{}
	for l.currentLineIndex < len(l.currentLine) && l.currentLine[l.currentLineIndex] >= 48 &&
		l.currentLine[l.currentLineIndex] <= 57 {
		numBuff.WriteString(string(l.currentLine[l.currentLineIndex]))
		l.currentLineIndex++
	}
	return numBuff.String()
}

func (l *lexer) readNumber() (num string, ty token.TokenType) {
	num = l.readInt()
	ty = token.IntegerToken

	if l.currentLineIndex < len(l.currentLine) && l.currentLine[l.currentLineIndex] == 46 {
		l.currentLineIndex++
		num = fmt.Sprintf("%v.%v", num, l.readInt())
		ty = token.FloatingPointToken
	}
	return
}

func (l *lexer) trailingTerminal() (t token.TokenType, ok bool) {
	if l.currentLineIndex < len(l.currentLine) {
		t = token.TokenType(l.currentLine[l.currentLineIndex])
		ok = true
	}
	l.currentLineIndex++
	return
}

func (l *lexer) NextToken() (t token.Token, err error) {
	var s = ""

	// parsed the entire line
	if l.currentLineIndex >= len(l.currentLine) {
		// read in next line
		l.currentLine, _, err = l.ReadLine()
		l.currentLineNumber += 1

		// if EOF
		if err != nil {
			t = token.NewToken(token.EOFToken, "EOF", l.currentLineNumber, l.fileName)
			return
		} else {
			// begin parsing new line
			l.currentLineIndex = 0
			return l.NextToken()
		}
	}

	// next string
	s = string(l.currentLine[l.currentLineIndex])
	l.currentLineIndex++

	// ignore white space
	for l.ignoreSpace && s == " " && l.currentLineIndex < len(l.currentLine) {
		s = string(l.currentLine[l.currentLineIndex])
		l.currentLineIndex++

		if s == " " && l.currentLineIndex == len(l.currentLine) {
			// begin parsing new line
			l.currentLineIndex = 0
			return l.NextToken()
		}

	}

	switch token.TokenType(s) {
	case token.AddToken:
		t = token.NewToken(token.AddToken, s, l.currentLineNumber, l.fileName)
	case token.SubToken:
		t = token.NewToken(token.SubToken, s, l.currentLineNumber, l.fileName)
	case token.DivToken:
		t = token.NewToken(token.DivToken, s, l.currentLineNumber, l.fileName)
	case token.MultToken:
		t = token.NewToken(token.MultToken, s, l.currentLineNumber, l.fileName)
	case token.AssignToken:
		tt, _ := l.trailingTerminal()
		switch tt {
		case token.AssignToken:
			t = token.NewToken(token.EqualToken, s+s, l.currentLineNumber, l.fileName)
		default:
			// shift back, unread trailing terminal
			l.currentLineIndex--
			t = token.NewToken(token.AssignToken, s, l.currentLineNumber, l.fileName)
		}
	case token.NegateToken:
		tt, _ := l.trailingTerminal()
		switch tt {
		case token.AssignToken:
			t = token.NewToken(token.NotEqualToken, s+string(tt), l.currentLineNumber, l.fileName)
		default:
			// shift back, unread trailing terminal
			l.currentLineIndex--
			t = token.NewToken(token.NegateToken, s, l.currentLineNumber, l.fileName)
		}

	case token.GThanToken:
		tt, _ := l.trailingTerminal()
		switch tt {
		case token.AssignToken:
			t = token.NewToken(token.GThanEqualToken, s+string(tt), l.currentLineNumber, l.fileName)
		default:
			// shift back, unread trailing terminal
			l.currentLineIndex--
			t = token.NewToken(token.GThanToken, s, l.currentLineNumber, l.fileName)
		}
	case token.LThanToken:
		tt, _ := l.trailingTerminal()
		switch tt {
		case token.AssignToken:
			t = token.NewToken(token.LThanEqualToken, s+string(tt), l.currentLineNumber, l.fileName)
		default:
			// shift back, unread trailing terminal
			l.currentLineIndex--
			t = token.NewToken(token.LThanToken, s, l.currentLineNumber, l.fileName)
		}
	case token.SemicolonToken:
		t = token.NewToken(token.SemicolonToken, s, l.currentLineNumber, l.fileName)
	case token.CommaToken:
		t = token.NewToken(token.CommaToken, s, l.currentLineNumber, l.fileName)
	case token.QuotationMarkToken:
		t = token.NewToken(token.QuotationMarkToken, s, l.currentLineNumber, l.fileName)
		l.ignoreSpace = !l.ignoreSpace // allow strings to have white spaces
	case token.LeftParenToken:
		t = token.NewToken(token.LeftParenToken, s, l.currentLineNumber, l.fileName)
	case token.RightParenToken:
		t = token.NewToken(token.RightParenToken, s, l.currentLineNumber, l.fileName)
	case token.LeftBraceToken:
		t = token.NewToken(token.LeftBraceToken, s, l.currentLineNumber, l.fileName)
	case token.RightBraceToken:
		t = token.NewToken(token.RightBraceToken, s, l.currentLineNumber, l.fileName)
	case token.LeftBracketToken:
		t = token.NewToken(token.LeftBracketToken, s, l.currentLineNumber, l.fileName)
	case token.RightBracketToken:
		t = token.NewToken(token.RightBracketToken, s, l.currentLineNumber, l.fileName)
	case token.AndToken:
		t = token.NewToken(token.AndToken, s, l.currentLineNumber, l.fileName)
	case token.OrToken:
		t = token.NewToken(token.OrToken, s, l.currentLineNumber, l.fileName)
	case token.ReturnToken:
		t = token.NewToken(token.ReturnToken, s, l.currentLineNumber, l.fileName)

	default:
		// account for token literals, integers and illegal tokens

		l.currentLineIndex -= 1

		// recognise keywords, booleans and IdentifierNodes
		if l.validVariableNameStartCharacter(l.currentLine[l.currentLineIndex]) {
			idt := string(l.readIdentifierNode())
			idtType := classifyTokenLiteral(idt)
			t = token.NewToken(idtType, idt, l.currentLineNumber, l.fileName)

			// [0,9]
		} else if l.currentLine[l.currentLineIndex] >= 48 && l.currentLine[l.currentLineIndex] <= 57 {
			numStr, ty := l.readNumber()
			t = token.NewToken(ty, numStr, l.currentLineNumber, l.fileName)

		} else {
			l.currentLineIndex++
			t = token.NewToken(token.IllegalToken, s, l.currentLineNumber, l.fileName)
		}
	}
	return
}
