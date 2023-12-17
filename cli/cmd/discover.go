package cmd

import (
	"bisecure/sdk"
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
		Use:   "discover",
		Short: "Discover Hörmann BiSecur gateways on the local network",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			err := discover(discoveryTime)
			if err != nil {
				log.Fatalf("%v", err)
				os.Exit(2)
			}
		},
	}

	rootCmd.AddCommand(discoverCmd)

	discoverCmd.Flags().IntVar(&discoveryTime, "time", 10, "Time in second while device discovery happens")
}

func discover(discoveryTime int) error {
	ctx := context.Background()
	discovery := sdk.NewDiscovery(ctx, log, func(gateway sdk.Gateway) {
		log.Infof("Response received: %+v", gateway)
	})

	log.Infof("Start discovery...")
	err := discovery.Start()
	if err != nil {
		return err
	}

	log.Infof("Waiting for responses...")
	time.Sleep(time.Second * time.Duration(discoveryTime))

	list := discovery.GetList()
	log.Infof("list: %+v", list)

	log.Infof("Stop disovery...")
	err = discovery.Stop()
	if err != nil {
		return nil
	}

	log.Infof("Terminated")
	return nil
}
