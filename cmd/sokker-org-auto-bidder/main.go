package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		wrongSubcommand()
	}

	addCmd := flag.NewFlagSet("bid", flag.ExitOnError)
	playerId := addCmd.Int("playerId", 0, "Player ID")
	maxPrice := addCmd.Int("maxPrice", 0, "Maxium price for player to bid")

	switch os.Args[1] {
	case "bid":
		fmt.Println("make bid for listed players")
	case "add":
		addCmd.Parse(os.Args[2:])

		if *playerId <= 0 {
			fmt.Println("playerId has to be greater than zero")
			os.Exit(1)
		}

		if *maxPrice <= 0 {
			fmt.Println("maxPrice has to be greater than zero")
			os.Exit(1)
		}

		fmt.Println("add player to bid list")
		fmt.Println("id:", *playerId)
		fmt.Println("max price:", *maxPrice)
	default:
		wrongSubcommand()
	}
}

func wrongSubcommand() {
	fmt.Println("expected 'bid' or 'add' subcommand")
	os.Exit(1)
}