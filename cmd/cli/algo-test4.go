package main

import (
	"context"
	"fmt"
	"github.com/ramadhia/mnc-test/internal/provider"
	"github.com/ramadhia/mnc-test/internal/usecase"
	"github.com/ramadhia/mnc-test/internal/usecase/algo"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cutiBersama, cutiDurasi int
var joinDate, cutiDate string

func AlgoTest4(pb provider.ProviderBuilder) *cobra.Command {
	cliCommand := &cobra.Command{
		Use:   "algo-test4",
		Short: "Run the test algorithm",
		Run: func(cmd *cobra.Command, args []string) {
			logger := logrus.WithField("method", "Algo")
			logger.Info("Algo Test4 - Enter a new line for input the arguments")

			app := provider.NewProvider()
			init := algo.NewAlgo(app)

			can, res, err := init.AlgoTest4(context.Background(), usecase.AlgoTest4Request{
				CutiBersama: cutiBersama,
				CutiDurasi:  cutiDurasi,
				JoinDate:    joinDate,
				CutiDate:    cutiDate,
			})
			if err != nil {
				logrus.WithError(err).Warning("Error when run cli")
				return
			}

			fmt.Println(fmt.Sprintf("Output: %v, %s", can, res))
		},
	}

	cliCommand.Flags().IntVarP(&cutiBersama, "cuti-bersama", "a", 0, "Masukan input cuti bersama")
	cliCommand.Flags().IntVarP(&cutiDurasi, "cuti-durasi", "b", 0, "Masukan input cuti bersama")
	cliCommand.Flags().StringVarP(&joinDate, "join-date", "c", "", "Masukan input tanggal join")
	cliCommand.Flags().StringVarP(&cutiDate, "cuti-date", "d", "", "Masukan input tanggal cuti")

	return cliCommand
}
