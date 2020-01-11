package eval

import (
	"Yum-Programming-Language-Interpreter/object"
)

type SymbolTable interface {
	EnterScope()
	ExitScope()
	SetVar(string, object.Object)
	UpdateVar(string, object.Object)
	DelVar(string)
	GetVar(string) (object.Object, bool)
	SetUserFunc(object.UserFunction)
	GetNativeFunc(string) (*object.NativeFunction, bool)
	DelUserFunc(string)
	GetUserFunc(string) (object.UserFunction, bool)
}

type symbolTable struct {
	nameSpace            []map[string]object.Object
	functionDeclarations map[string]object.UserFunction
	nativeFunctions      map[string]*object.NativeFunction
	scope                int
}

func NewSymbolTable() *symbolTable {
	globalScope := make(map[string]object.Object)
	return &symbolTable{
		nameSpace:            []map[string]object.Object{globalScope}, // initialise global scope
		functionDeclarations: map[string]object.UserFunction{},
		nativeFunctions:      object.NativeFunctions,
		scope:                0,
	}
}

func (st *symbolTable) SetVar(name string, object object.Object) {
	st.nameSpace[st.scope][name] = object
	return
}

func (st *symbolTable) UpdateVar(name string, o object.Object) {
	var (
		s  = st.scope
		ok bool
	)

	for s >= 0 && !ok {
		if _, ok := st.nameSpace[s][name]; ok {
			st.nameSpace[s][name] = o
		}
		s--
	}
	return
}

func (st *symbolTable) DelVar(name string) {
	return
}

func (st *symbolTable) GetVar(name string) (o object.Object, ok bool) {
	o, ok = st.nameSpace[st.scope][name]
	s := st.scope + 1
	for s >= 0 && !ok {
		o, ok = st.nameSpace[s][name]
		s++
	}
	return
}

func (st *symbolTable) SetUserFunc(f object.UserFunction) {
	st.functionDeclarations[f.Name] = f
	return
}

func (st *symbolTable) DelUserFunc(string) {
	return
}

func (st *symbolTable) GetUserFunc(name string) (o object.UserFunction, ok bool) {
	o, ok = st.functionDeclarations[name]
	return
}

func (st *symbolTable) GetNativeFunc(name string) (o *object.NativeFunction, ok bool) {
	o, ok = st.nativeFunctions[name]
	return
}

func (st *symbolTable) EnterScope() {
	st.scope++
	st.nameSpace = append(st.nameSpace, make(map[string]object.Object))
	return
}

func (st *symbolTable) ExitScope() {
	st.nameSpace = st.nameSpace[0:st.scope]
	st.scope--
	return
}
