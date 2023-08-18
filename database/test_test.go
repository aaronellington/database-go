package database_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aaronellington/database-go/database"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
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

func getTestConnection(t *testing.T, configurators ...database.ConnectionConfigurator) *database.Connection {
	pool, err := dockertest.NewPool("")
	if err != nil {
		panic(err)
	}

	err = pool.Client.Ping()
	if err != nil {
		panic(err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mysql",
		Tag:        "5.7",
		Env:        []string{"MYSQL_ROOT_PASSWORD=secret"},
	}, func(hc *docker.HostConfig) {
		hc.Mounts = append(hc.Mounts,
			docker.HostMount{
				Type:   "tmpfs",
				Target: "/var/lib/mysql",
			},
			docker.HostMount{
				Type:   "tmpfs",
				Target: "/var/log/mysql",
			},
		)
	})
	if err != nil {
		panic(err)
	}

	if err := mysql.SetLogger(log.New(io.Discard, "", log.LstdFlags)); err != nil {
		panic(err)
	}

	var db *sqlx.DB

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

	configurators = append(
		[]database.ConnectionConfigurator{
			database.WithPreRunHook(func(statement string, parameters map[string]any) error {
				log.Print(statement)
				log.Print(parameters)

				return nil
			}),
		},
		configurators...,
	)

	connection := database.NewConnection(db, configurators...)

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

	t.Cleanup(func() {
		if err := pool.Purge(resource); err != nil {
			panic(err)
		}
	})

	return connection
}
