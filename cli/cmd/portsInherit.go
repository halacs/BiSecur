package cmd

import (
	"halsecur/cli"
	"halsecur/cli/bisecur"
	"halsecur/cli/utils"
	"os"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var portsInheritCmd = &cobra.Command{
	Use:   "inherit",
	Short: "Pair a new door by transmitting the gateway's radio code (INHERIT_PORT)",
	Long: `Pair a new door by having the gateway transmit its own radio code.

The motor must be in learn mode (press P button on the motor).
The gateway should be within ~20-30cm of the motor for pairing.
Normal operation works at much greater range after pairing.
In testing, the inherited port supported position feedback via HM_GET_TRANSITION.
Timeout: ~40 seconds.`,
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

		cli.Log.Infof("Transmitting gateway radio code (~40s timeout). Put the motor in learn mode (press P button)...")

		portId, err := bisecur.InheritPort(localMac, mac, host, port, token)
		if err != nil {
			cli.Log.Fatalf("%v", err)
			os.Exit(2)
		}

		cli.Log.Infof("Success. New port ID: %d", portId)
	},
}

func init() {
	portsCmd.AddCommand(portsInheritCmd)
}
