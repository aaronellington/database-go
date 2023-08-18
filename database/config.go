package database

type ConnectionConfigurator func(connection *Connection)

type PreRunHook func(statement string, parameters map[string]any) error

func WithPreRunHook(preRunHook PreRunHook) ConnectionConfigurator {
	return func(connection *Connection) {
		connection.preRunHooks = append(connection.preRunHooks, preRunHook)
	}
}
