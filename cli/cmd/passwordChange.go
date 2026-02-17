package cmd

import (
	"halsecur/cli"
	"halsecur/cli/bisecur"
	"halsecur/cli/utils"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

			err = bisecur.UserPasswordChange(localMac, mac, host, port, token, userIdByte, newPassword)
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
