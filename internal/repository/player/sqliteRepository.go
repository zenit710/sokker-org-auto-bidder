package player

import (
	"database/sql"
	"fmt"
	"sokker-org-auto-bidder/internal/model"
	"sokker-org-auto-bidder/internal/repository"
	"time"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

// timeLayout store time.Layout used inside sqlite
const timeLayout = "2006-01-02 15:04:05"

var _ PlayerRepository = &sqlitePlayerRepository{}

// sqlitePlayerRepository handle sqlite connection for player bid list
type sqlitePlayerRepository struct {
	path string
	db *sql.DB
}

// NewSqlitePlayerRepository returns new repository with sqlite connection
func NewSqlitePlayerRepository(path string) *sqlitePlayerRepository {
	log.Trace("creating new sqlite player repository")
	return &sqlitePlayerRepository{path: path}
}

// OpenConnection opens sqlite db connection
func (r *sqlitePlayerRepository) OpenConnection() error {
	log.Tracef("opening sqlite3 connection for: %s", r.path)
	db, err := sql.Open("sqlite3", r.path)
	if err != nil {
		log.Error(err)
		return err
	}
	r.db = db 

	return nil
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
		return repository.NewErrRepositoryInitFailure(fmt.Sprintf("%q: %s", err.Error(), sqlStmt))
	}

	return nil
}

func (r *sqlitePlayerRepository) Init() error {
	log.Trace("sqlite player repository init")
	if err := r.OpenConnection(); err != nil {
		log.Error(err)
		return repository.NewErrRepositoryInitFailure(err.Error())
	}
	log.Debug("sqlite connection open")

	if err := r.CreateSchema(); err != nil {
		log.Error(err)
		return repository.NewErrRepositoryInitFailure(err.Error())
	}

	return nil
}

func (r *sqlitePlayerRepository) Add(player *model.Player) error {
	log.Tracef("add player (%v) to the bid list", player)

	return r.makeTransaction(
		`insert into players (playerId, maxPrice, deadline) values(?, ?, datetime(?))`,
		player.Id, player.MaxPrice, player.Deadline.Format(time.RFC3339),
	)
}

func (r *sqlitePlayerRepository) Delete(player *model.Player) error {
	log.Tracef("remove player (%d) from the bid list", player.Id)

	return r.makeTransaction(
		`delete from players where playerId = ?`,
		player.Id,
	)
}

func (r *sqlitePlayerRepository) List() ([]*model.Player, error) {
	log.Trace("fetch players on bid list")
	rows, err := r.db.Query(`select * from players where deadline > datetime("now")`)
	if err != nil {
		log.Error(err)
		return nil, err
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
			log.Error(err)
			return nil, err
		}

		log.Trace("parse deadline time")
		deadline, err := time.Parse(timeLayout, dt)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		player.Deadline = deadline
		players = append(players, player)
	}

	return players, nil
}

func (r *sqlitePlayerRepository) Update(player *model.Player) error {
	log.Tracef("update player (%d) on bid list", player.Id)

	return r.makeTransaction(
		`update players set maxPrice = ?, deadline = datetime(?) where playerId = ?`,
		player.MaxPrice, player.Deadline.Format(time.RFC3339), player.Id,
	)
}

func (r *sqlitePlayerRepository) Close() {
	log.Debug("close db connection")
	r.db.Close()
}

// makeTransaction handles DB transaction
func (r *sqlitePlayerRepository) makeTransaction(sql string, params ...interface{}) error {
	log.Trace("create db transaction")
	tx, err := r.db.Begin()
	if err != nil {
		log.Error(err)
		return err
	}

	log.Trace("prepare sql statement")
	stmt, err := tx.Prepare(sql)
	if err != nil {
		log.Error(err)
		return err
	}
	defer stmt.Close()
	
	log.Trace("execute transaction")
	_, err = stmt.Exec(params)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Trace("commit transaction")
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
