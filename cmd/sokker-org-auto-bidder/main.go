package main

import (
	"flag"
	"log"
	"os"
	"sokker-org-auto-bidder/internal/model"
	"sokker-org-auto-bidder/internal/repository/player"

	_ "github.com/mattn/go-sqlite3"
)

var playerRepository player.PlayerRepository

func main() {
	// check subcommand is provided
	if len(os.Args) < 2 {
		wrongSubcommand()
	}

	playerRepository = createPlayerRepository()
	defer playerRepository.Close()

	// choose subcommand to run
	switch os.Args[1] {
	case "bid":
		handleBidCommand()
	case "add":
		handleAddCommand()
	default:
		wrongSubcommand()
	}
}


func createPlayerRepository() player.PlayerRepository {
	playerRepository := player.NewSqlitePlayerRepository("./bidder.db")
	if err := playerRepository.Init(); err != nil {
		log.Fatal(err.Error())
	}

	return playerRepository
}

func handleAddCommand() {
	// define 'add' subcommand flags set
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	playerId := addCmd.Int("playerId", 0, "Player ID")
	maxPrice := addCmd.Int("maxPrice", 0, "Maxium price for player to bid")

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
	if err := playerRepository.Add(player); err != nil {
		log.Fatal(err)
	}

	log.Printf("player added to bid list: %v", player)
}

func handleBidCommand() {
	log.Print("make bid for listed players:")

	// get players to bid list
	players, err := playerRepository.GetList()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// print all players to bid
	for _, player := range players {
		log.Printf("%v", player)
	}
}

func wrongSubcommand() {
	log.Fatal("expected 'bid' or 'add' subcommand")
}