package relp

import (
	"bufio"
	"fmt"
	"io"

	"github.com/karamaru-alpha/monkey/compiler"
	"github.com/karamaru-alpha/monkey/lexer"
	"github.com/karamaru-alpha/monkey/parser"
	"github.com/karamaru-alpha/monkey/vm"
)

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	fmt.Println("console...")
	for {
		fmt.Print(">> ")
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		input := scanner.Text()
		if input == "exit" {
			fmt.Println("bye!")
			return
		}
		p := parser.New(lexer.New(input))
		program := p.ParseProgram()
		if len(p.Errors()) > 0 {
			for _, msg := range p.Errors() {
				io.WriteString(out, msg)
			}
			return
		}

		// Compiler

		comp := compiler.New()
		if err := comp.Compile(program); err != nil {
			fmt.Fprintf(out, "compile failed: \n %s\n", err)
			continue
		}

		machine := vm.New(comp.Bytecode())
		if err := machine.Run(); err != nil {
			fmt.Fprintf(out, "executing bytecode failed: \n %s\n", err)
			continue
		}

		stackTop := machine.StackTop()
		io.WriteString(out, stackTop.Inspect())
		io.WriteString(out, "\n")

		// Interpreter

		//environment := object.NewEnvironment()
		//evaluated := evaluator.Eval(program, environment)
		//if evaluated != nil {
		//	io.WriteString(out, evaluated.Inspect())
		//	io.WriteString(out, "\n")
		//}
	}
}
