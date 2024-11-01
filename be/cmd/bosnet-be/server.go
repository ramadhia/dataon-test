package main

import (
	"context"
	"github.com/ramadhia/bosnet/be/internal/handler/http"
	"github.com/ramadhia/bosnet/be/internal/provider"
	"github.com/ramadhia/bosnet/be/internal/repository/gorm"
	"github.com/ramadhia/bosnet/be/internal/usecase/history"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Server(pb provider.ProviderBuilder) *cobra.Command {
	cliCommand := &cobra.Command{
		Use:   "server",
		Short: "Run the REST API server",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := logrus.WithField("method", "Server")

			app, closeResourcesFn, err := pb.Build(context.Background())
			if err != nil {
				panic(err)
			}
			if closeResourcesFn != nil {
				defer closeResourcesFn()
			}

			db := app.DB()

			// -- set repo
			app.SetHistoryRepo(gorm.NewHistoryRepo(db))
			app.SetBalanceRepo(gorm.NewBalanceRepo(db))

			// -- set usecase
			app.SetHistoryUseCase(history.NewHistory(app))

			// Start Http Server
			if err := http.NewHttpServer(app).Start(); err != nil {
				logger.WithError(err).Error("Error starting web server")
				return err
			}

			return nil
		},
	}
	return cliCommand
}
