package utils

import (
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

// GetColumnFromConstraintName returns the column name from the constraint name
func GetColumnFromConstraint(pgErr *pgconn.PgError) string {
	return strings.Split(pgErr.ConstraintName, "_")[1]
}
