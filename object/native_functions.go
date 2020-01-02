package object

import "fmt"

var NativeFunctions = map[string]*NativeFunction{
	yumPrint.Name: yumPrint,
}

var yumPrint = NewNativeFunction("print", func(os ...Object) Object {
	for _, o := range os {
		fmt.Println(o.Literal())
	}
	return NewNull()
})
