/* The Torino virtual machine.

Author:  Ian Fisher (iafisher@protonmail.com)
Version: February 2019
*/
package vm

import (
	"errors"
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

func (vm *VirtualMachine) Execute(
	program []*compiler.Instruction, env *Environment) (data.TorinoValue, error) {
	/* Print the bytecode, for debugging.
	for _, inst := range program {
		fmt.Printf("%v\n", inst)
	}
	fmt.Println("DONE")
	*/

	pc := 0
	for pc < len(program) {
		inst := program[pc]
		err := vm.executeOne(inst, env)
		if err != nil {
			return nil, err
		}

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
		return vm.stack[len(vm.stack)-1], nil
	} else {
		return &data.TorinoNone{}, nil
	}
}

func (vm *VirtualMachine) executeOne(inst *compiler.Instruction, env *Environment) error {
	if inst.Name == "PUSH_CONST" {
		vm.pushStack(inst.Args[0])
	} else if inst.Name == "STORE_NAME" {
		key := inst.Args[0].(*data.TorinoString).Value
		_, ok := env.Get(key)
		if ok {
			return errors.New(fmt.Sprintf("cannot redefine symbol %s", key))
		}
		env.Put(key, vm.popStack())
	} else if inst.Name == "ASSIGN_NAME" {
		key := inst.Args[0].(*data.TorinoString).Value
		_, ok := env.Get(key)
		if !ok {
			return errors.New(fmt.Sprintf("undefined symbol %s", key))
		}
		env.Put(key, vm.popStack())
	} else if inst.Name == "PUSH_NAME" {
		key := inst.Args[0].(*data.TorinoString).Value
		val, ok := env.Get(key)
		if !ok {
			return errors.New(fmt.Sprintf("undefined symbol %s", key))
		}
		vm.pushStack(val)
	} else if inst.Name == "BINARY_ADD" {
		left, right, ok := vm.popTwoInts()
		if !ok {
			return errors.New("+ takes integer operands")
		}
		vm.pushStack(&data.TorinoInt{left.Value + right.Value})
	} else if inst.Name == "BINARY_SUB" {
		left, right, ok := vm.popTwoInts()
		if !ok {
			return errors.New("- takes integer operands")
		}
		vm.pushStack(&data.TorinoInt{left.Value - right.Value})
	} else if inst.Name == "BINARY_MUL" {
		left, right, ok := vm.popTwoInts()
		if !ok {
			return errors.New("* takes integer operands")
		}
		vm.pushStack(&data.TorinoInt{left.Value * right.Value})
	} else if inst.Name == "BINARY_DIV" {
		left, right, ok := vm.popTwoInts()
		if !ok {
			return errors.New("/ takes integer operands")
		}
		vm.pushStack(&data.TorinoInt{left.Value / right.Value})
	} else if inst.Name == "BINARY_EQ" {
		left, right, ok := vm.popTwoInts()
		if !ok {
			return errors.New("== takes integer operands")
		}
		vm.pushStack(&data.TorinoBool{left.Value == right.Value})
	} else if inst.Name == "BINARY_GT" {
		left, right, ok := vm.popTwoInts()
		if !ok {
			return errors.New("> takes integer operands")
		}
		vm.pushStack(&data.TorinoBool{left.Value > right.Value})
	} else if inst.Name == "BINARY_LT" {
		left, right, ok := vm.popTwoInts()
		if !ok {
			return errors.New("< takes integer operands")
		}
		vm.pushStack(&data.TorinoBool{left.Value < right.Value})
	} else if inst.Name == "BINARY_GE" {
		left, right, ok := vm.popTwoInts()
		if !ok {
			return errors.New(">= takes integer operands")
		}
		vm.pushStack(&data.TorinoBool{left.Value >= right.Value})
	} else if inst.Name == "BINARY_LE" {
		left, right, ok := vm.popTwoInts()
		if !ok {
			return errors.New("<= takes integer operands")
		}
		vm.pushStack(&data.TorinoBool{left.Value <= right.Value})
	} else if inst.Name == "BINARY_AND" {
		left, right, ok := vm.popTwoBools()
		if !ok {
			return errors.New("and takes boolean operands")
		}
		vm.pushStack(&data.TorinoBool{left.Value && right.Value})
	} else if inst.Name == "BINARY_OR" {
		left, right, ok := vm.popTwoBools()
		if !ok {
			return errors.New("or takes boolean operands")
		}
		vm.pushStack(&data.TorinoBool{left.Value || right.Value})
	} else if inst.Name == "BINARY_INDEX" {
		switch indexed := vm.popStack().(type) {
		case *data.TorinoList:
			index, ok := vm.popStack().(*data.TorinoInt)
			if !ok {
				return errors.New("index must be an integer")
			}

			if index.Value < 0 || index.Value >= len(indexed.Values) {
				return errors.New("index out of bounds")
			}
			vm.pushStack(indexed.Values[index.Value])
		case *data.TorinoMap:
			index := vm.popStack()
			val, ok := indexed.Get(index)
			if !ok {
				return errors.New("key not in map")
			}

			vm.pushStack(val)
		case *data.TorinoString:
			index, ok := vm.popStack().(*data.TorinoInt)
			if !ok {
				return errors.New("index must be an integer")
			}

			if index.Value < 0 || index.Value >= len(indexed.Value) {
				return errors.New("index out of bounds")
			}

			vm.pushStack(&data.TorinoString{string(indexed.Value[index.Value])})
		default:
			return errors.New("only lists and maps may be indexed")
		}
	} else if inst.Name == "UNARY_MINUS" {
		arg, ok := vm.popStack().(*data.TorinoInt)
		if !ok {
			return errors.New("unary - takes integer operand")
		}

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
			res, err := builtin.F(args...)
			if err != nil {
				return err
			}
			vm.pushStack(res)
		} else {
			f, ok := tos.(*compiler.TorinoFunction)
			if !ok {
				return errors.New("cannot apply non-function")
			}

			fEnv := NewEnv(env)

			if len(args) != len(f.Params) {
				return errors.New("wrong number of arguments to user-defined function")
			}

			for i, param := range f.Params {
				fEnv.Put(param.Value, args[i])
			}

			val, err := vm.Execute(f.Body, fEnv)
			if err != nil {
				return err
			}
			vm.pushStack(val)
		}
		// The following operations don't affect the stack.
	} else if inst.Name == "MAKE_LIST" {
		nelems := inst.Args[0].(*data.TorinoInt).Value

		values := []data.TorinoValue{}
		for i := 0; i < nelems; i++ {
			values = append(values, vm.popStack())
		}
		vm.pushStack(&data.TorinoList{values})
	} else if inst.Name == "MAKE_MAP" {
		nelems := inst.Args[0].(*data.TorinoInt).Value

		mapVal := &data.TorinoMap{map[string]data.TorinoValue{}}
		for i := 0; i < nelems; i++ {
			val := vm.popStack()
			key := vm.popStack()
			mapVal.Put(key, val)
		}
		vm.pushStack(mapVal)
	} else if inst.Name == "RETURN_VALUE" {
	} else if inst.Name == "REL_JUMP_IF_FALSE" {
	} else if inst.Name == "REL_JUMP" {
	} else {
		return errors.New(fmt.Sprintf("unknown instruction %s", inst.Name))
	}

	return nil
}

func (vm *VirtualMachine) pushStack(vals ...data.TorinoValue) {
	vm.stack = append(vm.stack, vals...)
}

func (vm *VirtualMachine) popStack() data.TorinoValue {
	ret := vm.stack[len(vm.stack)-1]
	vm.stack = vm.stack[:len(vm.stack)-1]
	return ret
}

func (vm *VirtualMachine) popTwoInts() (*data.TorinoInt, *data.TorinoInt, bool) {
	left, ok1 := vm.popStack().(*data.TorinoInt)
	right, ok2 := vm.popStack().(*data.TorinoInt)
	return left, right, ok1 && ok2
}

func (vm *VirtualMachine) popTwoBools() (*data.TorinoBool, *data.TorinoBool, bool) {
	left, ok1 := vm.popStack().(*data.TorinoBool)
	right, ok2 := vm.popStack().(*data.TorinoBool)
	return left, right, ok1 && ok2
}
