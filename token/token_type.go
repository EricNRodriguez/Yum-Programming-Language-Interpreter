package token

import (
	"errors"
	"fmt"
)

type TokenType string

func (tt TokenType) AssertEqual(ttTwo TokenType) (err error) {
	if tt != ttTwo {
		err = errors.New(fmt.Sprintf(" %v != %v", tt, ttTwo))
	}
	return
}

const (
	EOF     TokenType = "EOF"
	ILLEGAL TokenType = "ILLEGAL"

	IDEN   TokenType = "IDEN"
	FUNC   TokenType = "FUNC"
	VAR    TokenType = "VAR"
	IF     TokenType = "IF"
	ELSE   TokenType = "ELSE"
	RETURN TokenType = "RETURN"

	// Arithmetic operations
	ADD     TokenType = "+"
	SUB     TokenType = "-"
	DIV     TokenType = "/"
	MULT    TokenType = "*"
	GTHAN   TokenType = ">"
	GTEQUAL TokenType = ">="
	LTHAN   TokenType = "<"
	LTEQUAL TokenType = "<="

	// boolean operators
	NEGATE TokenType = "!"
	AND    TokenType = "&"
	OR     TokenType = "|"

	// general operators
	ASSIGN TokenType = "="
	EQUAL  TokenType = "=="
	NEQUAL TokenType = "!="

	// delimiters
	SEMICOLON TokenType = ";"
	COMMA     TokenType = ","

	LPAREN TokenType = "("
	RPAREN TokenType = ")"

	LBRACE TokenType = "{"
	RBRACE TokenType = "}"

	LBRACKET TokenType = "["
	RBRACKET TokenType = "]"

	INT TokenType = "INT"

	BOOLEAN TokenType = "BOOL"
)
