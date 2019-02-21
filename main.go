package main

import (
	"bufio"
	"fmt"
	"github.com/iafisher/torino/data"
	"github.com/iafisher/torino/eval"
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
	env := vm.NewEnv(nil)
	for {
		fmt.Print(">>> ")
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		oneline(line, env)
	}
}

func oneline(text string, env *vm.Environment) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("Error:", err)
		}
	}()

	val, err := eval.Eval(text, env)
	if err != nil {
		fmt.Println("Error:", err)
	}

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

	env := vm.NewEnv(nil)
	_, err = eval.Eval(text, env)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
