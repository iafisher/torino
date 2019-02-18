/* The Torino virtual machine.

Author:  Ian Fisher (iafisher@protonmail.com)
Version: February 2019
*/
package vm

import (
	"fmt"
	"github.com/iafisher/torino/compiler"
	"github.com/iafisher/torino/data"
)

type VirtualMachine struct {
	Stack []data.TorinoValue
}

func New() *VirtualMachine {
	return &VirtualMachine{}
}

func (vm *VirtualMachine) Execute(
	program []*compiler.Instruction, env *Environment,
) data.TorinoValue {
	for _, inst := range program {
		vm.executeOne(inst, env)
	}

	if len(vm.Stack) > 0 {
		return vm.Stack[len(vm.Stack)-1]
	} else {
		return &data.TorinoNone{}
	}
}

func (vm *VirtualMachine) executeOne(inst *compiler.Instruction, env *Environment) {
	if inst.Name == "PUSH_CONST" {
		vm.Stack = append(vm.Stack, inst.Args[0])
	} else if inst.Name == "STORE_NAME" {
		key := inst.Args[0].(*data.TorinoString).Value
		env.Put(key, vm.popStack())
	} else if inst.Name == "PUSH_NAME" {
		key := inst.Args[0].(*data.TorinoString).Value
		vm.Stack = append(vm.Stack, env.Get(key))
	} else if inst.Name == "ADD" {
		left, right := vm.popTwoInts()
		vm.Stack = append(vm.Stack, &data.TorinoInt{left.Value + right.Value})
	} else if inst.Name == "SUB" {
		left, right := vm.popTwoInts()
		vm.Stack = append(vm.Stack, &data.TorinoInt{left.Value - right.Value})
	} else if inst.Name == "MUL" {
		left, right := vm.popTwoInts()
		vm.Stack = append(vm.Stack, &data.TorinoInt{left.Value * right.Value})
	} else if inst.Name == "DIV" {
		left, right := vm.popTwoInts()
		vm.Stack = append(vm.Stack, &data.TorinoInt{left.Value / right.Value})
	} else {
		panic(fmt.Sprintf("VirtualMachine.Execute - unknown instruction %s", inst.Name))
	}
}

func (vm *VirtualMachine) popStack() data.TorinoValue {
	ret := vm.Stack[len(vm.Stack)-1]
	vm.Stack = vm.Stack[:len(vm.Stack)-1]
	return ret
}

func (vm *VirtualMachine) popTwoInts() (*data.TorinoInt, *data.TorinoInt) {
	leftInt := vm.popStack().(*data.TorinoInt)
	rightInt := vm.popStack().(*data.TorinoInt)
	return leftInt, rightInt
}
