package backend

import (
	"github.com/jackc/pgx/v4"
)

type Backend interface {
	Close() error
	Ping() error
	ExecTransaction(query string, arguments ...interface{}) error
	QueryRow(query string, arguments ...interface{}) pgx.Row
}
