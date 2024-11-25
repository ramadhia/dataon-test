package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ramadhia/dataon-test/internal/config"

	"github.com/golang-migrate/migrate/v4"
	sqlServerMigrate "github.com/golang-migrate/migrate/v4/database/sqlserver"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// GetSqlServerDb return sql connection
func GetSqlServerDb() *gorm.DB {
	dbConfig := config.Instance().DB

	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database)

	return SqlServerConn(dsn)
}

func SqlServerConn(dsn string) *gorm.DB {
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
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

// MigrateSqlServerDb init postgres database migration
func MigrateSqlServerDb(db *sql.DB, migrationFolder *string, rollback bool, versionToForce int) error {
	logger := logrus.WithField("method", "storage.MigrateSqlServerDb")
	dbConfig := config.Instance().DB

	var validMigrationFolder = dbConfig.Migration.Path
	if migrationFolder != nil && *migrationFolder != "" {
		validMigrationFolder = *migrationFolder
	}

	if validMigrationFolder == "" {
		return fmt.Errorf("empty migration folder")
	}
	logger.Infof("Migration folder: %s", validMigrationFolder)

	driver, err := sqlServerMigrate.WithInstance(db, &sqlServerMigrate.Config{})
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
