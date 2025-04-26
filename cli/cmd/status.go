package cmd

import (
	"bisecur/cli"
	"bisecur/cli/bisecur"
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
				cli.Log.Fatalf("%v", err)
				os.Exit(1)
			}

			status, err := bisecur.GetStatus(localMac, mac, host, port, byte(devicePort), token)
			if err != nil {
				cli.Log.Fatalf("%v", err)
				os.Exit(2)
			}

			cli.Log.WithField("status", status).Infof("Success")
		},
	}

	rootCmd.AddCommand(statusCmd)

	statusCmd.Flags().IntVar(&devicePort, ArgDevicePortName, 0, "Port number of the door")
	statusCmd.MarkFlagsOneRequired(ArgDevicePortName)
}
