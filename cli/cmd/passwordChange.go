package cmd

import (
	"bisecur/cli"
	"bisecur/cli/bisecur"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

func init() {
	var (
		userId      int
		newPassword string
	)

	passwordChangeCmd := &cobra.Command{
		Use:   "password-change",
		Short: "Change password of a gateway user",
		Long:  `Change password of a gateway user`,
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

			err = bisecur.UserPasswordChange(localMac, mac, host, port, token, byte(userId), newPassword)
			if err != nil {
				cli.Log.Fatalf("%v", err)
				os.Exit(2)
			}

			cli.Log.Infof("Password has been changed")
		},
	}

	usersCmd.AddCommand(passwordChangeCmd)

	passwordChangeCmd.Flags().IntVar(&userId, ArgNameUserId, 0, "ID of the user")
	passwordChangeCmd.MarkFlagsOneRequired(ArgNameUserId)

	passwordChangeCmd.Flags().StringVar(&newPassword, ArgNameNewPassword, "", "new password")
	passwordChangeCmd.MarkFlagsOneRequired(ArgNameNewPassword)
}
