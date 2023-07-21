package player

import (
	"database/sql"
	"errors"
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type dbPlayer struct {
	playerId, maxPrice int
	deadline           string
}

type listTest struct {
	players      []dbPlayer
	expectedSize int
}

func createSqlitePlayerRepository() (*sql.DB, sqlmock.Sqlmock, *sqlitePlayerRepository) {
	db, mock, _ := sqlmock.New()
	r := NewSqlitePlayerRepository(db)

	return db, mock, r
}

func TestCreateSchemaFailureWhenDbError(t *testing.T) {
	db, _, r := createSqlitePlayerRepository()
	defer db.Close()

	if err := r.CreateSchema(); err == nil || !errors.Is(err, ErrCreateSchemaFailed) {
		t.Errorf("expected '%v' but '%v' returned", ErrCreateSchemaFailed, err)
	}
}

func TestCreateSchemaSuccess(t *testing.T) {
	db, m, r := createSqlitePlayerRepository()
	defer db.Close()

	m.ExpectExec("create table .+").WillReturnResult(sqlmock.NewResult(0, 0))

	if err := r.CreateSchema(); err != nil {
		t.Errorf("expected <nil> but '%v' returned", err)
	}
}

func TestInitFailureWhenCanNotCreateDbSchema(t *testing.T) {
	db, _, r := createSqlitePlayerRepository()
	defer db.Close()

	if err := r.Init(); err == nil || !errors.Is(err, ErrCanNotCreateDbSchema) {
		t.Errorf("expected '%v' but '%v' returned", ErrCanNotCreateDbSchema, err)
	}
}

func TestInitSuccess(t *testing.T) {
	db, m, r := createSqlitePlayerRepository()
	defer db.Close()

	m.ExpectExec("create table .+").WillReturnResult(sqlmock.NewResult(0, 0))

	if err := r.Init(); err != nil {
		t.Errorf("expected <nil> but '%v' returned", err)
	}
}

func TestListFailureWhenCanNotFetchPlayers(t *testing.T) {
	db, m, r := createSqlitePlayerRepository()
	defer db.Close()

	m.ExpectQuery("select .+").WillReturnError(sql.ErrConnDone)

	if _, err := r.List(); err == nil || !errors.Is(err, ErrCanNotFetchPlayers) {
		t.Errorf("expected error '%v' but '%v' returned", ErrCanNotFetchPlayers, err)
	}
}

func TestList(t *testing.T) {
	cases := []listTest{
		{[]dbPlayer{}, 0},
		{[]dbPlayer{{-1, -1, ""}}, 0}, // map error, can't use negative player index
		{[]dbPlayer{{1, 1, ""}}, 0},   // deadline date can not be parsed
		{[]dbPlayer{{1, 1, "2006-01-02 15:04:05"}}, 1},
	}
	db, m, r := createSqlitePlayerRepository()
	defer db.Close()

	for i, c := range cases {
		rows := sqlmock.NewRows([]string{"playerId", "maxPrice", "deadline"})
		for _, p := range c.players {
			rows.AddRow(p.playerId, p.maxPrice, p.deadline)
		}

		m.ExpectQuery("select .+").WillReturnRows(rows)

		p, _ := r.List()
		pCount := len(p)
		if pCount != c.expectedSize {
			t.Errorf("(run no. %d): expected %d players in list but %d players returned", i, c.expectedSize, pCount)
		}
	}
}
