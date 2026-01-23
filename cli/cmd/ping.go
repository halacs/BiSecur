package cmd

import (
	"halsecur/cli"
	"halsecur/cli/bisecur"
	"halsecur/cli/utils"
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"

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
				cli.Log.Fatalf("Invalid delay value: %d", delay)
				os.Exit(1)
			}

			mac, err := utils.ParesMacString(deviceMac)
			if err != nil {
				cli.Log.Fatalf("%v", err)
				os.Exit(1)
			}

			err = bisecur.Ping(localMac, mac, host, port, count, time.Duration(delay)*time.Millisecond, token)
			if err != nil {
				cli.Log.Fatalf("%v", err)
				os.Exit(2)
			}
		},
	}
	rootCmd.AddCommand(pingCmd)

	pingCmd.Flags().IntVarP(&count, "count", "c", 3, "Number of ping packages")
	pingCmd.Flags().IntVarP(&delay, "delay", "d", 1000, fmt.Sprintf("Miliseconds between ping packets. Must be at least %d", minDalayValue))
}
