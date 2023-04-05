package main

import (
	"fmt"
	"os"
	"sokker-org-auto-bidder/internal/client"
	"sokker-org-auto-bidder/internal/repository/player"
	"sokker-org-auto-bidder/internal/subcommands"

	log "github.com/sirupsen/logrus"
)

// main is a central point of application
func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)

	log.Trace("create new http client instance")
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

// createPlayerRepository returns new player.PlayerRepository instance
func createPlayerRepository() player.PlayerRepository {
	playerRepository := player.NewSqlitePlayerRepository("./bidder.db")
	if err := playerRepository.Init(); err != nil {
		logError(err.Error())
	}

	return playerRepository
}

// logError logs error message and exits program
func logError(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
