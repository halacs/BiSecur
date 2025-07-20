package cmd

import (
	"bisecur/cli"
	"bisecur/sdk"
	"context"
	"os"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	var (
		discoveryTime = 20
	)

	discoverCmd := &cobra.Command{
		Use:     "discover",
		Short:   "Discover Hörmann BiSecur gateways on the local network",
		Long:    ``,
		PreRunE: preRunFuncs,
		Run: func(cmd *cobra.Command, args []string) {
			err := discover(discoveryTime)
			if err != nil {
				cli.Log.Fatalf("%v", err)
				os.Exit(2)
			}
		},
	}

	rootCmd.AddCommand(discoverCmd)

	discoverCmd.Flags().IntVar(&discoveryTime, "time", 10, "Time in second while device discovery happens")
}

func discover(discoveryTime int) error {
	ctx := context.Background()
	discovery := sdk.NewDiscovery(ctx, cli.Log, func(gateway sdk.Gateway) {
		cli.Log.Infof("Response received: %+v", gateway)
	})

	cli.Log.Infof("Start discovery...")
	err := discovery.Start()
	if err != nil {
		return err
	}

	cli.Log.Infof("Waiting for responses...")
	time.Sleep(time.Second * time.Duration(discoveryTime))

	list := discovery.GetList()
	cli.Log.Infof("groups: %+v", list)

	cli.Log.Infof("Stop disovery...")
	err = discovery.Stop()
	if err != nil {
		return nil
	}

	cli.Log.Infof("Terminated")
	return nil
}
