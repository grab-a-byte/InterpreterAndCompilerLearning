package main

import (
	"fmt"
	"monkey/repl"
	"os"
	"os/user"
)

func main() {
	currUser, err := user.Current()
	if err != nil {
		panic("Unknown User")
	}

	fmt.Printf("Hello %s , Welcome to the Monkey Programming Language", currUser.Username)
	fmt.Println("Feel free to type commands in!")
	repl.Start(os.Stdin, os.Stdout)
}
