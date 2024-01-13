package cmd

import (
	"bisecur/cli"
	"bisecur/sdk"
	"github.com/spf13/viper"
	"os"

	"github.com/spf13/cobra"
)

var usersListCmd = &cobra.Command{
	Use:     "list",
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
			log.Fatalf("%v", err)
			os.Exit(1)
		}

		err = listUsers(localMac, mac, host, port, token)
		if err != nil {
			log.Fatalf("%v", err)
			os.Exit(2)
		}
	},
}

func init() {
	usersCmd.AddCommand(usersListCmd)
}

func listUsers(localMac [6]byte, mac [6]byte, host string, port int, token uint32) error {
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

	var users *sdk.Users
	err = retry(func() error {
		var err2 error
		users, err2 = client.GetUsers()
		return err2
	})

	if err != nil {
		return err
	}

	log.Infof("Users: %s", users.String())

	return nil
}
