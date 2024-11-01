package main

import (
	"os"
	"strings"
	"time"

	"github.com/ramadhia/bosnet/be/internal/config"
	"github.com/ramadhia/bosnet/be/internal/provider"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bosnet-be",
	Short: "bostnet be",
}

func init() {
	// load config of the env
	err := config.Load()
	if err != nil {
		panic(err)
	}

	initLogging()
}

func main() {
	rootCmd := registerCommands(&provider.DefaultProviderBuilder{})
	if err := rootCmd.Execute(); err != nil {
		logrus.Error(err.Error())
		os.Exit(1)
	}
}

func registerCommands(builder provider.ProviderBuilder) *cobra.Command {
	rootCmd.AddCommand(Server(builder))
	rootCmd.AddCommand(Migrate())

	// we can register other commands like worker to consume the message broker
	//rootCmd.AddCommand(Worker(builder))

	return rootCmd
}

func initLogging() *logrus.Logger {
	cfg := config.Instance()
	log := logrus.StandardLogger()
	log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.RFC3339Nano,
	})
	if strings.ToLower(cfg.Log.Format) == "json" {
		log.SetFormatter(&logrus.JSONFormatter{})
	}
	level, err := logrus.ParseLevel(cfg.Log.Level)
	if err != nil {
		panic(err)
	}
	log.SetLevel(level)
	return log
}
