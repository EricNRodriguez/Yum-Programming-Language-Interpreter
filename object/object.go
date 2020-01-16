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
	Value int64
}

func NewInteger(i int64) *Integer {
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

type Float struct {
	Value float64
}

func NewFloat(f float64) *Float {
	return &Float{f}
}

func (f *Float) Type() ObjectType {
	return FLOAT
}

func (f *Float) Literal() string {
	return fmt.Sprintf("%f", f.Value)
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

type UserFunction struct {
	Name       string
	Parameters []string
	Body       []ast.Statement
}

func NewUserFunction(n string, params []string, body []ast.Statement) *UserFunction {
	return &UserFunction{
		Name:       n,
		Parameters: params,
		Body:       body,
	}
}

func (f *UserFunction) Type() ObjectType {
	return USER_FUNCTION
}

func (f *UserFunction) Literal() string {
	sBuff := bytes.Buffer{}
	sBuff.WriteString(fmt.Sprintf("func %v(%v) {", f.Name, strings.Join(f.Parameters, ", ")))
	for _, s := range f.Body {
		sBuff.WriteString(s.String())
	}
	sBuff.WriteString("};")
	return sBuff.String()
}

type NativeFunction struct {
	Name     string
	Function func(args ...Object) Object
}

func NewNativeFunction(n string, f func(args ...Object) Object) *NativeFunction {
	return &NativeFunction{
		Name:     n,
		Function: f,
	}
}

func (nf *NativeFunction) Type() ObjectType {
	return NATIVE_FUNCTION
}

func (nf *NativeFunction) Literal() string {
	return fmt.Sprintf("%v", *nf)
}
