package token

type Metadata interface {
	LineNumber() int
	FileName() string
}

type metadata struct {
	lineNumber int
	fileName   string
}

func NewMetatadata(lN int, fN string) Metadata {
	return &metadata{
		lineNumber: lN,
		fileName:   fN,
	}
}

func (m *metadata) LineNumber() int {
	return m.lineNumber
}

func (m *metadata) FileName() string {
	return m.fileName
}
