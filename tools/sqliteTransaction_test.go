package tools_test

import (
	"database/sql"
	"errors"
	"sokker-org-auto-bidder/tools"
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func initSqlMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestMakeTransactionBeginTransactionError(t *testing.T) {
	db, mock := initSqlMock(t)
	defer db.Close()

	mock.ExpectBegin().WillReturnError(tools.ErrTransactionBeginFailed)

	err := tools.MakeTransaction(db, "sql")
	if err == nil || !errors.Is(err, tools.ErrTransactionBeginFailed) {
		t.Errorf("expected '%v' error, got '%v'", tools.ErrTransactionBeginFailed, err)
	}
}

func TestMakeTransactionSqlPrepareError(t *testing.T) {
	db, mock := initSqlMock(t)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectPrepare("sql").WillReturnError(&tools.ErrSqlPrepareFailed{})

	var expectedErrType *tools.ErrSqlPrepareFailed
	err := tools.MakeTransaction(db, "sql")
	if err == nil || !errors.As(err, &expectedErrType) {
		t.Errorf("expected '%T' error, got '%T'", expectedErrType, err)
	}
}

func TestMakeTransactionSqlExecError(t *testing.T) {
	db, mock := initSqlMock(t)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectPrepare("sql")
	mock.ExpectExec("sql").WillReturnError(&tools.ErrSqlExecFailed{})

	var expectedErrType *tools.ErrSqlExecFailed
	err := tools.MakeTransaction(db, "sql")
	if err == nil || !errors.As(err, &expectedErrType) {
		t.Errorf("expected '%T' error, got '%T'", expectedErrType, err)
	}
}

func TestMakeTransactionCommitTransactionError(t *testing.T) {
	db, mock := initSqlMock(t)
	defer db.Close()

	var lastInsertID, affected int64
	result := sqlmock.NewResult(lastInsertID, affected)

	mock.ExpectBegin()
	mock.ExpectPrepare("sql")
	mock.ExpectExec("sql").WillReturnResult(result)
	mock.ExpectCommit().WillReturnError(tools.ErrTransactionCommitFailed)

	err := tools.MakeTransaction(db, "sql")
	if err == nil || !errors.Is(err, tools.ErrTransactionCommitFailed) {
		t.Errorf("expected '%s' error, got '%s'", tools.ErrTransactionCommitFailed, err)
	}
}

func TestMakeTransactionSuccess(t *testing.T) {
	db, mock := initSqlMock(t)
	defer db.Close()

	var lastInsertID, affected int64
	result := sqlmock.NewResult(lastInsertID, affected)

	mock.ExpectBegin()
	mock.ExpectPrepare("sql")
	mock.ExpectExec("sql").WillReturnResult(result)
	mock.ExpectCommit()

	if err := tools.MakeTransaction(db, "sql"); err != nil {
		t.Errorf("expected nil but got %T", err)
	}
}
