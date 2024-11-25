package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ramadhia/dataon-test/internal/config"

	"github.com/golang-migrate/migrate/v4"
	postgresMigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GetPostgresDb return sql connection
func GetPostgresDb() *gorm.DB {
	dbName := config.Instance().DB.Database
	return PostgresConn(&dbName)
}

func PostgresConn(dbName *string) *gorm.DB {
	dsn := getDsn(dbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("error: %v for %v", err.Error(), dsn))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("error: %v for %v", err.Error(), dsn))
	}

	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(90))
	//sqlDB.SetMaxOpenConns(0)
	sqlDB.SetMaxIdleConns(5)

	return db
}

//func CreatePostgresDb(dbName string) error {
//	dbConn := PostgresConn(nil)
//	defer func(dbConn *gorm.DB) {
//		err := dbConn.Close()
//		if err != nil {
//			logrus.WithError(err).Warning("Error when closing mysql db")
//		}
//	}(dbConn)
//	return dbConn.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName)).Error
//}

// MigratePostgresDb init postgres database migration
func MigratePostgresDb(db *sql.DB, migrationFolder *string, rollback bool, versionToForce int) error {
	logger := logrus.WithField("method", "storage.MigratePostgresDb")
	dbConfig := config.Instance().DB

	var validMigrationFolder = dbConfig.Migration.Path
	if migrationFolder != nil && *migrationFolder != "" {
		validMigrationFolder = *migrationFolder
	}

	if validMigrationFolder == "" {
		return fmt.Errorf("empty migration folder")
	}
	logger.Infof("Migration folder: %s", validMigrationFolder)

	driver, err := postgresMigrate.WithInstance(db, &postgresMigrate.Config{})
	if err != nil {
		logger.WithError(err).Warning("Error when instantiating driver")
		return err
	}
	m, err := migrate.NewWithDatabaseInstance("file://"+validMigrationFolder,
		dbConfig.Client,
		driver)
	if err != nil {
		logger.WithError(err).Warning("Error when instantiating migrate")
		return err
	}
	if rollback {
		logger.Info("About to Rolling back 1 step")
		err = m.Steps(-1)
	} else if versionToForce != -1 {
		logger.Info(fmt.Sprintf("About to force version %d", versionToForce))
		err = m.Force(versionToForce)
	} else {
		logger.Info("About to run migration")
		err = m.Up()
	}
	if err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
	}

	return nil
}

func CloseDB(db *gorm.DB) {
	if db == nil {
		return
	}

	if sqlDB, err := db.DB(); err != nil {
		logrus.Warnf("Error when get db connection: %s", err)
	} else {
		err = sqlDB.Close()
		if err != nil {
			logrus.Warnf("Error when closing db: %s", err)
		}
	}
}

func getDsn(dbName *string) string {
	dbConfig := config.Instance().DB

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, *dbName, dbConfig.Password)
	return dsn
}
