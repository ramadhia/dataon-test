package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/ramadhia/mnc-test/internal/provider"
	"github.com/ramadhia/mnc-test/internal/usecase/algo"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func AlgoTest1(pb provider.ProviderBuilder) *cobra.Command {
	cliCommand := &cobra.Command{
		Use:   "algo-test1",
		Short: "Run the test algorithm",
		Run: func(cmd *cobra.Command, args []string) {
			logger := logrus.WithField("method", "Algo")
			logger.Info("Algo Test1 - Enter a new line for input the arguments")

			app := provider.NewProvider()
			init := algo.NewAlgo(app)

			scanner := bufio.NewScanner(os.Stdin)
			var input []string
			for scanner.Scan() {
				line := scanner.Text()
				if line == "" {
					break
				}

				input = append(input, line)
			}

			ret, err := init.AlgoTest1(context.Background(), input)
			if err != nil {
				logrus.WithError(err).Warning("Error when run cli")
				return
			}
			if err := scanner.Err(); err != nil {
				fmt.Println(err)
			}

			fmt.Println(ret)
		},
	}
	return cliCommand
}
