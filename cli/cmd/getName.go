package cmd

import (
	"bisecur/cli"
	"bisecur/sdk"
	"github.com/spf13/viper"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	getNameCmd := &cobra.Command{
		Use:     "get-name",
		Short:   "Queries the name of the Hörmann BiSecur gateway",
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

			err = GetName(localMac, mac, host, port, token)
			if err != nil {
				log.Fatalf("%v", err)
				os.Exit(4)
			}
		},
	}

	rootCmd.AddCommand(getNameCmd)
}

func GetName(localMac, mac [6]byte, host string, port int, token uint32) error {
	client := sdk.NewClient(log, localMac, mac, host, port, token)
	err := client.Open()
	if err != nil {
		return err
	}

	defer func() {
		err2 := client.Close()
		if err2 != nil {
			log.Fatalf("%v", err2)
		}
	}()

	var name string
	err = retry(func() error {
		var err2 error
		name, err2 = client.GetName()
		return err2
	})

	if err != nil {
		return err
	}

	log.WithField("name", name).Infof("Success")

	return nil
}
