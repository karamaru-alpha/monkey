package main

import (
	"fmt"
	"os"

	"github.com/karamaru-alpha/monkey/relp"
)

func main() {
	fmt.Println("Hello! This is the Monkey console.")
	relp.Start(os.Stdin, os.Stdout)
}
