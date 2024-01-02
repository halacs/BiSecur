package cmd

import (
	"bisecure/cli"
	"bisecure/sdk"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"time"
)

func init() {
	loginCmd := &cobra.Command{
		Use:     LoginCmdName,
		Short:   "",
		Long:    ``,
		PreRunE: preRunFuncs,
		Run: func(cmd *cobra.Command, args []string) {
			err := loginCmdFunc()
			if err != nil {
				log.Fatalf("%v", err)
				os.Exit(2)
			}
		},
	}

	rootCmd.AddCommand(loginCmd)
}

func loginCmdFunc() error {
	deviceMac := viper.GetString(ArgNameDeviceMac)
	host := viper.GetString(ArgNameHost)
	port := viper.GetInt(ArgNamePort)
	username := viper.GetString(ArgNameUsername)
	password := viper.GetString(ArgNamePassword)

	mac, err := cli.ParesMacString(deviceMac)
	if err != nil {
		return err
	}

	token, err := login(localMac, mac, host, port, username, password)
	if err != nil {
		return err
	}

	log.Infof("Token: 0x%X", token)

	// Store token in persistent config
	viper.Set(ArgNameToken, token)
	viper.Set(ArgNameLastLoginTimeStamp, time.Now().UnixMicro())
	err = viper.WriteConfig()
	if err != nil {
		return fmt.Errorf("failed to save new configuration. %v", err)
	}

	return nil
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
