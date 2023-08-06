package database_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aaronellington/database-go/database"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest"
)

type DatabaseConfig struct {
	User string
	Pass string
	Host string
	Port string
	Name string
}

func (config DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&allowCleartextPasswords=1",
		config.User,
		config.Pass,
		config.Host,
		config.Port,
		config.Name,
	)
}

func getTestConnection() (*database.Connection, func()) {
	var db *sqlx.DB

	pool, err := dockertest.NewPool("")
	if err != nil {
		panic(err)
	}

	err = pool.Client.Ping()
	if err != nil {
		panic(err)
	}

	resource, err := pool.Run("mysql", "5.7", []string{"MYSQL_ROOT_PASSWORD=secret"})
	if err != nil {
		panic(err)
	}

	if err := mysql.SetLogger(log.New(io.Discard, "", log.LstdFlags)); err != nil {
		panic(err)
	}

	if err := pool.Retry(func() error {
		var err error

		db, err = sqlx.Open("mysql", DatabaseConfig{
			User: "root",
			Pass: "secret",
			Host: "localhost",
			Name: "mysql",
			Port: resource.GetPort("3306/tcp"),
		}.DSN())
		if err != nil {
			return err
		}

		return db.Ping()
	}); err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Minute * 1)

	connection := database.NewConnection(db, database.WithPreRunFunc(func(statement string, parameters map[string]any) error {
		log.Print(statement)
		log.Print(parameters)

		return nil
	}))

	fileContents, err := os.ReadFile("test_files/base.sql")
	if err != nil {
		panic(err)
	}

	statements := strings.Split(string(fileContents), ";")

	for _, statement := range statements {
		statement = strings.TrimSpace(statement)

		if statement != "" {
			if _, err := connection.Execute(context.TODO(), statement, nil); err != nil {
				panic(err)
			}
		}
	}

	return connection, func() {
		if err := pool.Purge(resource); err != nil {
			panic(err)
		}
	}
}
