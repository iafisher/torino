package eval

import (
	"github.com/iafisher/torino/data"
	"github.com/iafisher/torino/vm"
	"testing"
)

func TestEvalLetAndAssign(t *testing.T) {
	input := `
let abc = 666
abc = 42
abc
`
	val := evalHelper(input)
	checkInteger(t, val, 42)
}

func TestEvalArithmetic(t *testing.T) {
	val := evalHelper("(42 * (1 + 2 - 1)) / 2")
	checkInteger(t, val, 42)
}

func TestLetWithComplexArithmetic(t *testing.T) {
	input := `
let eighty = 40 * 2
let my_variable = (eighty + 6) / (1 + 1) - 1
my_variable
`
	val := evalHelper(input)
	checkInteger(t, val, 42)
}

func TestEvalIfStatement(t *testing.T) {
	input := `
let x = 42
if x > 100 {
	x = 666
}
x
`
	val := evalHelper(input)
	checkInteger(t, val, 42)
}

func TestEvalIfElseStatement(t *testing.T) {
	input := `
let x = 0
if true {
	x = 42
} else {
	x = 666
}
x
`
	val := evalHelper(input)
	checkInteger(t, val, 42)
}

func TestEvalIfElifStatement(t *testing.T) {
	input := `
let x = 0
if false {
	x = 666
} elif x == 0 {
	x = 42
}
x
`
	val := evalHelper(input)
	checkInteger(t, val, 42)
}

func TestEvalIfElifElseStatement(t *testing.T) {
	input := `
let x = 0
if false {
	x = 666
} elif x == 1 {
	x = 667
} else {
	x = 42
}
x
`
	val := evalHelper(input)
	checkInteger(t, val, 42)
}

// Helper functions

func evalHelper(text string) data.TorinoValue {
	env := vm.NewEnv()
	return Eval(text, env)
}

func checkInteger(t *testing.T, val data.TorinoValue, expected int64) {
	intVal, ok := val.(*data.TorinoInt)
	if !ok {
		t.Fatalf("Wrong Torino type: expected *TorinoInt, got %T", val)
	}

	if intVal.Value != expected {
		t.Fatalf("Wrong integer value: expected %d, got %d", expected, intVal.Value)
	}
}
