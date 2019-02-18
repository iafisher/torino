package compiler

import "github.com/iafisher/torino/data"

type Instruction struct {
	Name string
	Args []data.TorinoValue
}

func NewInst(name string, args ...data.TorinoValue) *Instruction {
	return &Instruction{name, args}
}
