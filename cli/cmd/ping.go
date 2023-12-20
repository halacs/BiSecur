package cmd

import (
	"bisecure/cli"
	"bisecure/sdk"
	"github.com/spf13/viper"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	var (
		count int
	)

	pingCmd := &cobra.Command{
		Use:    "ping",
		Short:  "Check if your Hörmann BiSecur gateway is reachable or not.",
		Long:   ``,
		PreRun: toggleDebug,
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

			err = ping(localMac, mac, host, port, count, token)
			if err != nil {
				log.Fatalf("%v", err)
				os.Exit(2)
			}
		},
	}
	rootCmd.AddCommand(pingCmd)

	pingCmd.Flags().IntVar(&count, "count", 3, "Amount of the ping packages will be sent to the device")
}

func ping(localMac [6]byte, mac [6]byte, host string, port int, count int, token uint32) error {
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

	err = client.Ping(count)
	if err != nil {
		return err
	}

	return nil
}
