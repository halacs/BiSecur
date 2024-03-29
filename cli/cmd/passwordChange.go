package cmd

import (
	"bisecur/cli"
	"bisecur/sdk"
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
				log.Fatalf("%v", err)
				os.Exit(1)
			}

			err = userPasswordChange(localMac, mac, host, port, token, byte(userId), newPassword)
			if err != nil {
				log.Fatalf("%v", err)
				os.Exit(2)
			}

			log.Infof("Password has been changed")
		},
	}

	usersCmd.AddCommand(passwordChangeCmd)

	passwordChangeCmd.Flags().IntVar(&userId, ArgNameUserId, 0, "ID of the user")
	passwordChangeCmd.MarkFlagsOneRequired(ArgNameUserId)

	passwordChangeCmd.Flags().StringVar(&newPassword, ArgNameNewPassword, "", "new password")
	passwordChangeCmd.MarkFlagsOneRequired(ArgNameNewPassword)
}

func userPasswordChange(localMac [6]byte, mac [6]byte, host string, port int, token uint32, userId byte, newPassword string) error {
	client := sdk.NewClient(log, localMac, mac, host, port, token)
	err := client.Open()
	if err != nil {
		return err
	}

	defer func() {
		err2 := client.Close()
		if err2 != nil {
			log.Errorf("%v", err)
		}
	}()

	err = retry(func() error {
		err2 := client.PasswordChange(userId, newPassword)
		return err2
	})

	if err != nil {
		return err
	}

	return nil
}
