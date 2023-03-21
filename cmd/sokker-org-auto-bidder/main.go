package main

import (
	"log"
	"os"
	"sokker-org-auto-bidder/internal/repository/player"
	"sokker-org-auto-bidder/internal/subcommands"

	_ "github.com/mattn/go-sqlite3"
)

var playerRepository player.PlayerRepository

func main() {
	// check subcommand is provided
	if len(os.Args) < 2 {
		wrongSubcommand()
	}

	// handle repository
	playerRepository = createPlayerRepository()
	defer playerRepository.Close()

	// get subcommand args
	args := os.Args[2:]

	// choose subcommand to run
	switch os.Args[1] {
	case "bid":
		bidCmd := subcommands.BidSubcommand{R: playerRepository, Args: args}
		if err := bidCmd.Run(); err != nil {
			log.Fatal(err)
		}
	case "add":
		addCmd := subcommands.PlayerAddSubcommand{R: playerRepository, Args: args}
		if err := addCmd.Run(); err != nil {
			log.Fatal(err)
		}
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

func wrongSubcommand() {
	log.Fatal("expected 'bid' or 'add' subcommand")
}