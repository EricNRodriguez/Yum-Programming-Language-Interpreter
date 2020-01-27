package lexer

import "Yum-Programming-Language-Interpreter/token"

var keywords = map[string]token.TokenType{
	"func":   token.FuncToken,
	"var":    token.VarToken,
	"if":     token.IfToken,
	"else":   token.ElseToken,
	"return": token.ReturnToken,
	"true":   token.BooleanToken,
	"false":  token.BooleanToken,
	"while":  token.WhileToken,
}

func classifyTokenLiteral(s string) (t token.TokenType) {
	t, ok := keywords[s]
	if !ok {
		t = token.IdentifierToken
	}
	return
}
