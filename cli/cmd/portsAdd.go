package cmd

import (
	"halsecur/cli"
	"halsecur/cli/bisecur"
	"halsecur/cli/utils"
	"os"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var portsAddCmd = &cobra.Command{
	Use:   AddCmdUse,
	Short: "Pair a new door by cloning a hand remote signal (ADD_PORT)",
	Long: `Pair a new door by cloning a hand remote's radio signal.

The gateway enters receive mode and listens for ~40 seconds.
Press the hand remote button within ~10cm of the gateway.
In testing, the cloned port did not support position feedback.`,
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

		cli.Log.Infof("Listening for hand remote signal (~40s timeout). Press the remote button within ~10cm of the gateway...")

		portId, err := bisecur.AddPort(localMac, mac, host, port, token)
		if err != nil {
			cli.Log.Fatalf("%v", err)
			os.Exit(2)
		}

		cli.Log.Infof("Success. New port ID: %d", portId)
	},
}

func init() {
	portsCmd.AddCommand(portsAddCmd)
}
