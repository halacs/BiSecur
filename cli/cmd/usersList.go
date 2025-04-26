package cmd

import (
	"bisecur/cli"
	"bisecur/cli/bisecur/users"
	"github.com/spf13/viper"
	"os"

	"github.com/spf13/cobra"
)

var usersListCmd = &cobra.Command{
	Use:     "groups",
	Short:   "List current gateway users",
	Long:    `List current gateway users`,
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

		err = users.ListUsers(localMac, mac, host, port, token)
		if err != nil {
			cli.Log.Fatalf("%v", err)
			os.Exit(2)
		}
	},
}

func init() {
	usersCmd.AddCommand(usersListCmd)
}
