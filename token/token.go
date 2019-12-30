package token

type TokenInterface interface {
	Type() TokenType
	Literal() string
	Metadata() MetadataInterface
	MetadataInterface
}

type Token struct {
	tokenType TokenType
	literal   string
	MetadataInterface
}

func NewToken(tt TokenType, l string, lN int, fN string) *Token {
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

func (t *Token) Metadata() MetadataInterface {
	return t.MetadataInterface
}
