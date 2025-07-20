package cmd

import (
	"bisecur/cli"
	"bisecur/cli/bisecur"
	"github.com/spf13/viper"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	getNameCmd := &cobra.Command{
		Use:     "get-name",
		Short:   "Queries the name of the Hörmann BiSecur gateway",
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

			name, err := bisecur.GetName(localMac, mac, host, port, token)
			if err != nil {
				cli.Log.Fatalf("%v", err)
				os.Exit(4)
			}

			cli.Log.WithField("name", name).Infof("Success")
		},
	}

	rootCmd.AddCommand(getNameCmd)
}
