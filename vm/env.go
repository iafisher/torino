package vm

import "github.com/iafisher/torino/data"

type Environment struct {
	symbols map[string]data.TorinoValue
}

func NewEnv() *Environment {
	env := &Environment{map[string]data.TorinoValue{}}
	env.Put("print", &data.TorinoBuiltin{builtinPrint})
	return env
}

func (env *Environment) Get(k string) (data.TorinoValue, bool) {
	// TODO: Could this be a one-liner?
	val, ok := env.symbols[k]
	return val, ok
}

func (env *Environment) Put(k string, v data.TorinoValue) {
	env.symbols[k] = v
}
