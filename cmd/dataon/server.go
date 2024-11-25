package main

import (
	"context"

	"github.com/ramadhia/dataon-test/internal/handler/http"
	"github.com/ramadhia/dataon-test/internal/provider"
	"github.com/ramadhia/dataon-test/internal/repository/gorm"
	"github.com/ramadhia/dataon-test/internal/usecase/organization"
	"github.com/ramadhia/dataon-test/internal/usecase/user"

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
			app.SetGroupRepo(gorm.NewGroupRepo(db))
			app.SetOrganizationRepo(gorm.NewOrganizationRepo(db))
			app.SetUserRepo(gorm.NewUserRepo(db))
			app.SetGroupRepo(gorm.NewGroupRepo(db))

			// -- set usecase
			app.SetOrganizationUseCase(organization.NewOrganization(app))
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
