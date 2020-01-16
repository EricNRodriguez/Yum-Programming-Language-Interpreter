package symbol_table

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
	SetUserFunc(*object.UserFunction)
	GetNativeFunc(string) (*object.NativeFunction, bool)
	DelUserFunc(string)
	GetUserFunc(string) (*object.UserFunction, bool)
	AvailableVar(string, bool) (ok bool)
	AvailableFunc(string) (ok bool)
	EnterFunction()
	ExitFunction()
	InFunctionCall() bool
}

type symbolTable struct {
	nameSpace            []map[string]object.Object
	functionDeclarations map[string]*object.UserFunction
	nativeFunctions      map[string]*object.NativeFunction
	cachedNameSpaces     [][]map[string]object.Object
	cachedScope          []int
	scope                int
}

func NewSymbolTable() *symbolTable {
	globalScope := make(map[string]object.Object)
	return &symbolTable{
		nameSpace:            []map[string]object.Object{globalScope}, // initialise global scope
		functionDeclarations: map[string]*object.UserFunction{},       // function declarations are global
		nativeFunctions:      object.NativeFunctions,
		scope:                0,
		cachedNameSpaces:     make([][]map[string]object.Object, 0),
		cachedScope:          make([]int, 0),
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
	s := st.scope
	for s >= 0 && !ok {
		if o, ok = st.nameSpace[s][name]; ok {
			st.nameSpace[s][name] = o
		}
		s--
	}
	return
}

// checks if var is available in the current scope
func (st *symbolTable) AvailableVar(name string, includeLowerScopes bool) bool {
	_, ok := st.nameSpace[st.scope][name]

	if includeLowerScopes {
		s := st.scope - 1
		for s >= 0 && !ok {
			_, ok = st.nameSpace[s][name]
			s--
		}
	}
	return !ok
}

func (st *symbolTable) SetUserFunc(f *object.UserFunction) {
	st.functionDeclarations[f.Name] = f
	return
}

// checks if func is available in the current scope
// native functions are not able to be overrided
func (st *symbolTable) AvailableFunc(name string) bool {
	_, okU := st.functionDeclarations[name]
	_, okN := st.nativeFunctions[name]
	return !(okU && okN)
}

func (st *symbolTable) DelUserFunc(string) {
	return
}

func (st *symbolTable) GetUserFunc(name string) (o *object.UserFunction, ok bool) {
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

func (st *symbolTable) EnterFunction() {
	st.cachedNameSpaces = append(st.cachedNameSpaces, st.nameSpace)
	st.cachedScope = append(st.cachedScope, st.scope)
	st.nameSpace = []map[string]object.Object{make(map[string]object.Object)}
	st.scope = 0
	return
}

func (st *symbolTable) ExitFunction() {
	st.nameSpace = st.cachedNameSpaces[len(st.cachedNameSpaces)-1]
	st.scope = st.cachedScope[len(st.cachedScope)-1]
	st.cachedNameSpaces = st.cachedNameSpaces[:len(st.cachedNameSpaces)-1]
	st.cachedScope = st.cachedScope[:len(st.cachedScope)-1]
	return
}

// true if current execution is within a function call, false if main file
func (st *symbolTable) InFunctionCall() bool {
	return len(st.cachedNameSpaces) != 0
}
