package token

type MetadataInterface interface {
	LineNumber() int
	FileName() string
}

type Metadata struct {
	lineNumber int
	fileName   string
}

func NewMetatadata(lN int, fN string) MetadataInterface {
	return &Metadata{
		lineNumber: lN,
		fileName:   fN,
	}
}

func (m *Metadata) LineNumber() int {
	return m.lineNumber
}

func (m *Metadata) FileName() string {
	return m.fileName
}
