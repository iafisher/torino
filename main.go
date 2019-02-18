package main

import (
	"bufio"
	"fmt"
	"github.com/iafisher/torino/compiler"
	"github.com/iafisher/torino/data"
	"github.com/iafisher/torino/lexer"
	"github.com/iafisher/torino/parser"
	"github.com/iafisher/torino/vm"
	"os"
)

func main() {
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
		p := parser.New(lexer.New(line))
		ast := p.Parse()

		cmp := compiler.New()
		program := cmp.Compile(ast)

		val := vm.Execute(program, env)
		_, isNone := val.(*data.TorinoNone)
		if !isNone {
			fmt.Println(val.String())
		}
	}
}
