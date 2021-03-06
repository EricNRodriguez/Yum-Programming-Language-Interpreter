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
	ErrInvalidIndexType              = "%v is not a valid index"

	// runtime errors
	ErrDivisionByZero   = "division by zero"
	ErrType             = "%v not of type %v"
	ErrTypeOperation    = "operation %v not available for type %v"
	ErrIndexOutOfBounds = "index out of bounds"
	ErrConditionType    = "condition does not evaluate to a boolean"

	// internal error
	ErrUnimplementedType = "unable to evaluate type %v"
	ErrFailedToReadFile  = "failed to read file %v | %v"

	// test errors
	ErrInvalidTokenTypeTest                   = "test case %v | token type %v received, expected %v"
	ErrInvalidTokenLiteralTest                = "test case %v | token literal %v received, expected %v"
	ErrInvalidASTNodeTypeTest                 = "test case %v | node type %v received, expected %v"
	ErrInvalidASTNodeLiteralTest              = "test case %v | node literal %v received, expected %v"
	ErrInvalidNumberOfErrorsTest              = "test case %v | expected %v errors, received %v"
	ErrInvalidNumberOfASTNodesTest            = "test case %v | expected %v ast nodes, received %v"
	ErrInvalidSyntaxSemanticAnalysisTestCases = "test case %v | invalid test case, %v syntax errors occurred"
	ErrInvalidSemanticsEvaluationTestCases    = "test case %v | invalid test case, %v semantic errors occurred"
	ErrInvalidSyntaxEvaluationTestCases       = "test case %v | invalid test case, %v syntax errors occurred"
	ErrMissingSymbolTest                      = "test case %v | %v not present after execution, expected %v : %v"
	ErrInvalidSymbolValueTest                 = "test case %v | expected %v : %v, received %v : %v"
	ErrUnexpectedRuntimeError                 = "test case %v | unexpected %v"

	// executing errors
	ErrFileNotProvided = "txt file required as argument"
	ErrFileNotFound    = "%v not found"
	ErrLoadFile        = "unable to load %v | %v"
)
