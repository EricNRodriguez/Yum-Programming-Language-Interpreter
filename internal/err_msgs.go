package internal

const (
	// syntax errors
	InvalidTokenErr          = "expected %v, received %v"
	InvalidPrefixOperatorErr = "%v is not a valid prefix operator"
	InvalidInfixOperatorErr  = "%v is not a valid infix operator"
	ErrEmptyFile             = "program file is empty"
	ErrInitParser            = "unable to initialise parser"
	ErrInvalidStatement      = "invalid statement beginning with %v"

	// semantic errors
	DeclaredVariableErr              = "%v already declared in current scope"
	ReturnLocationErr                = "unable to return outside of function"
	UndeclaredFunctionErr            = "%v not declared"
	DeclaredFunctionErr              = "%v already declared"
	InvalidFunctionCallParametersErr = "%v requires %v parameters, %v given"
	UndeclaredIdentifierErr          = "%v not declared"

	// runtime errors
	DivisionByZeroErr = "division by zero"
	TypeErr           = "%v not of type %v"
	MismatchedTypeErr = "mismatched types %v and %v"
	TypeOperationErr  = "operation %v not available for type %v"

	// internal error
	UnimplementedType = "unable to evaluate type %v"
)
