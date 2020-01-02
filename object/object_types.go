package object

type ObjectType int

const (
	INTEGER ObjectType = iota
	BOOLEAN
	RETURN
	USER_FUNCTION
	NATIVE_FUNCTION
	NULL
)
