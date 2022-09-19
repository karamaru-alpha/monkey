package relp

import (
	"bufio"
	"fmt"
	"io"

	"github.com/karamaru-alpha/monkey/evaluator"
	"github.com/karamaru-alpha/monkey/lexer"
	"github.com/karamaru-alpha/monkey/object"
	"github.com/karamaru-alpha/monkey/parser"
)

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	environment := object.NewEnvironment()

	for {
		fmt.Print(">> ")
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "exit" {
			fmt.Println("bye!")
			return
		}
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) > 0 {
			for _, msg := range p.Errors() {
				io.WriteString(out, msg)
			}
			return
		}
		evaluated := evaluator.Eval(program, environment)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}
