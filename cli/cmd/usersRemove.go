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
		userId int
	)

	usersDeleteCmd := &cobra.Command{
		Use:     "remove",
		Short:   "Delete a gateway user",
		Long:    `Delete a gateway user`,
		PreRunE: preRunFuncs,
		Run: func(cmd *cobra.Command, args []string) {
			deviceMac := viper.GetString(ArgNameDeviceMac)
			host := viper.GetString(ArgNameHost)
			port := viper.GetInt(ArgNamePort)
			token := viper.GetUint32(ArgNameToken)
			userId := viper.GetInt(ArgNameUserId)

			mac, err := cli.ParesMacString(deviceMac)
			if err != nil {
				log.Fatalf("%v", err)
				os.Exit(1)
			}

			err = userRemove(localMac, mac, host, port, token, byte(userId))
			if err != nil {
				log.Fatalf("%v", err)
				os.Exit(2)
			}

			log.Infof("Password has been removed")
		},
	}

	usersCmd.AddCommand(usersDeleteCmd)

	usersDeleteCmd.Flags().IntVar(&userId, ArgNameUserId, 0, "ID of the user to be deleted")
	usersDeleteCmd.MarkFlagsOneRequired(ArgNameUserId)
}

func userRemove(localMac [6]byte, mac [6]byte, host string, port int, token uint32, userId byte) error {
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
		err2 := client.RemoveUser(userId)
		return err2
	})

	if err != nil {
		return err
	}

	return nil
}
