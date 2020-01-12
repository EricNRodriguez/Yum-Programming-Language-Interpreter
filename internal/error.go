package internal

import (
	"Yum-Programming-Language-Interpreter/token"
	"fmt"
)

type ErrorType int

func (eT ErrorType) Name() string {
	switch eT {
	case SyntaxErr:
		return "SYNTAX ERROR"
	case RuntimeErr:
		return "RUNTIME ERROR"
	case InternalErr:
		return "INTERNAL PROGRAM ERROR"
	case SemanticErr:
		return "SEMANTIC ERROR"
	default:
		return "UNKNOWN ERROR TYPE"
	}
}

const (
	SyntaxErr ErrorType = iota
	RuntimeErr
	SemanticErr
	InternalErr
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
	return fmt.Sprintf("%v %v %v | %v", e.Type().Name(), e.Metadata.FileName(), e.Metadata.LineNumber(), e.msg)
}

func (e *Error) Type() ErrorType {
	return e.code
}
