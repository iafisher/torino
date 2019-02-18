package main

import (
	"bufio"
	"fmt"
	"github.com/iafisher/torino/compiler"
	"github.com/iafisher/torino/lexer"
	"github.com/iafisher/torino/parser"
	"github.com/iafisher/torino/vm"
	"os"
)

func main() {
	fmt.Println("The Torino programming language.")

	scanner := bufio.NewScanner(os.Stdin)
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

		vm := vm.New()
		vm.Execute(program)
	}
}
