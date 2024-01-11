package cmd

import (
	"bisecur/cli"
	"bisecur/sdk"
	"fmt"
	"github.com/spf13/viper"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	logoutCmd := &cobra.Command{
		Use:     LogoutCmdName,
		Short:   "",
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

			err = logout(localMac, mac, host, port, token)
			if err != nil {
				log.Fatalf("%v", err)
				os.Exit(2)
			}

			// Clear token in persistent config file
			viper.Set(ArgNameToken, 0)
			err = viper.WriteConfig()
			if err != nil {
				log.Errorf("Failed to save new configuration. %v", err)
			}
		},
	}

	rootCmd.AddCommand(logoutCmd)
}

func logout(localMac [6]byte, mac [6]byte, host string, port int, token uint32) error {
	if token == 0 {
		return fmt.Errorf("invalid token value: 0x%X", token)
	}

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

	client.SetToken(token)

	err = client.Logout()
	if err != nil {
		return err
	}

	return nil
}
