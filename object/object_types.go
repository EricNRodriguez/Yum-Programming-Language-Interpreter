package object

type ObjectType int

const (
	INTEGER ObjectType = iota
	BOOLEAN
	RETURN
	FUNCTION
	NULL
)
