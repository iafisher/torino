package compiler

import "github.com/iafisher/torino/data"

type Instruction struct {
	Name string
	Args []data.TorinoValue
}
