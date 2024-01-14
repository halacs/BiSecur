package cmd

import (
	"bisecur/cli"
	"bisecur/sdk"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const (
	minDalayValue = 500
)

func init() {
	var (
		count int
		delay int
	)

	pingCmd := &cobra.Command{
		Use:     PingCmdName,
		Short:   "Check if your Hörmann BiSecur gateway is reachable or not.",
		Long:    ``,
		PreRunE: preRunFuncs,
		Run: func(cmd *cobra.Command, args []string) {
			deviceMac := viper.GetString(ArgNameDeviceMac)
			host := viper.GetString(ArgNameHost)
			port := viper.GetInt(ArgNamePort)
			token := viper.GetUint32(ArgNameToken)

			// input validation. Try not to do DOS attack against the gateway.
			if delay < minDalayValue {
				log.Fatalf("Invalid delay value: %d", delay)
				os.Exit(1)
			}

			mac, err := cli.ParesMacString(deviceMac)
			if err != nil {
				log.Fatalf("%v", err)
				os.Exit(1)
			}

			err = ping(localMac, mac, host, port, count, time.Duration(delay)*time.Millisecond, token)
			if err != nil {
				log.Fatalf("%v", err)
				os.Exit(2)
			}
		},
	}
	rootCmd.AddCommand(pingCmd)

	pingCmd.Flags().IntVarP(&count, "count", "c", 3, "Number of ping packages")
	pingCmd.Flags().IntVarP(&delay, "delay", "d", 1000, fmt.Sprintf("Miliseconds between ping packets. Must be at least %d", minDalayValue))
}

func ping(localMac [6]byte, mac [6]byte, host string, port int, count int, delay time.Duration, token uint32) error {
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

	received := 0
	for i := 0; i < count; i++ {
		sentTimestamp, receivedTimestamp, err := client.Ping()

		if err != nil {
			log.Errorf("%v", err)
			continue
		}

		received = received + 1
		rtt := receivedTimestamp - sentTimestamp
		log.Infof("Response %d of %d received in %d ms", received, count, rtt)

		if i < count {
			time.Sleep(delay)
		}
	}

	return nil
}
