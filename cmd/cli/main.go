package main

import (
	"os"
	"strings"
	"time"

	"github.com/ramadhia/mnc-test/internal/config"
	"github.com/ramadhia/mnc-test/internal/provider"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cli-go",
	Short: "cli go",
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
	rootCmd.AddCommand(Worker(builder))
	rootCmd.AddCommand(Migrate())
	rootCmd.AddCommand(AlgoTest1(builder))
	rootCmd.AddCommand(AlgoTest2(builder))
	rootCmd.AddCommand(AlgoTest3(builder))
	rootCmd.AddCommand(AlgoTest4(builder))

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
