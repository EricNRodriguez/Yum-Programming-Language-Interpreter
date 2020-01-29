package token

type TokenType string

const (
	EOFToken     TokenType = "EOF"
	IllegalToken TokenType = "illegal token"

	FuncToken   TokenType = "func"
	VarToken    TokenType = "var"
	IfToken     TokenType = "if"
	ElseToken   TokenType = "else"
	ReturnToken TokenType = "return"
	WhileToken  TokenType = "while"

	// Arithmetic operations
	AddToken        TokenType = "+"
	SubToken        TokenType = "-"
	DivToken        TokenType = "/"
	MultToken       TokenType = "*"
	GThanToken      TokenType = ">"
	GThanEqualToken TokenType = ">="
	LThanToken      TokenType = "<"
	LThanEqualToken TokenType = "<="

	// boolean operators
	NegateToken TokenType = "!"
	AndToken    TokenType = "&"
	OrToken     TokenType = "|"

	// general operators
	AssignToken   TokenType = "="
	EqualToken    TokenType = "=="
	NotEqualToken TokenType = "!="

	// delimiters
	SemicolonToken TokenType = ";"
	CommaToken     TokenType = ","

	LeftParenToken  TokenType = "("
	RightParenToken TokenType = ")"

	LeftBraceToken  TokenType = "{"
	RightBraceToken TokenType = "}"

	LeftBracketToken  TokenType = "["
	RightBracketToken TokenType = "]"

	// string
	QuotationMarkToken TokenType = "\""

	IdentifierToken    TokenType = "identifier"
	IntegerToken       TokenType = "integer"
	FloatingPointToken TokenType = "floating point number"
	BooleanToken       TokenType = "boolean"
)
