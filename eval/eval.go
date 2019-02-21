/* A helper function to directly evaluate a string in an environment.

Author:  Ian Fisher (iafisher@protonmail.com)
Version: February 2019
*/
package eval

import (
	"errors"
	"github.com/iafisher/torino/compiler"
	"github.com/iafisher/torino/data"
	"github.com/iafisher/torino/lexer"
	"github.com/iafisher/torino/parser"
	"github.com/iafisher/torino/vm"
)

func Eval(text string, env *vm.Environment) (data.TorinoValue, error) {
	vm := vm.New()

	p := parser.New(lexer.New(text))
	ast, ok := p.Parse()
	if !ok {
		return nil, errors.New(p.Errors()[0])
	}

	cmp := compiler.New()
	program, err := cmp.Compile(ast)
	if err != nil {
		return nil, err
	}

	return vm.Execute(program, env), nil
}
