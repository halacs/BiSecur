package cmd

import (
	"bisecure/sdk"
	"context"
	"fmt"
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
				fmt.Printf("%v\n", err)
				os.Exit(2)
			}
		},
	}

	rootCmd.AddCommand(discoverCmd)

	discoverCmd.Flags().IntVar(&discoveryTime, "time", 10, "Time in second while device discovery happens")
}

func discover(discoveryTime int) error {
	ctx := context.Background()
	discovery := sdk.NewDiscovery(ctx, func(gateway sdk.Gateway) {
		fmt.Printf("Response received: %+v\n", gateway)
	})

	fmt.Printf("Start discovery...\n")
	err := discovery.Start()
	if err != nil {
		return err
	}

	fmt.Printf("Waiting for responses...\n")
	time.Sleep(time.Second * time.Duration(discoveryTime))

	list := discovery.GetList()
	fmt.Printf("list: %+v\n", list)

	fmt.Printf("Stop disovery...\n")
	err = discovery.Stop()
	if err != nil {
		return nil
	}

	fmt.Printf("Terminated\n")
	return nil
}
