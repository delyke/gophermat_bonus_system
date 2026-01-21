package postgres

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestIsUniqueViolation(t *testing.T) {
	if IsUniqueViolation(nil) {
		t.Fatal("expected false for nil error")
	}

	if IsUniqueViolation(errors.New("boom")) {
		t.Fatal("expected false for non pg error")
	}

	if IsUniqueViolation(&pgconn.PgError{Code: "99999"}) {
		t.Fatal("expected false for non unique violation")
	}

	if !IsUniqueViolation(&pgconn.PgError{Code: "23505"}) {
		t.Fatal("expected true for unique violation")
	}
}
