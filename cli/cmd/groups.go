package cmd

import (
	"bisecure/cli"
	"bisecure/sdk"
	"github.com/spf13/viper"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	groupsCmd := &cobra.Command{
		Use:    "groups",
		Short:  "Manages users defined in your Hörmann BiSecur gateway.",
		Long:   ``,
		PreRun: toggleDebug,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO implement query user list and rights, add and delete user, password change of an already existing user

			deviceMac := viper.GetString(ArgNameDeviceMac)
			host := viper.GetString(ArgNameHost)
			port := viper.GetInt(ArgNamePort)
			token := viper.GetUint32(ArgNameToken)

			mac, err := cli.ParesMacString(deviceMac)
			if err != nil {
				log.Fatalf("%v", err)
				os.Exit(1)
			}

			err = listGroups(localMac, mac, host, port, token)
			if err != nil {
				log.Fatalf("%v", err)
				os.Exit(2)
			}
		},
	}

	rootCmd.AddCommand(groupsCmd)
}

func listGroups(localMac [6]byte, mac [6]byte, host string, port int, token uint32) error {
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

	groups, err := client.GetGroups()
	if err != nil {
		return err
	}

	log.Infof("Groups: %s", groups.String())

	return nil
}
