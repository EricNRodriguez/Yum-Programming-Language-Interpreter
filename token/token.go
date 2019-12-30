package token

type TokenInterface interface {
	Type() TokenType
	Literal() string
	LineNumber() int
	FileName() string
}

type Token struct {
	tokenType TokenType
	literal   string
	MetadataInterface
}

func NewToken(tt TokenType, l string, lN int, fN string) TokenInterface {
	return &Token{
		tokenType:         tt,
		literal:           l,
		MetadataInterface: NewMetatadata(lN, fN),
	}
}

func (t *Token) Type() TokenType {
	return t.tokenType
}

func (t *Token) Literal() string {
	return t.literal
}
