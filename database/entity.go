package database

type Entity interface {
	TableName() string
	Joins() string
}
