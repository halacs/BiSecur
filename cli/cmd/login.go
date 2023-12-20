package cmd

import (
	"bisecure/cli"
	"bisecure/sdk"
	"github.com/spf13/viper"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	loginCmd := &cobra.Command{
		Use:    "login",
		Short:  "",
		Long:   ``,
		PreRun: toggleDebug,
		Run: func(cmd *cobra.Command, args []string) {
			deviceMac := viper.GetString(ArgNameDeviceMac)
			host := viper.GetString(ArgNameHost)
			port := viper.GetInt(ArgNamePort)
			username := viper.GetString(ArgNameUsername)
			password := viper.GetString(ArgNamePassword)

			mac, err := cli.ParesMacString(deviceMac)
			if err != nil {
				log.Fatalf("%v", err)
				os.Exit(1)
			}

			token, err := login(localMac, mac, host, port, username, password)
			if err != nil {
				log.Fatalf("%v", err)
				os.Exit(2)
			}

			log.Infof("Token: 0x%X", token)

			// Store token in persistent config
			viper.Set(ArgNameToken, token)
			err = viper.WriteConfig()
			if err != nil {
				log.Errorf("Failed to save new configuration. %v", err)
			}
		},
	}

	rootCmd.AddCommand(loginCmd)
}

func login(localMac [6]byte, mac [6]byte, host string, port int, username string, password string) (uint32, error) {
	client := sdk.NewClient(log, localMac, mac, host, port, 0)
	err := client.Open()
	if err != nil {
		return 0, err
	}

	defer func() {
		err2 := client.Close()
		if err2 != nil {
			log.Errorf("%v", err)
		}
	}()

	err = client.Login(username, password)
	if err != nil {
		return 0, err
	}

	return client.GetToken(), nil
}
