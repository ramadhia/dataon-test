package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ramadhia/mnc-test/internal/provider"
	"github.com/ramadhia/mnc-test/internal/usecase/algo"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func AlgoTest3(pb provider.ProviderBuilder) *cobra.Command {
	cliCommand := &cobra.Command{
		Use:   "algo-test3",
		Short: "Run the test algorithm",
		Run: func(cmd *cobra.Command, args []string) {
			logger := logrus.WithField("method", "Algo")
			logger.Info("Algo Test3 - Enter a new line for input the arguments")

			app := provider.NewProvider()
			init := algo.NewAlgo(app)

			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				line := strings.ReplaceAll(scanner.Text(), "\n", "")
				ret, err := init.AlgoTest3(context.Background(), line)
				if err != nil {
					logrus.WithError(err).Warning("Error when run cli")
					return
				}

				fmt.Println(fmt.Sprintf("\nOutput: %v\n", ret))
			}

			if err := scanner.Err(); err != nil {
				fmt.Println(err)
			}
		},
	}
	return cliCommand
}
