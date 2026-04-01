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

	portsRemoveCmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove a paired port from the gateway",
		Long:  `Remove a paired port (radio channel) from the gateway.`,
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

			err = bisecur.RemovePort(localMac, mac, host, port, devicePortByte, token)
			if err != nil {
				cli.Log.Fatalf("%v", err)
				os.Exit(2)
			}

			cli.Log.Infof("Success. Port %d removed.", devicePort)
		},
	}

	portsCmd.AddCommand(portsRemoveCmd)

	portsRemoveCmd.Flags().IntVar(&devicePort, ArgDevicePortName, 0, "Port number to remove")
	portsRemoveCmd.MarkFlagsOneRequired(ArgDevicePortName)
}
