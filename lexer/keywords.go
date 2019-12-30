package lexer

import "Yum-Programming-Language-Interpreter/token"

var keywords = map[string]token.TokenType{
	"func":   token.FUNC,
	"var":    token.VAR,
	"if":     token.IF,
	"else":   token.ELSE,
	"return": token.RETURN,
	"true":   token.BOOLEAN,
	"false":  token.BOOLEAN,
}

func classifyTokenLiteral(s string) (t token.TokenType) {
	t, ok := keywords[s]
	if !ok {
		t = token.IDEN
	}
	return
}
