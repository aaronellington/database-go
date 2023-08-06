package database

import (
	"context"
)

func NewRepository[T Entity](connection *Connection) *Repository[T] {
	e := *new(T)

	return &Repository[T]{
		ReadOnlyRepository: NewReadOnlyRepository[T](connection),
		saveStatement:      generateSave(e),
		deleteStatement:    generateDelete(e),
		insertStatement:    generateInsert(e),
		updateStatement:    generateUpdate(e),
	}
}

type Repository[T Entity] struct {
	*ReadOnlyRepository[T]

	saveStatement   string
	deleteStatement string
	insertStatement string
	updateStatement string
}

func (repo Repository[T]) Save(ctx context.Context, entity T) error {
	_, err := repo.connection.Execute(ctx, repo.saveStatement, convertStructToParams(entity))
	if err != nil {
		return err
	}

	return nil
}

func (repo Repository[T]) Delete(ctx context.Context, entity T) error {
	_, err := repo.connection.Execute(ctx, repo.deleteStatement, convertStructToParams(entity))
	if err != nil {
		return err
	}

	return nil
}

func (repo Repository[T]) Insert(ctx context.Context, entity T) (uint64, error) {
	result, err := repo.connection.Execute(ctx, repo.insertStatement, convertStructToParams(entity))
	if err != nil {
		return 0, err
	}

	newID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(newID), nil
}

func (repo Repository[T]) Update(ctx context.Context, entity T) error {
	_, err := repo.connection.Execute(ctx, repo.updateStatement, convertStructToParams(entity))
	if err != nil {
		return err
	}

	return nil
}
