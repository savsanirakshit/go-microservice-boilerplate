package db

import (
	"errors"
	"flag"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"golang-microservice-boilerplate/common"
	"golang-microservice-boilerplate/logger"
	"golang-microservice-boilerplate/model"
	"log"
	"strconv"
)

var Conn *pg.DB

func Connect() (*pg.DB, error) {
	currentDir := common.CurrentWorkingDir()
	dbHost := common.GetEnv("DB_HOST", "localhost")
	dbPort := common.GetEnv("DB_PORT", strconv.Itoa(5432))
	dbUser := common.GetEnv("DB_USER", "postgres")
	dbPass := common.GetEnv("DB_PASS", "postgres")
	dbName := common.GetEnv("DB_NAME", "microservice")
	recreateDB := common.GetEnv("RECREATE_DB", "true")
	migrationFilePath := common.GetEnv("MIGRATION_PATH", currentDir+"/db/migrations")

	dbConnString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)
	Conn = pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%s", dbHost, dbPort),
		User:     dbUser,
		Password: dbPass,
		Database: dbName,
	})

	reCreateDb, _ := strconv.ParseBool(recreateDB)
	if reCreateDb {
		dropAndCreateSchema(Conn, dbUser)
		err := createSchemaTable(Conn)
		if err != nil {
			return nil, err
		}
	} else {
		flag.Parse()
		m, err := migrate.New("file:"+migrationFilePath, dbConnString)
		if err != nil {
			log.Fatal(err)
		}
		if err := m.Up(); err != nil {
			if errors.As(err, &migrate.ErrDirty{}) {
				version, _, _ := m.Version()
				m.Force(int(version) - 1)
				if err := m.Up(); err != nil {
					m.Force(int(version) - 1)
					logger.ServiceLogger.Error(err)
				}
				logger.ServiceLogger.Warn("Database is already up-to-date with version : ", version)
			} else {
				logger.ServiceLogger.Error(err)
			}
		} else {
			version, _, _ := m.Version()
			logger.ServiceLogger.Info("Database successfully migrated with version : ", version)
		}
	}
	return Conn, nil
}

func dropAndCreateSchema(conn *pg.DB, dbUser string) error {
	tx, err := conn.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err := tx.Rollback(); err != nil {
			fmt.Printf(err.Error())
		}
	}()

	_, err = tx.Exec("DROP SCHEMA public cascade")
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}

	_, err = tx.Exec("CREATE SCHEMA public AUTHORIZATION " + dbUser)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func createSchemaTable(conn *pg.DB) error {
	models := models()
	for _, model := range models {
		err := conn.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func models() []interface{} {
	models := []interface{}{
		(*model.User)(nil),
	}
	return models
}
