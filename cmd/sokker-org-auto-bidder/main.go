package main

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"sokker-org-auto-bidder/internal/model"
	"sokker-org-auto-bidder/internal/repository"

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
	r := repository.NewDbRepository(db)
	if err := r.Init(); err != nil {
		log.Print(err.Error())
		return
	}

	// define 'add' subcommand flags set
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	playerId := addCmd.Int("playerId", 0, "Player ID")
	maxPrice := addCmd.Int("maxPrice", 0, "Maxium price for player to bid")

	// choose subcommand to run
	switch os.Args[1] {
	case "bid":
		log.Print("make bid for listed players:")

		// get players to bid list
		players, err := r.GetList()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		// print all players to bid
		for _, player := range players {
			log.Printf("%v", player)
		}
	case "add":
		// parse cmd flags
		addCmd.Parse(os.Args[2:])

		// validate playerId
		if *playerId <= 0 {
			log.Fatal("playerId has to be greater than zero")
		}

		// validate maxPrice
		if *maxPrice <= 0 {
			log.Fatal("maxPrice has to be greater than zero")
		}

		// create player model
		player := &model.Player{
			Id: uint(*playerId),
			MaxPrice: uint(*maxPrice),
		}

		// save player into the DB
		if err := r.Add(player); err != nil {
			log.Fatal(err)
		}

		log.Printf("player added to bid list: %v", player)
	default:
		wrongSubcommand()
	}
}

func wrongSubcommand() {
	log.Fatal("expected 'bid' or 'add' subcommand")
}