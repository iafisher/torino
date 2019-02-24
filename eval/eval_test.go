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
	val := evalHelper(t, input)
	checkInteger(t, val, 42)
}

func TestEvalArithmetic(t *testing.T) {
	val := evalHelper(t, "(42 * (1 + 2 - 1)) / 2")
	checkInteger(t, val, 42)
}

func TestLetWithComplexArithmetic(t *testing.T) {
	input := `
let eighty = 40 * 2
let my_variable = (eighty + 6) / (1 + 1) - 1
my_variable
`
	val := evalHelper(t, input)
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
	val := evalHelper(t, input)
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
	val := evalHelper(t, input)
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
	val := evalHelper(t, input)
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
	val := evalHelper(t, input)
	checkInteger(t, val, 42)
}

func TestEvalFunctionDeclaration(t *testing.T) {
	input := `
fn return42() {
	return 42
}
let x = return42()
x
`
	val := evalHelper(t, input)
	checkInteger(t, val, 42)
}

func TestEvalFunctionWithGlobalVariable(t *testing.T) {
	input := `
let FORTY_TWO = 42

fn return42() {
	return FORTY_TWO
}
return42()
`
	val := evalHelper(t, input)
	checkInteger(t, val, 42)
}

func TestEvalWhileLoop(t *testing.T) {
	input := `
let x = 0
while x < 42 {
	x = x + 1
}
x
`
	val := evalHelper(t, input)
	checkInteger(t, val, 42)
}

func TestEvalList(t *testing.T) {
	val := evalHelper(t, "[1, 2, 3]")

	listVal := checkList(t, val, 3)
	checkInteger(t, listVal.Values[0], 1)
	checkInteger(t, listVal.Values[1], 2)
	checkInteger(t, listVal.Values[2], 3)
}

func TestEvalIndex(t *testing.T) {
	val := evalHelper(t, "[\"a\", \"b\", \"c\"][2]")

	checkString(t, val, "c")
}

func TestEvalMap(t *testing.T) {
	val := evalHelper(t, "{\"one\": 1, \"two\": 1+1}")

	mapVal := checkMap(t, val, 2)

	first, ok := mapVal.Get(&data.TorinoString{"one"})
	if !ok {
		t.Fatalf("Expected \"one\" to be in the map")
	}
	checkInteger(t, first, 1)

	second, ok := mapVal.Get(&data.TorinoString{"two"})
	if !ok {
		t.Fatalf("Expected \"two\" to be in the map")
	}
	checkInteger(t, second, 2)
}

func TestEvalIndexMap(t *testing.T) {
	input := `
let m = {"one": 1}
let x = m["one"]
x
`
	val := evalHelper(t, input)

	checkInteger(t, val, 1)
}

func TestEvalIndexStr(t *testing.T) {
	val := evalHelper(t, "\"abc\"[0]")

	checkString(t, val, "a")
}

func TestEvalForLoop(t *testing.T) {
	input := `
let x = 0
for i in range(6) {
	x = x + 7
}
x
`
	val := evalHelper(t, input)

	checkInteger(t, val, 42)
}

// Helper functions

func evalHelper(t *testing.T, text string) data.TorinoValue {
	env := vm.NewEnv(nil)
	val, err := Eval(text, env)
	if err != nil {
		t.Fatalf("Eval error: %s", err)
	}
	return val
}

func checkInteger(t *testing.T, val data.TorinoValue, expected int) {
	intVal, ok := val.(*data.TorinoInt)
	if !ok {
		t.Fatalf("Wrong Torino type: expected *TorinoInt, got %T", val)
	}

	if intVal.Value != expected {
		t.Fatalf("Wrong integer value: expected %d, got %d", expected, intVal.Value)
	}
}

func checkString(t *testing.T, val data.TorinoValue, expected string) {
	strVal, ok := val.(*data.TorinoString)
	if !ok {
		t.Fatalf("Wrong Torino type: expected *TorinoString, got %T", val)
	}

	if strVal.Value != expected {
		t.Fatalf("Wrong string value: expected %q, got %q", expected, strVal.Value)
	}
}

func checkList(t *testing.T, val data.TorinoValue, nelems int) *data.TorinoList {
	listVal, ok := val.(*data.TorinoList)
	if !ok {
		t.Fatalf("Wrong Torino type: expected *TorinoList, got %T", val)
	}

	if len(listVal.Values) != nelems {
		t.Fatalf("Wrong number of list elements: expected %d, got %d",
			nelems, len(listVal.Values))
	}

	return listVal
}

func checkMap(t *testing.T, val data.TorinoValue, nelems int) *data.TorinoMap {
	mapVal, ok := val.(*data.TorinoMap)
	if !ok {
		t.Fatalf("Wrong Torino type: expected *TorinoMap, got %T", val)
	}

	if len(mapVal.Values) != nelems {
		t.Fatalf("Wrong number of map elements: expected %d, got %d",
			nelems, len(mapVal.Values))
	}

	return mapVal
}
