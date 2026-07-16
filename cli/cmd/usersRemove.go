package cmd

import (
	"halsecur/cli"
	"halsecur/cli/bisecur/users"
	"halsecur/cli/utils"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	var (
		userId int
	)

	usersDeleteCmd := &cobra.Command{
		Use:     UsersRemoveCmdName,
		Short:   "Delete a gateway user",
		Long:    `Delete a gateway user`,
		PreRunE: preRunFuncs,
		Run: func(cmd *cobra.Command, args []string) {
			deviceMac := viper.GetString(ArgNameDeviceMac)
			host := viper.GetString(ArgNameHost)
			port := viper.GetInt(ArgNamePort)
			token := viper.GetUint32(ArgNameToken)
			userId := viper.GetInt(ArgNameUserId)

			mac, err := utils.ParesMacString(deviceMac)
			if err != nil {
				cli.Log.Fatalf("%v", err)
				os.Exit(1)
			}

			userIdByte, err := utils.SafeLen(userId)
			if err != nil {
				cli.Log.Fatalf("too big userId value. %v", err)
				os.Exit(2)
			}
			err = users.UserRemove(localMac, mac, host, port, token, userIdByte)
			if err != nil {
				cli.Log.Fatalf("%v", err)
				os.Exit(2)
			}

			cli.Log.Infof("Password has been removed")
		},
	}

	usersCmd.AddCommand(usersDeleteCmd)

	usersDeleteCmd.Flags().IntVar(&userId, ArgNameUserId, 0, "ID of the user to be deleted")
	usersDeleteCmd.MarkFlagsOneRequired(ArgNameUserId)
}
