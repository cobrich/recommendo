package repo

import (
	"context"
	"database/sql"
)

// DBTX - это интерфейс, которому удовлетворяют и *sql.DB, и *sql.Tx.
// Это позволяет вашим методам репозитория работать внутри транзакций.
type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}
