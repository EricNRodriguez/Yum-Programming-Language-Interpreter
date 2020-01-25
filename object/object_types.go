package object

type ObjectType int

func (oT ObjectType) String() string {
	switch oT {
	case INTEGER:
		return "INTEGER"
	case BOOLEAN:
		return "BOOLEAN"
	case RETURN:
		return "RETURN"
	case USER_FUNCTION:
		return "USER FUNCTION"
	case NATIVE_FUNCTION:
		return "NATIVE FUNCTION"
	case NULL:
		return "NULL"
	case FLOAT:
		return "FLOAT"
	case ARRAY:
		return "ARRAY"
	case STRING:
		return "STRING"
	default:
		return "UNKNOWN TYPE"

	}
}

const (
	INTEGER ObjectType = iota
	FLOAT
	BOOLEAN
	STRING
	RETURN
	USER_FUNCTION
	NATIVE_FUNCTION
	ARRAY
	NULL
)
