package cmd

import (
	"halsecur/cli"
	"halsecur/cli/bisecur"
	"halsecur/cli/utils"
	"os"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

func init() {
	var devicePort int

	setStateCmd := &cobra.Command{
		Use:     SetStateCmdName,
		Short:   "Open or close a door connected to your Hörmann BiSecur gateway.",
		Long:    ``,
		PreRunE: preRunFuncs,
		Run: func(cmd *cobra.Command, args []string) {
			deviceMac := viper.GetString(ArgNameDeviceMac)
			host := viper.GetString(ArgNameHost)
			port := viper.GetInt(ArgNamePort)
			token := viper.GetUint32(ArgNameToken)

			mac, err := utils.ParesMacString(deviceMac)
			if err != nil {
				cli.Log.Fatalf("%v", err)
				os.Exit(1)
			}

			devicePortByte, err := utils.SafeLen(devicePort)
			if err != nil {
				cli.Log.Fatalf("too big devicePort value. %v", err)
				os.Exit(2)
			}

			err = bisecur.SetState(localMac, mac, host, port, devicePortByte, token)
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
