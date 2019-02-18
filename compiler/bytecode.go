package compiler

type Instruction interface {
	Name() string
	Args() []int
}
