package cmd

import (
	"bisecure/cli"
	"bisecure/sdk"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	groupsCmd := &cobra.Command{
		Use:   "groups",
		Short: "Manages users defined in your Hörmann BiSecur gateway.",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO implement query user list and rights, add and delete user, password change of an already existing user

			mac, err := cli.ParesMacString(deviceMac)
			if err != nil {
				log.Fatalf("%v", err)
				os.Exit(1)
			}

			err = listGroups(localMac, mac, host, port)
			if err != nil {
				log.Fatalf("%v", err)
				os.Exit(2)
			}
		},
	}

	rootCmd.AddCommand(groupsCmd)
}

func listGroups(localMac [6]byte, mac [6]byte, host string, port int) error {
	client := sdk.NewClient(log, localMac, mac, host, port)
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
