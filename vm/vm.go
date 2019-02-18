/* The Torino virtual machine.

Author:  Ian Fisher (iafisher@protonmail.com)
Version: February 2019
*/
package vm

import "github.com/iafisher/torino/compiler"

type VirtualMachine struct {
}

func New() *VirtualMachine {
	return &VirtualMachine{}
}

func (vm *VirtualMachine) Execute(program []compiler.Instruction) {
	for _, inst := range program {
		inst.Args()
	}
}
