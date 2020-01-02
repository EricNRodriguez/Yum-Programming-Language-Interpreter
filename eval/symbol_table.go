package eval

import (
	"Yum-Programming-Language-Interpreter/object"
)

type SymbolTable interface {
	EnterScope()
	ExitScope()
	SetVar(string, object.Object)
	DelVar(string)
	GetVar(string) (object.Object, bool)
	SetFunc(object.Function)
	DelFunc(string)
	GetFunc(string) (object.Function, bool)
}

type symbolTable struct {
	nameSpace []map[string]object.Object
	functionDeclarations map[string]object.Function
	scope int
}

func NewSymbolTable() *symbolTable {
	globalScope := make(map[string]object.Object)
	return &symbolTable{
		nameSpace: []map[string]object.Object{globalScope}, // initialise global scope
		functionDeclarations: map[string]object.Function{},
		scope: 0,
	}
}

func (st *symbolTable) SetVar(name string, object object.Object) {
	st.nameSpace[st.scope][name] = object
	return
}

func (st *symbolTable) DelVar(name string) {
	return
}

func (st *symbolTable) GetVar(name string) (o object.Object, ok bool) {
	o, ok = st.nameSpace[st.scope][name]
	s := st.scope - 1
	for s >= 0 && !ok {
		o, ok = st.nameSpace[st.scope][name]
		s--
	}
	return
}

func (st *symbolTable) SetFunc(f object.Function) {
	st.functionDeclarations[f.Name] = f
	return
}


func (st *symbolTable) DelFunc(string) {
	return
}

func (st *symbolTable) GetFunc(name string) (o object.Function, ok bool){
	o, ok = st.functionDeclarations[name]
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

