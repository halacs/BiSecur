package cmd

import (
	"bisecure/cli"
	"bisecure/sdk"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	var (
		count = 4
	)

	pingCmd := &cobra.Command{
		Use:   "ping",
		Short: "Check if your Hörmann BiSecur gateway is reachable or not.",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			mac, err := cli.ParesMacString(deviceMac)
			if err != nil {
				fmt.Printf("%v\n", err)
				os.Exit(1)
			}

			err = ping(localMac, mac, host, port, count)
			if err != nil {
				fmt.Printf("%v\n", err)
				os.Exit(2)
			}
		},
	}
	rootCmd.AddCommand(pingCmd)

	pingCmd.Flags().IntVar(&count, "count", 5, "Amount of the ping packages will be sent to the device")
}

func ping(localMac, mac [6]byte, host string, port int, count int) error {
	client := sdk.NewClient(localMac, mac, host, port)
	err := client.Open()
	if err != nil {
		return err
	}

	defer func() {
		err2 := client.Close()
		if err2 != nil {
			fmt.Printf("%v", err) // TODO add log message
		}
	}()

	err = client.Ping(count)
	if err != nil {
		return err
	}

	return nil
}
