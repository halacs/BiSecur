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

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Queries the status (open/closed/etc) of your door.",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			mac, err := cli.ParesMacString(deviceMac)
			if err != nil {
				log.Fatalf("%v", err)
				os.Exit(1)
			}

			err = getStatus(localMac, mac, host, port, username, password, byte(devicePort))
			if err != nil {
				log.Fatalf("%v", err)
				os.Exit(2)
			}
		},
	}

	rootCmd.AddCommand(statusCmd)

	statusCmd.Flags().IntVar(&devicePort, devicePortName, 0, "Port number of the door")
	statusCmd.MarkFlagsOneRequired(devicePortName)
}

func getStatus(localMac [6]byte, mac [6]byte, host string, port int, username string, password string, devicePort byte) error {
	client := sdk.NewClient(log, localMac, mac, host, port)
	err := client.Open()
	if err != nil {
		fmt.Printf("%v", err)
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

	status, err := client.GetTransition(devicePort)
	if err != nil {
		return err
	}

	log.Infof("Transition: %+v", status)

	return nil
}
