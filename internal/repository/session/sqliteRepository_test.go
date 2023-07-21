package session

import (
	"database/sql"
	"errors"
	"sokker-org-auto-bidder/internal/repository"
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func createSqliteSessionRepository() (*sql.DB, sqlmock.Sqlmock, SessionRepository) {
	db, mock, _ := sqlmock.New()
	r := NewSqliteSessionRepository(db)

	return db, mock, r
}

func TestInitFailureWhenCanNotCreateDbSchema(t *testing.T) {
	db, _, r := createSqliteSessionRepository()
	defer db.Close()

	if err := r.Init(); err == nil || !errors.Is(err, repository.ErrCanNotCreateDbSchema) {
		t.Errorf("expected '%v' but '%v' returned", repository.ErrCanNotCreateDbSchema, err)
	}
}

func TestInitSuccess(t *testing.T) {
	db, m, r := createSqliteSessionRepository()
	defer db.Close()

	m.ExpectExec("create table .+").WillReturnResult(sqlmock.NewResult(0, 0))

	if err := r.Init(); err != nil {
		t.Errorf("expected <nil> but '%v' returned", err)
	}
}

func TestGetFailureWhenNoSessionKey(t *testing.T) {
	db, _, r := createSqliteSessionRepository()
	defer db.Close()

	if _, err := r.Get(); err == nil || !errors.Is(err, ErrNoSessionKey) {
		t.Errorf("expected '%v' but '%v' returned", ErrNoSessionKey, err)
	}
}

func TestGetSuccess(t *testing.T) {
	expectedKey := "foo"
	db, m, r := createSqliteSessionRepository()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"key"}).AddRow(expectedKey)
	m.ExpectQuery("select .+").WillReturnRows(rows)

	key, err := r.Get()
	if err != nil {
		t.Errorf("expected <nil> error but '%v' returned", err)
	}
	if key != expectedKey {
		t.Errorf("expected '%s' key but '%s' returned", expectedKey, key)
	}
}
