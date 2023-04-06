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

	log.Trace("create player repository")
	playerRepository := createPlayerRepository()
	defer playerRepository.Close()

	log.Trace("create subcommands registry, register subcommands")
	subCmdRegistry := subcommands.NewSubcommandRegistry()
	subCmdRegistry.Register("bid", subcommands.NewBidSubcommand(playerRepository, client))
	subCmdRegistry.Register("add", subcommands.NewPlayerAddSubcommand(playerRepository, client))
	subCmdRegistry.Register("check-auth", subcommands.NewCheckAuthSubcommand(client))

	log.Trace("check is subcommand provided")
	if len(os.Args) < 2 {
		printError(fmt.Sprintf("No subcommand provided. Try one of %v", subCmdRegistry.GetSubcommandNames()))
	}

	subcommand := os.Args[1]
	log.Tracef("%s subcommand chosen", subcommand)
	args := os.Args[2:]
	log.Tracef("subcommand args: %v", args)

	log.Trace("execute subcommand")
	if err := subCmdRegistry.Run(subcommand, args); err != nil {
		log.Error(err)
		printError("Command execution failed. Run with -v flag for more information")
	}
}

// createPlayerRepository returns new player.PlayerRepository instance
func createPlayerRepository() player.PlayerRepository {
	playerRepository := player.NewSqlitePlayerRepository("./bidder.db")
	if err := playerRepository.Init(); err != nil {
		log.Error(err)
		printError("Could not open player bid list database. Run with -v flag for more information")
	}

	return playerRepository
}

// printError logs error message and exits program
func printError(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
