package tests

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/sirupsen/logrus"
)

func readDataFromFile(file string) (string, error) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func getDeleteSQL() (string, error) {
	basePath, _ := os.Getwd()
	baseDir := filepath.Dir(basePath)
	schemaFilePath := filepath.Join(baseDir, "src", "repository", "drop_schema.sql")
	return readDataFromFile(schemaFilePath)
}

func getCreateSQL() (string, error) {
	basePath, _ := os.Getwd()
	baseDir := filepath.Dir(basePath)
	schemaFilePath := filepath.Join(baseDir, "src", "repository", "schema.sql")
	return readDataFromFile(schemaFilePath)
}

func SetupTestDB() (*pgxpool.Pool, func(), error) {
	PG_PORT_TEST := "5455"
	PG_HOST_TEST := "localhost"
	PG_DATABASE_TEST := "db_test"
	PG_USER_TEST := "postgres"
	PG_PASSWORD_TEST := "1234"

	postgresUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		PG_USER_TEST,
		PG_PASSWORD_TEST,
		PG_HOST_TEST,
		PG_PORT_TEST,
		PG_DATABASE_TEST,
	)
	config, err := pgxpool.ParseConfig(postgresUrl)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to parse database config: %s", err.Error())
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to connect to database: %s", err.Error())
	}

	delete_sql, err_d := getDeleteSQL()
	create_sql, err_c := getCreateSQL()

	if err_d != nil || err_c != nil {
		var err error
		if err_c != nil {
			err = err_c
		} else {
			err = err_d
		}
		return nil, nil, fmt.Errorf("unable to get sql commands: %s", err.Error())
	}

	_, err_del := pool.Exec(context.Background(), delete_sql)
	if err_del != nil {
		return nil, nil, fmt.Errorf("unable to perform deleting tables: %s", err_del.Error())
	}
	_, err_create := pool.Exec(context.Background(), create_sql)
	if err_create != nil {
		return nil, nil, fmt.Errorf("unable to perform creation tables: %s", err_create.Error())
	}

	cleanup := func() {
		_, err := pool.Exec(context.Background(), delete_sql)
		if err != nil {
			log.Printf("Failed to clean database: %s\n", err.Error())
		}
		pool.Close()
	}

	return pool, cleanup, nil
}

func SetupLogger() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	log.SetOutput(os.Stdout)

	log.SetLevel(logrus.InfoLevel)
	return log
}
