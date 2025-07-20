package cmd

import (
	"bisecur/cli"
	"bisecur/cli/bisecur"
	"github.com/spf13/viper"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	logoutCmd := &cobra.Command{
		Use:     LogoutCmdName,
		Short:   "",
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

			err = bisecur.Logout(localMac, mac, host, port, token)
			if err != nil {
				cli.Log.Fatalf("%v", err)
				os.Exit(2)
			}

			// Clear token in persistent config file
			viper.Set(ArgNameToken, 0)
			err = viper.WriteConfig()
			if err != nil {
				cli.Log.Errorf("Failed to save new configuration. %v", err)
			}

			cli.Log.Infof("Success")
		},
	}

	rootCmd.AddCommand(logoutCmd)
}
