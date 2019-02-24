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

func builtinRange(vals ...data.TorinoValue) (data.TorinoValue, error) {
	for _, v := range vals {
		_, ok := v.(*data.TorinoInt)
		if !ok {
			return nil, errors.New("range takes integer arguments")
		}
	}

	var lo, hi, step int
	if len(vals) == 1 {
		lo = 0
		hi = vals[0].(*data.TorinoInt).Value
		step = 1
	} else if len(vals) == 2 {
		lo = vals[0].(*data.TorinoInt).Value
		hi = vals[1].(*data.TorinoInt).Value
		step = 1
	} else if len(vals) == 3 {
		lo = vals[0].(*data.TorinoInt).Value
		hi = vals[1].(*data.TorinoInt).Value
		step = vals[2].(*data.TorinoInt).Value
	} else {
		return nil, errors.New("range takes between one and three arguments")
	}

	lst := []data.TorinoValue{}
	for i := lo; i < hi; i += step {
		lst = append(lst, &data.TorinoInt{i})
	}
	return &data.TorinoList{lst}, nil
}
