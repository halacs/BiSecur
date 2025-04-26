package cmd

import (
	"bisecur/cli"
	"bisecur/cli/bisecur"
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
				cli.Log.Fatalf("%v", err)
				os.Exit(2)
			}

			cli.Log.Infof("Successful login")
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

	token, err := bisecur.Login(localMac, mac, host, port, username, password)
	if err != nil {
		return err
	}

	cli.Log.Infof("Token: 0x%X", token)

	// Store token in persistent config
	viper.Set(ArgNameToken, token)
	viper.Set(ArgNameLastLoginTimeStamp, time.Now().UnixMicro())
	err = viper.WriteConfig()
	if err != nil {
		return fmt.Errorf("failed to save new configuration. %v", err)
	}

	return nil
}
