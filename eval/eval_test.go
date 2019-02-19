package eval

import (
	"github.com/iafisher/torino/data"
	"github.com/iafisher/torino/vm"
	"testing"
)

func TestEvalIfStatement(t *testing.T) {
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

// Helper functions

func evalHelper(text string) data.TorinoValue {
	env := vm.NewEnv()
	return Eval(text, env)
}

func checkInteger(t *testing.T, val data.TorinoValue, expected int64) {
	intVal, ok := val.(*data.TorinoInt)
	if !ok {
		t.Fatalf("Wrong Torino type: expected *TorinoInt, got %T", intVal)
	}

	if intVal.Value != expected {
		t.Fatalf("Wrong integer value: expected %d, got %d", expected, intVal.Value)
	}
}