package main

import (
	"context"

	"github.com/ramadhia/mnc-test/internal/handler/http"
	"github.com/ramadhia/mnc-test/internal/provider"
	"github.com/ramadhia/mnc-test/internal/repository/gorm"
	"github.com/ramadhia/mnc-test/internal/usecase/transcation"
	"github.com/ramadhia/mnc-test/internal/usecase/user"

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

			// -- set service

			// -- set repo
			app.SetTransactionRepo(gorm.NewTransactionRepo(db))
			app.SetBalanceRepo(gorm.NewBalanceRepo(db))
			app.SetUserRepo(gorm.NewUserRepo(db))

			// -- set usecase
			app.SetTransactionUseCase(transcation.NewTransaction(app))
			app.SetUserUseCase(user.NewUser(app))

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
