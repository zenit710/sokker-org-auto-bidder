package repository

import (
	"database/sql"
	"fmt"
	"sokker-org-auto-bidder/internal/model"

	_ "github.com/mattn/go-sqlite3"
)

type SqlRepository struct {
	db *sql.DB
}

func NewDbRepository(db *sql.DB) *SqlRepository {
	return &SqlRepository{db}
}

func (r *SqlRepository) Init() error {
	sqlStmt := `create table if not exists players (playerId integer not null primary key, maxPrice integer not null);`

	if _, err := r.db.Exec(sqlStmt); err != nil {
		return &ErrRepositoryInitFailure{fmt.Sprintf("%q: %s", err.Error(), sqlStmt)}		
	}

	return nil
}

func (r *SqlRepository) Add(player *model.Player) error {
	// start db transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// create player to bid insert statement
	stmt, err := tx.Prepare("insert into players(playerId, maxPrice) values(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	
	// set transaction values
	_, err = stmt.Exec(player.Id, player.MaxPrice)
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

func (r *SqlRepository) GetList() ([]*model.Player, error) {
	// fetch players to bid from db
	rows, err := r.db.Query(`select * from players`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := []*model.Player{}

	// convert db rows to player model
	for rows.Next() {
		player := &model.Player{}

		err = rows.Scan(&player.Id, &player.MaxPrice)
		if err != nil {
			return nil, err
		}

		players = append(players, player)
	}

	return players, nil
}