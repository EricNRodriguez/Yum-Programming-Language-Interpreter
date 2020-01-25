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
		if o[0].Type() != ARRAY {
			err = errors.New(fmt.Sprintf(internal.TypeErr, o[0].Type(), ARRAY))
			return
		}
		arr := o[0].(*Array)
		l = NewInteger(arr.Length)
		return
	})

	NativeFunctions = map[string]*NativeFunction{
		print.Name: print,
		length.Name: length,
	}

)

