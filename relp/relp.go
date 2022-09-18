package relp

import (
	"bufio"
	"fmt"
	"io"

	"github.com/karamaru-alpha/monkey/lexer"
	"github.com/karamaru-alpha/monkey/parser"
)

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

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
		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}
}
