package tools

import (
	"database/sql"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
)

var ErrTransactionBeginFailed = errors.New("db transaction could not begin")
var ErrTransactionCommitFailed = errors.New("could not commit sql transaction")

type ErrSqlPrepareFailed struct {
	sql string
}

func (e *ErrSqlPrepareFailed) Error() string {
	return fmt.Sprintf("could not prepare sql statement: '%s'", e.sql)
}

type ErrSqlExecFailed struct {
	sql    string
	params []interface{}
}

func (e *ErrSqlExecFailed) Error() string {
	return fmt.Sprintf("sql '%s' execution with params (%v) failed", e.sql, e.params)
}

// MakeTransaction handles DB transaction
func MakeTransaction(db *sql.DB, sql string, params ...interface{}) error {
	log.Trace("create db transaction")
	tx, err := db.Begin()
	if err != nil {
		log.Error(err)
		return ErrTransactionBeginFailed
	}

	log.Trace("prepare sql statement")
	stmt, err := tx.Prepare(sql)
	if err != nil {
		log.Error(err)
		return &ErrSqlPrepareFailed{sql}
	}
	defer stmt.Close()

	log.Trace("execute transaction")
	_, err = stmt.Exec(params...)
	if err != nil {
		log.Error(err)
		return &ErrSqlExecFailed{sql, params}
	}

	log.Trace("commit transaction")
	err = tx.Commit()
	if err != nil {
		log.Error(err)
		return ErrTransactionCommitFailed
	}

	return nil
}
