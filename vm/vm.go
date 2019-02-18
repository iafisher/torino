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
		_, ok := env.Get(key)
		if ok {
			panic(fmt.Sprintf("cannot redefine symbol %s", key))
		}
		env.Put(key, vm.popStack())
	} else if inst.Name == "PUSH_NAME" {
		key := inst.Args[0].(*data.TorinoString).Value
		val, ok := env.Get(key)
		if !ok {
			panic(fmt.Sprintf("undefined symbol %s", key))
		}
		vm.Stack = append(vm.Stack, val)
	} else if inst.Name == "BINARY_ADD" {
		left, right := vm.popTwoInts()
		vm.Stack = append(vm.Stack, &data.TorinoInt{left.Value + right.Value})
	} else if inst.Name == "BINARY_SUB" {
		left, right := vm.popTwoInts()
		vm.Stack = append(vm.Stack, &data.TorinoInt{left.Value - right.Value})
	} else if inst.Name == "BINARY_MUL" {
		left, right := vm.popTwoInts()
		vm.Stack = append(vm.Stack, &data.TorinoInt{left.Value * right.Value})
	} else if inst.Name == "BINARY_DIV" {
		left, right := vm.popTwoInts()
		vm.Stack = append(vm.Stack, &data.TorinoInt{left.Value / right.Value})
	} else if inst.Name == "BINARY_EQ" {
		left, right := vm.popTwoInts()
		vm.Stack = append(vm.Stack, &data.TorinoBool{left.Value == right.Value})
	} else if inst.Name == "BINARY_GT" {
		left, right := vm.popTwoInts()
		vm.Stack = append(vm.Stack, &data.TorinoBool{left.Value > right.Value})
	} else if inst.Name == "BINARY_LT" {
		left, right := vm.popTwoInts()
		vm.Stack = append(vm.Stack, &data.TorinoBool{left.Value < right.Value})
	} else if inst.Name == "BINARY_GE" {
		left, right := vm.popTwoInts()
		vm.Stack = append(vm.Stack, &data.TorinoBool{left.Value >= right.Value})
	} else if inst.Name == "BINARY_LE" {
		left, right := vm.popTwoInts()
		vm.Stack = append(vm.Stack, &data.TorinoBool{left.Value <= right.Value})
	} else if inst.Name == "BINARY_AND" {
		left, right := vm.popTwoBools()
		vm.Stack = append(vm.Stack, &data.TorinoBool{left.Value && right.Value})
	} else if inst.Name == "BINARY_OR" {
		left, right := vm.popTwoBools()
		vm.Stack = append(vm.Stack, &data.TorinoBool{left.Value || right.Value})
	} else if inst.Name == "UNARY_MINUS" {
		arg := vm.popStack().(*data.TorinoInt)
		vm.Stack = append(vm.Stack, &data.TorinoInt{-arg.Value})
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
	left := vm.popStack().(*data.TorinoInt)
	right := vm.popStack().(*data.TorinoInt)
	return left, right
}

func (vm *VirtualMachine) popTwoBools() (*data.TorinoBool, *data.TorinoBool) {
	left := vm.popStack().(*data.TorinoBool)
	right := vm.popStack().(*data.TorinoBool)
	return left, right
}
