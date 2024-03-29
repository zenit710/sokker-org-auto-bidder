package session

import (
	"database/sql"
	"errors"
	"fmt"
	"sokker-org-auto-bidder/internal/repository"
	"sokker-org-auto-bidder/tools"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

var (
	_               SessionRepository = &sqliteSessionRepository{}
	ErrNoSessionKey                   = errors.New("could not get session key")
)

// sqliteSessionRepository handle sqlite connection for sokker.org session
type sqliteSessionRepository struct {
	db *sql.DB
}

// NewSqliteSessionRepository returns new repository with sqlite connection
func NewSqliteSessionRepository(db *sql.DB) *sqliteSessionRepository {
	log.Trace("creating new sqlite session repository")
	return &sqliteSessionRepository{db}
}

func (r *sqliteSessionRepository) Get() (string, error) {
	log.Trace("get newest session key")
	key := ""

	row := r.db.QueryRow(`select key from sessions order by id desc`)
	if err := row.Scan(&key); err != nil {
		log.Error(err)
		return key, ErrNoSessionKey
	}

	return key, nil
}

func (r *sqliteSessionRepository) Init() error {
	log.Trace("sqlite session repository init")
	if err := r.CreateSchema(); err != nil {
		log.Error(err)
		return repository.ErrCanNotCreateDbSchema
	}

	return nil
}

func (r *sqliteSessionRepository) Save(sess string) error {
	log.Tracef("save session in the db")

	return tools.MakeTransaction(
		r.db,
		`insert into sessions (key, created) values(?, datetime('now'))`,
		sess,
	)
}

// CreateSchema creates db structure if it is not created yet
func (r *sqliteSessionRepository) CreateSchema() error {
	sqlStmt := `create table if not exists sessions (
		id integer not null primary key autoincrement,
		key text not null,
		created text not null
		);`

	log.Trace("create database schema if not exists")
	if _, err := r.db.Exec(sqlStmt); err != nil {
		log.Error(err)
		return fmt.Errorf("create schema sql execution failed")
	}

	return nil
}
