package player

import (
	"database/sql"
	"fmt"
	"sokker-org-auto-bidder/internal/model"
	"sokker-org-auto-bidder/internal/repository"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const timeLayout = "2006-01-02 15:04:05"

var _ PlayerRepository = &sqlitePlayerRepository{}

type sqlitePlayerRepository struct {
	path string
	db *sql.DB
}

func NewSqlitePlayerRepository(path string) *sqlitePlayerRepository {
	return &sqlitePlayerRepository{path: path}
}

func (r *sqlitePlayerRepository) OpenConnection() error {
	db, err := sql.Open("sqlite3", r.path)
	if err != nil {
		return err
	}
	r.db = db 

	return nil
}

func (r *sqlitePlayerRepository) CreateSchema() error {
	sqlStmt := `create table if not exists players (
		playerId integer not null primary key,
		maxPrice integer not null,
		deadline text not null
		);`

	if _, err := r.db.Exec(sqlStmt); err != nil {
		return repository.NewErrRepositoryInitFailure(fmt.Sprintf("%q: %s", err.Error(), sqlStmt))
	}

	return nil
}

func (r *sqlitePlayerRepository) Init() error {
	if err := r.OpenConnection(); err != nil {
		return repository.NewErrRepositoryInitFailure(err.Error())
	}

	if err := r.CreateSchema(); err != nil {
		return repository.NewErrRepositoryInitFailure(err.Error())
	}

	return nil
}

func (r *sqlitePlayerRepository) Add(player *model.Player) error {
	// start db transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// create player to bid insert statement
	stmt, err := tx.Prepare(`insert into players (playerId, maxPrice, deadline) values(?, ?, datetime(?))`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	
	// set transaction values
	_, err = stmt.Exec(player.Id, player.MaxPrice, player.Deadline.Format(time.RFC3339))
	if err != nil {
		return err
	}

	// commit db transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r *sqlitePlayerRepository) Delete(player *model.Player) error {
	// start db transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// delete player query
	stmt, err := tx.Prepare(`delete from players where playerId = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	
	// set transaction values
	_, err = stmt.Exec(player.Id)
	if err != nil {
		return err
	}

	// commit db transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r *sqlitePlayerRepository) List() ([]*model.Player, error) {
	// fetch players to bid from db
	rows, err := r.db.Query(`select * from players where deadline > datetime("now")`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := []*model.Player{}

	// convert db rows to player model
	for rows.Next() {
		player := &model.Player{}
		var dt string

		err = rows.Scan(&player.Id, &player.MaxPrice, &dt)
		if err != nil {
			return nil, err
		}

		deadline, err := time.Parse(timeLayout, dt)
		if err != nil {
			return nil, err
		}

		player.Deadline = deadline
		players = append(players, player)
	}

	return players, nil
}

func (r *sqlitePlayerRepository) Update(player *model.Player) error {
	// start db transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// update player query
	stmt, err := tx.Prepare(`update players set maxPrice = ?, deadline = datetime(?) where playerId = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	
	// set transaction values
	_, err = stmt.Exec(player.MaxPrice, player.Deadline.Format(time.RFC3339), player.Id)
	if err != nil {
		return err
	}

	// commit db transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r *sqlitePlayerRepository) Close() {
	r.db.Close()
}