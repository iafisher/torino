package vm

import (
	"fmt"
	"github.com/iafisher/torino/data"
)

func builtinPrint(vals ...data.TorinoValue) data.TorinoValue {
	if len(vals) != 1 {
		panic("print takes one argument")
	}

	fmt.Print(vals[0].String())
	return &data.TorinoNone{}
}

func builtinPrintln(vals ...data.TorinoValue) data.TorinoValue {
	if len(vals) != 1 {
		panic("println takes one argument")
	}

	fmt.Println(vals[0].String())
	return &data.TorinoNone{}
}
