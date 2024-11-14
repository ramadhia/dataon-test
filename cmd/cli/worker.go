package main

import (
	"context"
	"github.com/ramadhia/mnc-test/internal/repository/gorm"
	"github.com/ramadhia/mnc-test/internal/usecase/transcation"
	"strings"

	"github.com/ramadhia/mnc-test/internal/handler/messaging"
	"github.com/ramadhia/mnc-test/internal/provider"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var flagWorkerTopics string
var flagWorkerNotTopics string

func Worker(pb provider.ProviderBuilder) *cobra.Command {
	cliCommand := &cobra.Command{
		Use:   "worker",
		Short: "Run the worker",
		Run: func(cmd *cobra.Command, args []string) {
			app, closeResourcesFn, err := pb.Build(context.Background())
			if err != nil {
				panic(err)
			}
			if closeResourcesFn != nil {
				defer closeResourcesFn()
			}

			db := app.DB()

			// -- set repo
			app.SetTransactionRepo(gorm.NewTransactionRepo(db))
			app.SetBalanceRepo(gorm.NewBalanceRepo(db))

			// -- set usecase
			app.SetTransactionUseCase(transcation.NewTransaction(app))

			// topic whitelist
			var topics []string
			if flagWorkerTopics != "" {
				topics = strings.Split(flagWorkerTopics, ",")
				for idx, topic := range topics {
					topics[idx] = strings.TrimSpace(topic)
				}
			}

			err = messaging.RegisterSubscribers(app, topics)
			if err != nil {
				panic(err)
			}

			logrus.Info("Start consuming. Press Ctrl+C to stop")
			err = app.MessageBus().StartConsuming()
			if err != nil {
				panic(err)
			}
		},
	}

	cliCommand.Flags().StringVarP(&flagWorkerTopics, "topics", "t", "", "The topics to listen to (comma delimited) or empty to listen to all topics")
	cliCommand.Flags().StringVarP(&flagWorkerNotTopics, "not-topics", "n", "", "The topics excluded by the worker (comma delimited)")

	return cliCommand
}
