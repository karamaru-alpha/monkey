package main

import (
	"os"

	"github.com/karamaru-alpha/monkey/relp"
)

func main() {
	relp.Start(os.Stdin, os.Stdout)
}
