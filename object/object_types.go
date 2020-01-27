package object

type ObjectType string

const (
	IntegerObject        = "integer"
	FloatingPointObject  = "float"
	BooleanObject        = "boolean"
	StringObject         = "string"
	ReturnObject         = "return"
	UserFunctionObject   = "user function"
	NativeFunctionObject = "native function"
	ArrayObject          = "ArrayNode"
	NullObject           = "null"
)
