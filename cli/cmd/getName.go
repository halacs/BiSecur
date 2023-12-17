package cmd

import (
	"bisecure/cli"
	"bisecure/sdk"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	getNameCmd := &cobra.Command{
		Use:   "get-name",
		Short: "Queries the name of the Hörmann BiSecur gateway",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			mac, err := cli.ParesMacString(deviceMac)
			if err != nil {
				log.Fatalf("%v", err)
				os.Exit(1)
			}

			err = GetName(localMac, mac, host, port)
			if err != nil {
				log.Fatalf("%v", err)
				os.Exit(4)
			}
		},
	}

	rootCmd.AddCommand(getNameCmd)
}

func GetName(localMac, mac [6]byte, host string, port int) error {
	client := sdk.NewClient(log, localMac, mac, host, port)
	err := client.Open()
	if err != nil {
		return err
	}

	defer func() {
		err2 := client.Close()
		if err2 != nil {
			log.Fatalf("%v", err2)
		}
	}()

	name, err := client.GetName()
	if err != nil {
		return err
	}

	log.Infof("Received name: %s", name)

	return nil
}
