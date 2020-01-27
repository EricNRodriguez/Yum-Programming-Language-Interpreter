package object

import (
	"Yum-Programming-Language-Interpreter/internal"
	"errors"
	"fmt"
)

var (
	print = NewNativeFunction("print", -1, func(os ...Object) (o Object, err error) {
		for _, o := range os {
			fmt.Println(o.Literal())
		}
		return NewNull(), nil
	})

	length = NewNativeFunction("length", 1, func(o ...Object) (l Object, err error) {
		if o[0].Type() != ArrayObject {
			err = errors.New(fmt.Sprintf(internal.ErrType, o[0].Type(), ArrayObject))
			return
		}
		arr := o[0].(*ArrayNode)
		l = NewInteger(arr.Length)
		return
	})

	isNull = NewNativeFunction("isNull", 1, func(o ...Object) (l Object, err error) {
		l = NewBoolean(o[0].Type() == NullObject)
		return
	})

	NativeFunctions map[string]*NativeFunction
)

func init() {

	NativeFunctions = map[string]*NativeFunction{
		print.Name:  print,
		length.Name: length,
		isNull.Name: isNull,
	}
}
