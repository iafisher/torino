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
	stack []data.TorinoValue
}

func New() *VirtualMachine {
	return &VirtualMachine{}
}

func (vm *VirtualMachine) Execute(program []*compiler.Instruction, env *Environment) data.TorinoValue {
	/* Print the bytecode, for debugging.
	for _, inst := range program {
		fmt.Printf("%v\n", inst)
	}
	fmt.Println("DONE")
	*/

	pc := 0
	for pc < len(program) {
		inst := program[pc]
		vm.executeOne(inst, env)

		// Update the program counter.
		if inst.Name == "REL_JUMP_IF_FALSE" {
			cond := vm.popStack().(*data.TorinoBool)
			if !cond.Value {
				pc += int(inst.Args[0].(*data.TorinoInt).Value)
			} else {
				pc += 1
			}
		} else if inst.Name == "REL_JUMP" {
			pc += int(inst.Args[0].(*data.TorinoInt).Value)
		} else if inst.Name == "RETURN_VALUE" {
			break
		} else {
			pc += 1
		}
	}

	if len(vm.stack) > 0 {
		return vm.stack[len(vm.stack)-1]
	} else {
		return &data.TorinoNone{}
	}
}

func (vm *VirtualMachine) executeOne(inst *compiler.Instruction, env *Environment) {
	if inst.Name == "PUSH_CONST" {
		vm.pushStack(inst.Args[0])
	} else if inst.Name == "STORE_NAME" {
		key := inst.Args[0].(*data.TorinoString).Value
		_, ok := env.Get(key)
		if ok {
			panic(fmt.Sprintf("cannot redefine symbol %s", key))
		}
		env.Put(key, vm.popStack())
	} else if inst.Name == "ASSIGN_NAME" {
		key := inst.Args[0].(*data.TorinoString).Value
		_, ok := env.Get(key)
		if !ok {
			panic(fmt.Sprintf("undefined symbol %s", key))
		}
		env.Put(key, vm.popStack())
	} else if inst.Name == "PUSH_NAME" {
		key := inst.Args[0].(*data.TorinoString).Value
		val, ok := env.Get(key)
		if !ok {
			panic(fmt.Sprintf("undefined symbol %s", key))
		}
		vm.pushStack(val)
	} else if inst.Name == "BINARY_ADD" {
		left, right := vm.popTwoInts()
		vm.pushStack(&data.TorinoInt{left.Value + right.Value})
	} else if inst.Name == "BINARY_SUB" {
		left, right := vm.popTwoInts()
		vm.pushStack(&data.TorinoInt{left.Value - right.Value})
	} else if inst.Name == "BINARY_MUL" {
		left, right := vm.popTwoInts()
		vm.pushStack(&data.TorinoInt{left.Value * right.Value})
	} else if inst.Name == "BINARY_DIV" {
		left, right := vm.popTwoInts()
		vm.pushStack(&data.TorinoInt{left.Value / right.Value})
	} else if inst.Name == "BINARY_EQ" {
		left, right := vm.popTwoInts()
		vm.pushStack(&data.TorinoBool{left.Value == right.Value})
	} else if inst.Name == "BINARY_GT" {
		left, right := vm.popTwoInts()
		vm.pushStack(&data.TorinoBool{left.Value > right.Value})
	} else if inst.Name == "BINARY_LT" {
		left, right := vm.popTwoInts()
		vm.pushStack(&data.TorinoBool{left.Value < right.Value})
	} else if inst.Name == "BINARY_GE" {
		left, right := vm.popTwoInts()
		vm.pushStack(&data.TorinoBool{left.Value >= right.Value})
	} else if inst.Name == "BINARY_LE" {
		left, right := vm.popTwoInts()
		vm.pushStack(&data.TorinoBool{left.Value <= right.Value})
	} else if inst.Name == "BINARY_AND" {
		left, right := vm.popTwoBools()
		vm.pushStack(&data.TorinoBool{left.Value && right.Value})
	} else if inst.Name == "BINARY_OR" {
		left, right := vm.popTwoBools()
		vm.pushStack(&data.TorinoBool{left.Value || right.Value})
	} else if inst.Name == "UNARY_MINUS" {
		arg := vm.popStack().(*data.TorinoInt)
		vm.pushStack(&data.TorinoInt{-arg.Value})
	} else if inst.Name == "CALL_FUNCTION" {
		// Get the function itself.
		tos := vm.popStack()

		// Gather the arguments for the function.
		args := []data.TorinoValue{}
		for i := 0; i < inst.Args[0].(*data.TorinoInt).Value; i++ {
			args = append(args, vm.popStack())
		}

		builtin, ok := tos.(*data.TorinoBuiltin)
		if ok {
			vm.pushStack(builtin.F(args...))
		} else {
			f := tos.(*compiler.TorinoFunction)

			fEnv := NewEnv(env)

			if len(args) != len(f.Params) {
				panic("Too few arguments to user-defined function")
			}

			for i, param := range f.Params {
				fEnv.Put(param.Value, args[i])
			}

			vm.pushStack(vm.Execute(f.Body, fEnv))
		}
		// The following operations don't affect the stack.
	} else if inst.Name == "RETURN_VALUE" {
	} else if inst.Name == "REL_JUMP_IF_FALSE" {
	} else if inst.Name == "REL_JUMP" {
	} else {
		panic(fmt.Sprintf("VirtualMachine.Execute - unknown instruction %s", inst.Name))
	}
}

func (vm *VirtualMachine) pushStack(vals ...data.TorinoValue) {
	vm.stack = append(vm.stack, vals...)
}

func (vm *VirtualMachine) popStack() data.TorinoValue {
	ret := vm.stack[len(vm.stack)-1]
	vm.stack = vm.stack[:len(vm.stack)-1]
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
