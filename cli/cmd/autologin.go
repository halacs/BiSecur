package cmd

import (
	"fmt"
	"halsecur/cli"
	"halsecur/cli/bisecur"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func loginRequired(name string) bool {
	loginRequiredCommands := []string{StatusCmdUse, SetStateCmdUse, AddCmdUse, InheritCmdUse, RemoveCmdUse}

	for _, i := range loginRequiredCommands {
		if name == i {
			return true
		}
	}

	return false
}

func autoLogin(cmd *cobra.Command, args []string) error {
	if !loginRequired(cmd.Use) {
		cli.Log.Debugf("Login not required. Don't need to auto login for this command.")
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
	} else {
		cli.Log.Debugf("Token is still valid.")
	}

	return nil
}
