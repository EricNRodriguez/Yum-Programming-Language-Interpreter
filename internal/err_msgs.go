package internal

const (
	// syntax errors
	ErrInvalidToken          = "expected %v, received %v"
	ErrInvalidPrefixOperator = "%v is not a valid prefix operator"
	ErrInvalidInfixOperator  = "%v is not a valid infix operator"
	ErrEmptyFile             = "ProgramNode file is empty"
	ErrInitParser            = "unable to initialise parser"
	ErrInvalidStatement      = "invalid statement beginning with %v"
	ErrEndOfFile             = "unexpected EOF at line %v"

	// semantic errors
	ErrDeclaredVariable              = "%v already declared in current scope"
	ErrReturnLocation                = "unable to return outside of function"
	ErrUndeclaredFunction            = "%v not declared"
	ErrDeclaredFunction              = "%v declared in file"
	ErrInvalidFunctionCallParameters = "%v requires %v parameters, %v given"
	ErrUndeclaredIdentifierNode      = "%v not declared"

	// runtime errors
	ErrDivisionByZero   = "division by zero"
	ErrType             = "%v not of type %v"
	ErrTypeOperation    = "operation %v not available for type %v"
	ErrIndexOutOfBounds = "index out of bounds"
	ErrConditionType    = "condition does not evaluate to a boolean"

	// internal error
	ErrUnimplementedType = "unable to evaluate type %v"
	ErrFailedToReadFile = "failed to read file %v | %v"
)
