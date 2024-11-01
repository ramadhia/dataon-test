package main

import (
	"os"

	"github.com/ramadhia/bosnet/be/internal/config"
	"github.com/ramadhia/bosnet/be/internal/storage"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var migrationPath string
var rollback bool
var versionToForce int

func Migrate() *cobra.Command {
	cliCommand := &cobra.Command{
		Use:   "migrate",
		Short: "Run the database migration",
		Run:   migrate,
	}
	cliCommand.Flags().StringVarP(&migrationPath, "path", "p", config.Instance().DB.Migration.Path, "The migration folder")
	cliCommand.Flags().BoolVarP(&rollback, "rollback", "r", false, "Rollback to prev migration (-1 step)")
	cliCommand.Flags().IntVarP(&versionToForce, "force", "f", -1, "Force to specific version")

	return cliCommand
}

func migrate(cmd *cobra.Command, args []string) {
	db := storage.GetSqlServerDb()
	defer storage.CloseDB(db)

	sqlDb, err := db.DB()
	if err != nil {
		logrus.Errorf("Error when migration: %s", err.Error())
		os.Exit(1)
	}

	err = storage.MigrateSqlServerDb(sqlDb, &migrationPath, rollback, versionToForce)
	if err != nil {
		logrus.Errorf("Error when migration: %s", err.Error())
		os.Exit(1)
	}

	logrus.Info("Finish migrating database")
}
