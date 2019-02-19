/* A helper function to directly evaluate a string in an environment.

Author:  Ian Fisher (iafisher@protonmail.com)
Version: February 2019
*/
package eval

import (
	"github.com/iafisher/torino/compiler"
	"github.com/iafisher/torino/data"
	"github.com/iafisher/torino/lexer"
	"github.com/iafisher/torino/parser"
	"github.com/iafisher/torino/vm"
)

func Eval(text string, env *vm.Environment) data.TorinoValue {
	vm := vm.New()

	p := parser.New(lexer.New(text))
	ast := p.Parse()

	cmp := compiler.New()
	program := cmp.Compile(ast)

	return vm.Execute(program, env)
}
