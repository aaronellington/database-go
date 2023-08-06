package database

type connectionConfig func(connection *Connection)

type PreRunFunc func(statement string, parameters map[string]any) error

func WithPreRunFunc(preRunFunc PreRunFunc) connectionConfig {
	return func(connection *Connection) {
		connection.preRunFuncs = append(connection.preRunFuncs, preRunFunc)
	}
}
