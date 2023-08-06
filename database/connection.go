package database

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

func NewConnection(
	db *sqlx.DB,
	configOptions ...connectionConfig,
) *Connection {
	connection := &Connection{
		db:          db,
		preRunFuncs: []PreRunFunc{},
	}

	for _, configOption := range configOptions {
		configOption(connection)
	}

	return connection
}

type Connection struct {
	db          *sqlx.DB
	preRunFuncs []PreRunFunc
}

func (c *Connection) Select(ctx context.Context, target any, statement string, parameters map[string]any) error {
	statement = sanitizeStatement(statement)

	if parameters == nil {
		parameters = map[string]any{}
	}

	for _, preRunFunc := range c.preRunFuncs {
		if err := preRunFunc(statement, parameters); err != nil {
			return err
		}
	}

	s, err := c.db.PrepareNamed(statement)
	if err != nil {
		return err
	}

	if err := s.Select(target, parameters); err != nil {
		return err
	}

	return nil
}

func (c *Connection) Execute(ctx context.Context, statement string, parameters map[string]any) (sql.Result, error) {
	statement = sanitizeStatement(statement)

	if parameters == nil {
		parameters = map[string]any{}
	}

	for _, preRunFunc := range c.preRunFuncs {
		if err := preRunFunc(statement, parameters); err != nil {
			return nil, err
		}
	}

	return c.db.NamedExec(statement, parameters)
}
