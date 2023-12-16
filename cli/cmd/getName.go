package cmd

import (
	"bisecure/cli"
	"bisecure/sdk"
	"fmt"
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
				fmt.Printf("%v\n", err)
				os.Exit(1)
			}

			name, err := GetName(localMac, mac, host, port)
			if err != nil {
				fmt.Printf("%v\n", err)
				os.Exit(4)
			}
			fmt.Printf("Received name: %s\n", name)
		},
	}

	rootCmd.AddCommand(getNameCmd)
}

func GetName(localMac, mac [6]byte, host string, port int) (string, error) {
	client := sdk.NewClient(localMac, mac, host, port)
	err := client.Open()
	if err != nil {
		fmt.Printf("%v", err)
	}

	defer func() {
		err2 := client.Close()
		if err2 != nil {
			fmt.Printf("%v", err2)
			os.Exit(2)
		}
	}()

	name, err := client.GetName()
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(3)
	}

	return name, nil
}
