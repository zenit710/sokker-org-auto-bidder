package player

import (
	"database/sql"
	"errors"
	"fmt"
	"sokker-org-auto-bidder/internal/model"
	"sokker-org-auto-bidder/internal/repository"
	"sokker-org-auto-bidder/tools"
	"time"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

// timeLayout store time.Layout used inside sqlite
const timeLayout = "2006-01-02 15:04:05"

var (
	_                     PlayerRepository = &sqlitePlayerRepository{}
	ErrCanNotFetchPlayers                  = errors.New("could not fetch players to bid")
)

// sqlitePlayerRepository handle sqlite connection for player bid list
type sqlitePlayerRepository struct {
	db *sql.DB
}

// NewSqlitePlayerRepository returns new repository with sqlite connection
func NewSqlitePlayerRepository(db *sql.DB) *sqlitePlayerRepository {
	log.Trace("creating new sqlite player repository")
	return &sqlitePlayerRepository{db}
}

// CreateSchema creates db structure if it is not created yet
func (r *sqlitePlayerRepository) CreateSchema() error {
	sqlStmt := `create table if not exists players (
		playerId integer not null primary key,
		maxPrice integer not null,
		deadline text not null
		);`

	log.Trace("create database schema if not exists")
	if _, err := r.db.Exec(sqlStmt); err != nil {
		log.Error(err)
		return errors.New("create schema sql execution failed")
	}

	return nil
}

func (r *sqlitePlayerRepository) Init() error {
	log.Trace("sqlite player repository init")
	if err := r.CreateSchema(); err != nil {
		log.Error(err)
		return repository.ErrCanNotCreateDbSchema
	}

	return nil
}

func (r *sqlitePlayerRepository) Add(player *model.Player) error {
	log.Tracef("add player (%v) to the bid list", player)

	return tools.MakeTransaction(
		r.db,
		`insert into players (playerId, maxPrice, deadline) values(?, ?, datetime(?))`,
		player.Id, player.MaxPrice, player.Deadline.Format(time.RFC3339),
	)
}

func (r *sqlitePlayerRepository) Delete(player *model.Player) error {
	log.Tracef("remove player (%d) from the bid list", player.Id)

	return tools.MakeTransaction(
		r.db,
		`delete from players where playerId = ?`,
		player.Id,
	)
}

func (r *sqlitePlayerRepository) List() ([]*model.Player, error) {
	log.Trace("fetch players on bid list")
	rows, err := r.db.Query(`select * from players where deadline > datetime("now")`)
	if err != nil {
		log.Error(err)
		return nil, ErrCanNotFetchPlayers
	}
	defer rows.Close()

	players := []*model.Player{}

	log.Trace("convert db entries to player models")
	for rows.Next() {
		player := &model.Player{}
		var dt string

		log.Trace("map entry to player model properties")
		err = rows.Scan(&player.Id, &player.MaxPrice, &dt)
		if err != nil {
			log.Warn(fmt.Errorf("could not map player from db into the player model: %v", err))
			continue
		}

		log.Trace("parse deadline time")
		deadline, err := time.Parse(timeLayout, dt)
		if err != nil {
			log.Warn(fmt.Errorf("could not parse deadline time for player (%d): %v", player.Id, err))
			continue
		}

		player.Deadline = deadline
		players = append(players, player)
	}

	return players, nil
}

func (r *sqlitePlayerRepository) Update(player *model.Player) error {
	log.Tracef("update player (%d) on bid list", player.Id)

	return tools.MakeTransaction(
		r.db,
		`update players set maxPrice = ?, deadline = datetime(?) where playerId = ?`,
		player.MaxPrice, player.Deadline.Format(time.RFC3339), player.Id,
	)
}
