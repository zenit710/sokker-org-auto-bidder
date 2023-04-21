package tools_test

import (
	"errors"
	"fmt"
	"sokker-org-auto-bidder/tools"
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestMakeTransactionBeginTransactionError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin().WillReturnError(fmt.Errorf("db transaction could not begin"))

	expectedError := fmt.Errorf("db transaction could not begin")
	err = tools.MakeTransaction(db, "sql")
	if err == nil || !errors.Is(err, expectedError) {
		// need to make export error from tools
		t.Errorf("expected '%v' error, got '%v'", expectedError, err)
	}
}
