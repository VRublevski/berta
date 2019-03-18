package main

import (
	"fmt"
	"github.com/user/courseWork/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s from Berta programming language!\n", user.Username)
	fmt.Printf("You can type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}
