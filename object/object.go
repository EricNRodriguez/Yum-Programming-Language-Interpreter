package object

import (
	"Yum-Programming-Language-Interpreter/ast"
	"bytes"
	"fmt"
	"strings"
)

// CACHE TRUE FALSE AND NULL

type Object interface {
	Type() ObjectType
	Literal() string
}

type Integer struct {
	Value int
}

func NewInteger(i int) *Integer {
	return &Integer{
		Value: i,
	}
}

func (i *Integer) Type() ObjectType {
	return INTEGER
}

func (i *Integer) Literal() string {
	return fmt.Sprintf("%v", i.Value)
}

type Boolean struct {
	Value bool
}

func NewBoolean(b bool) *Boolean {
	return &Boolean{
		Value: b,
	}
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN
}

func (b *Boolean) Literal() string {
	return fmt.Sprintf("%v", b.Value)
}

type Null struct{}

func NewNull() *Null {
	return &Null{}
}

func (n *Null) Type() ObjectType {
	return NULL
}

func (n *Null) Literal() string {
	return "null"
}

type ReturnValue struct {
	Value Object
}

func NewReturnValue(o Object) *ReturnValue {
	return &ReturnValue{
		Value: o,
	}
}

func (r *ReturnValue) Type() ObjectType {
	return RETURN
}

func (r *ReturnValue) Literal() string {
	return r.Value.Literal()
}

type Function struct {
	Name string
	Parameters []string
	Body []ast.Statement
}

func NewFunction(n string, params []string, body []ast.Statement) *Function {
	return &Function{
		Name: n,
		Parameters: params,
		Body: body,
	}
}

func (f *Function) Type() ObjectType {
	return FUNCTION
}

func (f *Function) Literal() string {
	sBuff := bytes.Buffer{}
	sBuff.WriteString(fmt.Sprintf("func %v(%v) {", f.Name, strings.Join(f.Parameters, ", ")))
	for _, s := range f.Body {
		sBuff.WriteString(s.String())
	}
	sBuff.WriteString("};")
	return sBuff.String()
}