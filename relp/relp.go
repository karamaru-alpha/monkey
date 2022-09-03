package relp

import (
	"bufio"
	"fmt"
	"io"

	"github.com/karamaru-alpha/monkey/lexer"
	"github.com/karamaru-alpha/monkey/token"
)

func Start(in io.Reader, _ io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(">> ")
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "exit" {
			return
		}
		l := lexer.New(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
