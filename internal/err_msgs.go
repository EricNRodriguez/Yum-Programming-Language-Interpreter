package internal

const (
	// syntax errors
	InvalidTokenErr          = "expected %v, received %v"
	InvalidPrefixOperatorErr = "%v is not a valid prefix operator"
	InvalidInfixOperatorErr  = "%v is not a valid infix operator"
	ErrEmptyFile             = "program file is empty"
	ErrInitParser            = "unable to initialise parser"
	ErrInvalidStatement      = "invalid statement beginning with %v"
	EndOfFileErr = "unexpected EOF at line %v"

	// semantic errors
	DeclaredVariableErr              = "%v already declared in current scope"
	ReturnLocationErr                = "unable to return outside of function"
	UndeclaredFunctionErr            = "%v is not declared function"
	DeclaredFunctionErr              = "%v already declared in file"
	InvalidFunctionCallParametersErr = "%v requires %v parameters, %v given"
	UndeclaredVariableErr          = "%v not a declared variable"
	UnknownImportErr = "imported file %v not found"
	InvalidFileImportErr = "invalid file %v imported | %v"
	ImportedUndefinedFunctionErr = "imported undefined function %v.%v"
	ImportedDeclaredFunctionErr = "%v already declared/imported"


	// runtime errors
	DivisionByZeroErr   = "division by zero"
	TypeErr             = "%v not of type %v"
	TypeOperationErr    = "operation %v not available for type %v"
	IndexOutOfBoundsErr = "index out of bounds"
	ConditionTypeErr = "condition does not evaluate to a boolean"

	// internal error
	UnimplementedType = "unable to evaluate type %v"
)
