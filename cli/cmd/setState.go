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

	setStateCmd := &cobra.Command{
		Use:     SetStateCmdUse,
		Short:   "Open or close a door connected to your Hörmann BiSecur gateway.",
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

			err = bisecur.SetState(localMac, mac, host, port, byte(devicePort), token)
			if err != nil {
				cli.Log.Fatalf("%v", err)
				os.Exit(2)
			}

			cli.Log.Infof("Success")
		},
	}

	rootCmd.AddCommand(setStateCmd)

	setStateCmd.Flags().IntVar(&devicePort, ArgDevicePortName, 0, "Port number of the door")
	setStateCmd.MarkFlagsOneRequired(ArgDevicePortName)
}
