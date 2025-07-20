package cmd

import (
	"bisecur/cli"
	"bisecur/cli/bisecur"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"time"
)

func loginRequired(name string) bool {
	loginRequiredCommands := []string{StatusCmdUse, SetStateCmdUse}

	for _, i := range loginRequiredCommands {
		if name == i {
			return true
		}
	}

	return false
}

func autoLogin(cmd *cobra.Command, args []string) error {
	if !loginRequired(cmd.Use) {
		cli.Log.Debugf("Login not required. Don't need to auto login.")
		return nil
	}

	autoLogin := viper.GetBool(ArgNameAutoLogin)
	if !autoLogin {
		cli.Log.Debugf("Auto login is disabled.")
		return nil
	}

	lastLoginTimeStamp := viper.GetInt64(ArgNameLastLoginTimeStamp)
	t := time.UnixMicro(lastLoginTimeStamp)

	if t.Add(bisecur.TokenExpirationTime).Before(time.Now()) {
		cli.Log.Infof("Token expired. Logging in...")
		err := loginCmdFunc()
		if err != nil {
			return fmt.Errorf("failed to auto login. %v", err)
		}
	}

	cli.Log.Debugf("Token is still valid.")
	return nil
}
