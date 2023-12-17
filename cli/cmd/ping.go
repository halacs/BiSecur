package cmd

import (
	"bisecure/cli"
	"bisecure/sdk"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	var (
		count int
	)

	pingCmd := &cobra.Command{
		Use:   "ping",
		Short: "Check if your Hörmann BiSecur gateway is reachable or not.",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			mac, err := cli.ParesMacString(deviceMac)
			if err != nil {
				log.Fatalf("%v", err)
				os.Exit(1)
			}

			err = ping(localMac, mac, host, port, count)
			if err != nil {
				log.Fatalf("%v", err)
				os.Exit(2)
			}
		},
	}
	rootCmd.AddCommand(pingCmd)

	pingCmd.Flags().IntVar(&count, "count", 5, "Amount of the ping packages will be sent to the device")
}

func ping(localMac [6]byte, mac [6]byte, host string, port int, count int) error {
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

	err = client.Ping(count)
	if err != nil {
		return err
	}

	return nil
}
