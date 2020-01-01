package lexer

import (
	"Yum-Programming-Language-Interpreter/token"
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

type LexerInterface interface {
	NextToken() (token.Token, error)
	Close() error
}

type Lexer struct {
	*bufio.Reader
	io.Closer
	currentLine       []byte
	currentLineNumber int
	currentLineIndex  int
	fileName          string
}

func NewLexer(f *os.File) (l LexerInterface, err error) {
	var (
		r    *bufio.Reader
		line []byte
	)

	r = bufio.NewReader(f)

	if line, _, err = r.ReadLine(); err != nil {
		err = errors.New(fmt.Sprintf("Empty file, unable to parse | %s", err.Error()))
		return
	}

	l = &Lexer{
		Reader:            r,
		Closer:            f,
		currentLineNumber: 1,
		currentLine:       line,
		currentLineIndex:  0,
		fileName:          f.Name(),
	}

	return

}

func (l *Lexer) readChars(n int) (chars []byte, err error) {
	chars = make([]byte, n)
	_, err = l.Reader.Read(chars)
	return
}

func (l *Lexer) validVariableNameCharacter(b byte) bool {
	return (b >= 65 && b <= 90) || (b >= 97 && b <= 122)
}

func (l *Lexer) readIdentifier() (idt []byte) {
	idt = make([]byte, 0)
	// ascii characters
	for l.currentLineIndex < len(l.currentLine) && l.validVariableNameCharacter(l.currentLine[l.currentLineIndex]) {
		idt = append(idt, l.currentLine[l.currentLineIndex])
		l.currentLineIndex++
	}
	return
}

func (l *Lexer) readNumber() (num []byte) {
	num = make([]byte, 0)
	// [0-9]
	for l.currentLineIndex < len(l.currentLine) && l.currentLine[l.currentLineIndex] >= 48 &&
		l.currentLine[l.currentLineIndex] <= 57 {
		num = append(num, l.currentLine[l.currentLineIndex])
		l.currentLineIndex++
	}
	return
}

func (l *Lexer) trailingTerminal() (t token.TokenType, ok bool) {
	if l.currentLineIndex < len(l.currentLine) {
		t = token.TokenType(l.currentLine[l.currentLineIndex])
		ok = true
	}
	l.currentLineIndex++
	return
}

func (l *Lexer) NextToken() (t token.Token, err error) {
	var s = ""

	// parsed the entire line
	if l.currentLineIndex >= len(l.currentLine) {
		// read in next line
		l.currentLine, _, err = l.ReadLine()
		l.currentLineNumber += 1

		// if EOF
		if err != nil {
			t = token.NewToken(token.EOF, "EOF", l.currentLineNumber, l.fileName)
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
	for s == " " && l.currentLineIndex < len(l.currentLine) {
		s = string(l.currentLine[l.currentLineIndex])
		l.currentLineIndex++

		if s == " " && l.currentLineIndex == len(l.currentLine) {
			// begin parsing new line
			l.currentLineIndex = 0
			return l.NextToken()
		}

	}

	switch token.TokenType(s) {
	case token.ADD:
		t = token.NewToken(token.ADD, s, l.currentLineNumber, l.fileName)
	case token.SUB:
		t = token.NewToken(token.SUB, s, l.currentLineNumber, l.fileName)
	case token.DIV:
		t = token.NewToken(token.DIV, s, l.currentLineNumber, l.fileName)
	case token.MULT:
		t = token.NewToken(token.MULT, s, l.currentLineNumber, l.fileName)
	case token.ASSIGN:
		tt, _ := l.trailingTerminal()
		switch tt {
		case token.ASSIGN:
			t = token.NewToken(token.EQUAL, s, l.currentLineNumber, l.fileName)
		default:
			// shift back, unread trailing terminal
			l.currentLineIndex--
			t = token.NewToken(token.ASSIGN, s+s, l.currentLineNumber, l.fileName)
		}
	case token.NEGATE:
		tt, _ := l.trailingTerminal()
		switch tt {
		case token.ASSIGN:
			t = token.NewToken(token.NEQUAL, s+string(tt), l.currentLineNumber, l.fileName)
		default:
			// shift back, unread trailing terminal
			l.currentLineIndex--
			t = token.NewToken(token.NEGATE, s, l.currentLineNumber, l.fileName)
		}

	case token.GTHAN:
		tt, _ := l.trailingTerminal()
		switch tt {
		case token.ASSIGN:
			t = token.NewToken(token.GTEQUAL, s+string(tt), l.currentLineNumber, l.fileName)
		default:
			// shift back, unread trailing terminal
			l.currentLineIndex--
			t = token.NewToken(token.GTHAN, s, l.currentLineNumber, l.fileName)
		}
	case token.LTHAN:
		tt, _ := l.trailingTerminal()
		switch tt {
		case token.ASSIGN:
			t = token.NewToken(token.LTEQUAL, s+string(tt), l.currentLineNumber, l.fileName)
		default:
			// shift back, unread trailing terminal
			l.currentLineIndex--
			t = token.NewToken(token.LTHAN, s, l.currentLineNumber, l.fileName)
		}
	case token.SEMICOLON:
		t = token.NewToken(token.SEMICOLON, s, l.currentLineNumber, l.fileName)
	case token.COMMA:
		t = token.NewToken(token.COMMA, s, l.currentLineNumber, l.fileName)
	case token.LPAREN:
		t = token.NewToken(token.LPAREN, s, l.currentLineNumber, l.fileName)
	case token.RPAREN:
		t = token.NewToken(token.RPAREN, s, l.currentLineNumber, l.fileName)
	case token.LBRACE:
		t = token.NewToken(token.LBRACE, s, l.currentLineNumber, l.fileName)
	case token.RBRACE:
		t = token.NewToken(token.RBRACE, s, l.currentLineNumber, l.fileName)
	case token.LBRACKET:
		t = token.NewToken(token.LBRACKET, s, l.currentLineNumber, l.fileName)
	case token.RBRACKET:
		t = token.NewToken(token.RBRACKET, s, l.currentLineNumber, l.fileName)
	case token.RETURN:
		t = token.NewToken(token.RETURN, s, l.currentLineNumber, l.fileName)

	default:
		// account for token literals, integers and illegal tokens

		l.currentLineIndex -= 1

		// recognise keywords, booleans and identifiers
		if l.validVariableNameCharacter(l.currentLine[l.currentLineIndex]) {
			idt := string(l.readIdentifier())
			idtType := classifyTokenLiteral(idt)
			t = token.NewToken(idtType, idt, l.currentLineNumber, l.fileName)

			// [0,9]
		} else if l.currentLine[l.currentLineIndex] >= 48 && l.currentLine[l.currentLineIndex] <= 57 {
			num := l.readNumber()
			t = token.NewToken(token.INT, string(num), l.currentLineNumber, l.fileName)

		} else {
			l.currentLineIndex++
			t = token.NewToken(token.ILLEGAL, s, l.currentLineNumber, l.fileName)
		}
	}
	return
}
