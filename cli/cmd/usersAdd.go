package cmd

import (
	"bisecur/cli"
	"bisecur/cli/bisecur/users"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

func init() {
	var (
		userName string
		password string
	)

	usersCreateCmd := &cobra.Command{
		Use:     "add",
		Short:   "Create a new gateway user",
		Long:    `Create a new gateway user`,
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

			userId, err := users.UserAdd(localMac, mac, host, port, token, userName, password)
			if err != nil {
				cli.Log.Fatalf("%v", err)
				os.Exit(2)
			}

			cli.Log.Infof("User has been added. User ID: %d", userId)
		},
	}

	usersCmd.AddCommand(usersCreateCmd)

	usersCreateCmd.Flags().StringVar(&userName, ArgNameUsername, "", "name of the new user")
	usersCreateCmd.MarkFlagsOneRequired(ArgNameUsername)

	usersCreateCmd.Flags().StringVar(&password, ArgNameNewPassword, "", "password of the new user")
	usersCreateCmd.MarkFlagsOneRequired(ArgNameNewPassword)
}
