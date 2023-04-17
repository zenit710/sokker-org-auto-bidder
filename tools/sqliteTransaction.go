package tools

import (
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// MakeTransaction handles DB transaction
func MakeTransaction(db *sql.DB, sql string, params ...interface{}) error {
	log.Trace("create db transaction")
	tx, err := db.Begin()
	if err != nil {
		log.Error(err)
		return fmt.Errorf("db transaction could not begin")
	}

	log.Trace("prepare sql statement")
	stmt, err := tx.Prepare(sql)
	if err != nil {
		log.Error(err)
		return fmt.Errorf("could not prepare sql statement: '%s'", sql)
	}
	defer stmt.Close()
	
	log.Trace("execute transaction")
	_, err = stmt.Exec(params...)
	if err != nil {
		log.Error(err)
		return fmt.Errorf("sql '%s' execution with params (%v) failed", sql, params)
	}

	log.Trace("commit transaction")
	err = tx.Commit()
	if err != nil {
		log.Error(err)
		return fmt.Errorf("could not commit sql transaction")
	}

	return nil
}
