package main

import (
	"flag"
	"fmt"
	"os"
	"sokker-org-auto-bidder/internal/client"
	"sokker-org-auto-bidder/internal/repository/player"
	"sokker-org-auto-bidder/internal/subcommands"

	log "github.com/sirupsen/logrus"
)

var logTraceLvl, logDebugLvl, logWarnLvl bool

// init initiaties module before run
func init() {
	flag.BoolVar(&logWarnLvl, "v", false, "show logs up to warning level")
	flag.BoolVar(&logDebugLvl, "vv", false, "show logs up to debug level")
	flag.BoolVar(&logTraceLvl, "vvv", false, "show all logs including trace messages")
}

// main is a central point of application
func main() {
	flag.Parse()
	args := flag.Args()

	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)
	log.SetLevel(getLogLevel()) // TODO: -v Warning, -vv Debug, -vvv Trace

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
	if len(args) == 0 {
		printError(fmt.Sprintf("No subcommand provided. Try one of %v", subCmdRegistry.GetSubcommandNames()))
	}

	subcommand := args[0]
	log.Infof("%s subcommand chosen", subcommand)
	subcommandArgs := args[1:]
	log.Infof("subcommand args: %v", subcommandArgs)

	log.Trace("execute subcommand")
	if err := subCmdRegistry.Run(subcommand, subcommandArgs); err != nil {
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

// getLogLevel returns log level on command flags base
func getLogLevel() log.Level {
	if logTraceLvl {return log.TraceLevel}
	if logDebugLvl {return log.DebugLevel}
	if logWarnLvl {return log.WarnLevel}
	return log.PanicLevel
}
