package main

import (
	"context"
	"fmt"
	"github.com/ramadhia/mnc-test/internal/provider"
	"github.com/ramadhia/mnc-test/internal/usecase/algo"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var totalBelanja float32
var uangDibayar float32

func AlgoTest2(pb provider.ProviderBuilder) *cobra.Command {
	cliCommand := &cobra.Command{
		Use:   "algo-test2",
		Short: "Run the test algorithm",
		Run: func(cmd *cobra.Command, args []string) {
			logger := logrus.WithField("method", "Algo")
			logger.Info("Algo Test2 - Enter a new line for input the arguments")

			app := provider.NewProvider()
			init := algo.NewAlgo(app)

			ret, err := init.AlgoTest2(context.Background(), totalBelanja, uangDibayar)
			if err != nil {
				logrus.WithError(err).Warning("Error when run cli")
				return
			}

			fmt.Println(ret)
		},
	}

	cliCommand.Flags().Float32VarP(&totalBelanja, "total-belanja", "t", 0.0, "Masukan total belanja")
	cliCommand.Flags().Float32VarP(&uangDibayar, "uang-dibayar", "u", 0.0, "Masukan uang dibayar")

	return cliCommand
}
