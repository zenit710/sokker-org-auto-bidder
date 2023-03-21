package main

import (
	"flag"
	"log"
	"os"
	"sokker-org-auto-bidder/internal/repository/player"
	"sokker-org-auto-bidder/internal/subcommands"

	_ "github.com/mattn/go-sqlite3"
)

var playerRepository player.PlayerRepository

func main() {
	// create player repository
	playerRepository = createPlayerRepository()
	defer playerRepository.Close()

	// create subcommand registry
	subCmdRegistry := subcommands.NewSubcommandRegistry()
	subCmdRegistry.Register("bid", &subcommands.BidSubcommand{R: playerRepository})
	subCmdRegistry.Register("add", &subcommands.PlayerAddSubcommand{R: playerRepository})

	flag.Parse()

	// check subcommand provided
	if len(os.Args) < 2 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// get subcommand args
	subcommand := os.Args[1]
	args := os.Args[2:]

	// handle subcommand
	if err := subCmdRegistry.Run(subcommand, args); err != nil {
		log.Fatal(err)
	}
}

func createPlayerRepository() player.PlayerRepository {
	playerRepository := player.NewSqlitePlayerRepository("./bidder.db")
	if err := playerRepository.Init(); err != nil {
		log.Fatal(err.Error())
	}

	return playerRepository
}
