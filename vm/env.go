package vm

import "github.com/iafisher/torino/data"

type Environment struct {
	symbols   map[string]data.TorinoValue
	enclosing *Environment
}

func NewEnv(enclosing *Environment) *Environment {
	env := &Environment{map[string]data.TorinoValue{}, enclosing}
	env.Put("print", &data.TorinoBuiltin{builtinPrint})
	env.Put("println", &data.TorinoBuiltin{builtinPrintln})
	return env
}

func (env *Environment) Get(k string) (data.TorinoValue, bool) {
	// TODO: Could this be a one-liner?
	val, ok := env.symbols[k]
	if !ok && env.enclosing != nil {
		return env.enclosing.Get(k)
	} else {
		return val, ok
	}
}

func (env *Environment) Put(k string, v data.TorinoValue) {
	env.symbols[k] = v
}
