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
				log.Fatalf("%v", err)
				os.Exit(1)
			}

			userId, err := userAdd(localMac, mac, host, port, token, userName, password)
			if err != nil {
				log.Fatalf("%v", err)
				os.Exit(2)
			}

			log.Infof("User has been added. User ID: %d", userId)
		},
	}

	usersCmd.AddCommand(usersCreateCmd)

	usersCreateCmd.Flags().StringVar(&userName, ArgNameUsername, "", "name of the new user")
	usersCreateCmd.MarkFlagsOneRequired(ArgNameUsername)

	usersCreateCmd.Flags().StringVar(&password, ArgNameNewPassword, "", "password of the new user")
	usersCreateCmd.MarkFlagsOneRequired(ArgNameNewPassword)
}

func userAdd(localMac [6]byte, mac [6]byte, host string, port int, token uint32, userName string, password string) (byte, error) {
	client := sdk.NewClient(log, localMac, mac, host, port, token)
	err := client.Open()
	if err != nil {
		return 0, err
	}

	defer func() {
		err2 := client.Close()
		if err2 != nil {
			log.Errorf("%v", err)
		}
	}()

	var userId byte
	err = retry(func() error {
		var err2 error
		userId, err2 = client.AddUser(userName, password)
		return err2
	})

	if err != nil {
		return 0, err
	}

	return userId, nil
}
