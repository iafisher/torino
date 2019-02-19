package main

import (
	"bufio"
	"fmt"
	"github.com/iafisher/torino/compiler"
	"github.com/iafisher/torino/data"
	"github.com/iafisher/torino/lexer"
	"github.com/iafisher/torino/parser"
	"github.com/iafisher/torino/vm"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		repl()
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		fmt.Println("Error: too many command-line arguments supplied.")
	}
}

func repl() {
	fmt.Println("The Torino programming language.\n")

	scanner := bufio.NewScanner(os.Stdin)
	env := vm.NewEnv()
	vm := vm.New()
	for {
		fmt.Print(">>> ")
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		oneline(line, vm, env)
	}
}

func oneline(text string, vm *vm.VirtualMachine, env *vm.Environment) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("Error:", err)
		}
	}()

	p := parser.New(lexer.New(text))
	ast := p.Parse()

	cmp := compiler.New()
	program := cmp.Compile(ast)

	val := vm.Execute(program, env)
	_, isNone := val.(*data.TorinoNone)
	if !isNone {
		fmt.Println(val.Repr())
	}
}

func runFile(path string) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	text := string(contents)

	p := parser.New(lexer.New(text))
	ast := p.Parse()

	cmp := compiler.New()
	program := cmp.Compile(ast)

	env := vm.NewEnv()
	vm := vm.New()
	vm.Execute(program, env)
}
