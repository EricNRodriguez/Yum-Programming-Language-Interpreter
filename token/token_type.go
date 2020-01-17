package token

type TokenType string

const (
	EOF     TokenType = "EOF"
	ILLEGAL TokenType = "ILLEGAL"

	FUNC   TokenType = "FUNC"
	VAR    TokenType = "VAR"
	IF     TokenType = "IF"
	ELSE   TokenType = "ELSE"
	RETURN TokenType = "RETURN"
	WHILE  TokenType = "WHILE"

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

	// string
	QUOTATION_MARK TokenType = "\""

	IDEN    TokenType = "IDENTIFIER"
	INT     TokenType = "INT"
	FLOAT   TokenType = "FLOAT"
	BOOLEAN TokenType = "BOOL"
	STRING  TokenType = "STRING"
)
