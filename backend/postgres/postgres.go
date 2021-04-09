package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/juju/errors"
	"github.com/ovh/configstore"
)

// SQL represents a sql instance.
type SQL struct {
	// conn is the database handler.
	conn *pgxpool.Pool
}

var dbNameKey = "dbName"
var dbUserKey = "dbUser"
var dbPasswordKey = "dbPassword"
var dbHostKey = "dbHost"
var dbPortKey = "dbPort"
var dbSSLModeKey = "dbSSLMode"

// New initializes a new sql instance.
func New() (*SQL, error) {
	connectionString, err := buildConnectionString()
	if err != nil {
		return nil, err
	}

	conn, err := pgxpool.Connect(context.Background(), connectionString)
	if err != nil {
		return nil, err
	}

	return &SQL{
		conn: conn,
	}, nil

}

// Close gracefully closes the database connection.
func (s *SQL) Close() error {
	s.conn.Close()
	return nil
}

// Ping pings the database.
func (s *SQL) Ping() error {
	return s.conn.Ping(context.Background())
}

// ExecTransaction builds and executes a query.
func (s *SQL) ExecTransaction(query string, arguments ...interface{}) error {
	tx, err := s.conn.Begin(context.Background())
	if err != nil {
		return err
	}

	// Rollback is safe to call even if the tx is already closed, so if
	// the tx commits successfully, this is a no-op
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), query, arguments...)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

// QueryRow builds and executes a query to select data.
func (s *SQL) QueryRow(query string, arguments ...interface{}) pgx.Row {
	return s.conn.QueryRow(context.Background(), query, arguments...)
}

// buildConnectionString builds the connection string
// used to connect to the sql database
func buildConnectionString() (string, error) {
	contextErr := "building connection string"
	dbName, err := configstore.GetItemValue(dbNameKey)
	if err != nil {
		return "", errors.Annotatef(err, contextErr)
	}

	dbUser, err := configstore.GetItemValue(dbUserKey)
	if err != nil {
		return "", errors.Annotatef(err, contextErr)
	}

	dbPassword, err := configstore.GetItemValue(dbPasswordKey)
	if err != nil {
		return "", errors.Annotatef(err, contextErr)
	}

	dbHost, err := configstore.GetItemValue(dbHostKey)
	if err != nil {
		return "", errors.Annotatef(err, contextErr)
	}

	dbPort, err := configstore.GetItemValue(dbPortKey)
	if err != nil {
		return "", errors.Annotatef(err, contextErr)
	}

	dbSSLMode, err := configstore.GetItemValue(dbSSLModeKey)
	if err != nil {
		return "", errors.Annotatef(err, contextErr)
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode), nil
}
