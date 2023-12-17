package cmd

import (
	"bisecure/cli"
	"bisecure/sdk"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	var devicePort int

	setStateCmd := &cobra.Command{
		Use:   "set-state",
		Short: "Open or close a door connected to your Hörmann BiSecur gateway.",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			mac, err := cli.ParesMacString(deviceMac)
			if err != nil {
				fmt.Printf("%v\n", err)
				os.Exit(1)
			}

			err = setStatus(localMac, mac, host, port, username, password, byte(devicePort))
			if err != nil {
				log.Fatalf("%v", err)
				os.Exit(2)
			}
		},
	}

	rootCmd.AddCommand(setStateCmd)

	setStateCmd.Flags().IntVar(&devicePort, devicePortName, 0, "Port number of the door")
	setStateCmd.MarkFlagsOneRequired(devicePortName)
}

func setStatus(localMac [6]byte, mac [6]byte, host string, port int, username string, password string, devicePort byte) error {
	client := sdk.NewClient(log, localMac, mac, host, port)
	err := client.Open()
	if err != nil {
		return err
	}

	defer func() {
		err2 := client.Close()
		if err2 != nil {
			log.Errorf("%v", err2)
		}
	}()

	err = client.Login(username, password)
	if err != nil {
		return err
	}

	log.Infof("Logged in successfully.")

	err = client.SetState(devicePort)
	if err != nil {
		return err
	}

	log.Infof("Done")

	return nil
}
