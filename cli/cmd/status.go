package cmd

import (
	"bisecur/cli"
	"bisecur/sdk"
	"bisecur/sdk/payload"
	"github.com/spf13/viper"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	var devicePort int

	statusCmd := &cobra.Command{
		Use:     StatusCmdUse,
		Short:   "Queries the status (open/closed/etc) of your door.",
		Long:    ``,
		PreRunE: preRunFuncs,
		Run: func(cmd *cobra.Command, args []string) {
			deviceMac := viper.GetString(ArgNameDeviceMac)
			host := viper.GetString(ArgNameHost)
			port := viper.GetInt(ArgNamePort)
			token := viper.GetUint32(ArgNameToken)

			mac, err := cli.ParesMacString(deviceMac)
			if err != nil {
				log.Fatalf("%v", err)
				os.Exit(1)
			}

			err = getStatus(localMac, mac, host, port, byte(devicePort), token)
			if err != nil {
				log.Fatalf("%v", err)
				os.Exit(2)
			}
		},
	}

	rootCmd.AddCommand(statusCmd)

	statusCmd.Flags().IntVar(&devicePort, devicePortName, 0, "Port number of the door")
	statusCmd.MarkFlagsOneRequired(devicePortName)
}

func getStatus(localMac [6]byte, mac [6]byte, host string, port int, devicePort byte, token uint32) error {
	client := sdk.NewClient(log, localMac, mac, host, port, token)
	err := client.Open()
	if err != nil {
		return err
	}

	defer func() {
		err2 := client.Close()
		if err2 != nil {
			log.Errorf("%v", err2)
		}
	}()

	var status *payload.HmGetTransitionResponse
	err = retry(func() error {
		var err2 error
		status, err2 = client.GetTransition(devicePort)
		return err2
	})

	if err != nil {
		return err
	}

	log.WithField("status", status).Infof("Success")

	return nil
}
