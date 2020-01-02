package ast

type NodeType int

const (
	PROGRAM    NodeType = iota
	IDENTIFIER          // should delete and just use identifier expression !
	PREFIX_EXPRESSION
	INFIX_EXPRESSION
	INTEGER_EXPRESSION
	BOOLEAN_EXPRESSION
	FUNC_CALL_EXPRESSION
	IDENTIFIER_EXPRESSION
	VAR_STATEMENT
	RETURN_STATEMENT
	EXPRESSION_STATEMENT
	IF_STATEMENT
	FUNCTION_DECLARATION_STATEMENT
)
