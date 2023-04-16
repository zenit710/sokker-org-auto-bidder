package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"sokker-org-auto-bidder/internal/client"
	"sokker-org-auto-bidder/internal/repository/player"
	"sokker-org-auto-bidder/internal/repository/session"
	"sokker-org-auto-bidder/internal/subcommands"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

const dbPath = "./bidder.db"

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
	log.SetLevel(getLogLevel())

	log.Trace("open db connection")
	db, err := openDbConnection(dbPath)
	if err != nil {
		log.Error(err)
		printError("Can not open internal database. Run with -v flag for more information")
	}
	defer db.Close()

	log.Trace("create player repository")
	playerRepository := createPlayerRepository(db)

	log.Trace("create session repository")
	sessionRepository := createSessionRepository(db)

	log.Trace("create new http client instance")
	var client client.Client = client.NewHttpClient(
		os.Getenv("SOKKER_USER"),
		os.Getenv("SOKKER_PASS"),
		sessionRepository,
	)

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
		switch err.(type) {
		case *subcommands.ErrMissingFlags: printError(fmt.Sprintf("Bad input: %v", err))
		default:
			log.Error(err)
			printError("Command execution failed. Run with -v flag for more information")
		}
	}
}

// openDbConnection opens sqlite db connection
func openDbConnection(path string) (*sql.DB, error) {
	log.Tracef("opening sqlite3 connection for: %s", path)
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Error(err)
		return nil, fmt.Errorf("could not open sqlite db %s", path)
	}

	return db, nil
}

// createPlayerRepository returns new initialized player.PlayerRepository instance
func createPlayerRepository(db *sql.DB) player.PlayerRepository {
	repo := player.NewSqlitePlayerRepository(db)
	if err := repo.Init(); err != nil {
		log.Error(err)
		printError("Could not init player bid list database. Run with -v flag for more information")
	}

	return repo
}

// createSessionRepository returns new initialized session.SessionRepository instance
func createSessionRepository(db *sql.DB) session.SessionRepository {
	repo := session.NewSqliteSessionRepository(db)
	if err := repo.Init(); err != nil {
		log.Error(err)
		printError("Could not init sessions database. Run with -v flag for more information")
	}

	return repo
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
