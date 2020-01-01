package token

type Token interface {
	Type() TokenType
	Literal() string
	Data() Metadata
	Metadata
}

type token struct {
	tokenType TokenType
	literal   string
	Metadata
}

func NewToken(tt TokenType, l string, lN int, fN string) *token {
	return &token{
		tokenType:         tt,
		literal:           l,
		Metadata: NewMetatadata(lN, fN),
	}
}

func (t *token) Type() TokenType {
	return t.tokenType
}

func (t *token) Literal() string {
	return t.literal
}

func (t *token) Data() Metadata {
	return t.Metadata
}
