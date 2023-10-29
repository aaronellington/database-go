package database

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

func NewConnection(
	db *sqlx.DB,
	configOptions ...ConnectionConfigurator,
) *Connection {
	connection := &Connection{
		db:          db,
		preRunHooks: []PreRunHook{},
	}

	for _, configOption := range configOptions {
		configOption(connection)
	}

	return connection
}

type Connection struct {
	db          *sqlx.DB
	preRunHooks []PreRunHook
}

func (c *Connection) Select(ctx context.Context, target any, statement string, parameters map[string]any) error {
	statement = sanitizeStatement(statement)

	if parameters == nil {
		parameters = map[string]any{}
	}

	for _, preRunHook := range c.preRunHooks {
		if err := preRunHook(statement, parameters); err != nil {
			return err
		}
	}

	s, err := c.db.PrepareNamed(statement)
	if err != nil {
		return err
	}
	defer s.Close()

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

	for _, preRunHook := range c.preRunHooks {
		if err := preRunHook(statement, parameters); err != nil {
			return nil, err
		}
	}

	return c.db.NamedExec(statement, parameters)
}
