package cmd

import (
	"bisecur/cli"
	"bisecur/sdk"
	"github.com/spf13/viper"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	var devicePort int

	setStateCmd := &cobra.Command{
		Use:     SetStateCmdUse,
		Short:   "Open or close a door connected to your Hörmann BiSecur gateway.",
		Long:    ``,
		PreRunE: preRunFuncs,
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

			err = setStatus(localMac, mac, host, port, byte(devicePort), token)
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

func setStatus(localMac [6]byte, mac [6]byte, host string, port int, devicePort byte, token uint32) error {
	client := sdk.NewClient(log, localMac, mac, host, port, token)
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

	err = retry(func() error {
		err2 := client.SetState(devicePort)
		return err2
	})

	if err != nil {
		return err
	}

	log.Infof("Done")

	return nil
}
