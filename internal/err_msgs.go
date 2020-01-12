package internal

const (
	ERR_INVALID_TOKEN                = "expected %v, received %v"
	ERR_INVALID_PREFIX_OPERATOR      = "%v is not a valid prefix operator"
	ERR_INVALID_INFIX_OPERATOR       = "%v is not a valid infix operator"
	ERR_EMPTY_FILE                   = "program file is empty"
	ERR_INIT_PARSER                  = "unable to initialise parser"
	ErrInvalidStatement              = "invalid statement beginning with %v"
	TypeErr                          = "%v is not of type %v"
	DeclaredVariableErr              = "%v already declared in current scope"
	ReturnLocationErr                = "unable to return outside of function"
	UndeclaredFunctionErr            = "%v not declared"
	DeclaredFunctionErr              = "%v already declared"
	InvalidFunctionCallParametersErr = "%v requires %v parameters, %v given"
	UndeclaredIdentifierErr = "%v not declared"
)
