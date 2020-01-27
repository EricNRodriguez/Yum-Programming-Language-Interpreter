package internal

import (
	"Yum-Programming-Language-Interpreter/token"
	"fmt"
)

type ErrorType string

const (
	SyntaxErr ErrorType = "syntax error"
	RuntimeErr ErrorType = "runtime error"
	SemanticErr ErrorType = "semantic error"
	InternalErr ErrorType = "internal error"
)

type Error struct {
	token.Metadata
	msg  string
	code ErrorType
}

func NewError(md token.Metadata, msg string, code ErrorType) *Error {
	return &Error{md, msg, code}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v %v %v | %v", e.Type(), e.Metadata.FileName(), e.Metadata.LineNumber(), e.msg)
}

func (e *Error) Type() ErrorType {
	return e.code
}
