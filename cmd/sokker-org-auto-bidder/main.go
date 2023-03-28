package main

import (
	"fmt"
	"os"
	"sokker-org-auto-bidder/internal/client"
	"sokker-org-auto-bidder/internal/repository/player"
	"sokker-org-auto-bidder/internal/subcommands"
)

func main() {
	// create client
	var client client.Client = client.NewHttpClient(os.Getenv("SOKKER_USER"), os.Getenv("SOKKER_PASS"))

	// create player repository
	playerRepository := createPlayerRepository()
	defer playerRepository.Close()

	// create subcommand registry
	subCmdRegistry := subcommands.NewSubcommandRegistry()
	subCmdRegistry.Register("bid", subcommands.NewBidSubcommand(playerRepository, client))
	subCmdRegistry.Register("add", subcommands.NewPlayerAddSubcommand(playerRepository, client))
	subCmdRegistry.Register("check-auth", subcommands.NewCheckAuthSubcommand(client))

	// check subcommand provided
	if len(os.Args) < 2 {
		logError(fmt.Sprintf("No subcommand provided. Try one of %v", subCmdRegistry.GetSubcommandNames()))
	}

	// get subcommand args
	subcommand := os.Args[1]
	args := os.Args[2:]

	// handle subcommand
	if err := subCmdRegistry.Run(subcommand, args); err != nil {
		logError(err.Error())
	}
}

func createPlayerRepository() player.PlayerRepository {
	playerRepository := player.NewSqlitePlayerRepository("./bidder.db")
	if err := playerRepository.Init(); err != nil {
		logError(err.Error())
	}

	return playerRepository
}

func logError(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
