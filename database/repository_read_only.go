package database

import "context"

func NewReadOnlyRepository[T Entity](connection *Connection) *ReadOnlyRepository[T] {
	e := *new(T)

	return &ReadOnlyRepository[T]{
		connection:     connection,
		queryStatement: GenerateQuery(e),
	}
}

type ReadOnlyRepository[T Entity] struct {
	connection     *Connection
	queryStatement Query
}

func (repo ReadOnlyRepository[T]) Find(ctx context.Context, query Query, parameters map[string]any) ([]T, error) {
	// Override the values we need
	query.Select = repo.queryStatement.Select
	query.From = repo.queryStatement.From
	query.Join = repo.queryStatement.Join

	target := []T(nil)

	if err := repo.connection.Select(ctx, &target, query.String(), parameters); err != nil {
		return []T(nil), err
	}

	return target, nil
}

func (repo ReadOnlyRepository[T]) FindOne(ctx context.Context, query Query, parameters map[string]any) (T, error) {
	// Override the values we need
	query.Select = repo.queryStatement.Select
	query.From = repo.queryStatement.From
	query.Join = repo.queryStatement.Join

	target := []T(nil)

	if err := repo.connection.Select(ctx, &target, query.String(), parameters); err != nil {
		return *new(T), err
	}

	return target[0], nil
}
