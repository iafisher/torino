package vm

import (
	"errors"
	"fmt"
	"github.com/iafisher/torino/data"
)

func builtinPrint(vals ...data.TorinoValue) (data.TorinoValue, error) {
	if len(vals) != 1 {
		return nil, errors.New("print takes one argument")
	}

	fmt.Print(vals[0].String())
	return &data.TorinoNone{}, nil
}

func builtinPrintln(vals ...data.TorinoValue) (data.TorinoValue, error) {
	if len(vals) != 1 {
		return nil, errors.New("println takes one argument")
	}

	fmt.Println(vals[0].String())
	return &data.TorinoNone{}, nil
}
