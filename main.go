package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		wrongSubcommand()
	}

	switch os.Args[1] {
	case "bid":
		fmt.Println("make bid for listed players")
	case "add":
		fmt.Println("add player to bid list")
	default:
		wrongSubcommand()
	}
}

func wrongSubcommand() {
	fmt.Println("expected 'bid' or 'add' subcommand")
	os.Exit(1)
}