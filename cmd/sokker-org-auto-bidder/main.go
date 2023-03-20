package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// check subcommand is provided
	if len(os.Args) < 2 {
		wrongSubcommand()
	}

	// init db connections
	db, err := sql.Open("sqlite3", "./bidder.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// ensure db tables structure created
	sqlStmt := `create table if not exists players (playerId integer not null primary key, maxPrice integer not null);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	// define 'add' subcommand flags set
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	playerId := addCmd.Int("playerId", 0, "Player ID")
	maxPrice := addCmd.Int("maxPrice", 0, "Maxium price for player to bid")

	// choose subcommand to run
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