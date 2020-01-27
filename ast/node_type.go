package ast

type NodeType string

const (
	ProgramNode                      = "program expression"
	IdentifierExpressionNode         = "identifier expression"
	ArrayExpressionNode              = "array expression"
	ArrayIndexExpressionNode         = "array index expression"
	PrefixExpressionNode             = "prefix expression"
	InfixExpressionNode              = "infix expression"
	IntegerExpressionNode            = "integer expression"
	FloatingPointExpressionNode      = "floating point expression"
	StringExpressionNode             = "string expression"
	BooleanExpressionNode            = "boolean expression"
	FunctionCallExpressionNode       = "function call expression"
	VarStatementNode                 = "variable declaration statement"
	AssignmentStatementNode          = "assignment statement"
	ReturnStatementNode              = "return statement"
	WhileStatementNode               = "while statement"
	IfStatementNode                  = "if statement"
	FunctionDeclarationStatementNode = "function declaration statement"
	FunctionCallStatementNode        = "function call statement"
)
