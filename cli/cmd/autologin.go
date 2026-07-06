package cmd

import (
	"fmt"
	"halsecur/cli"

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

// autoLogin logs in before a command that needs a session (status / set-state). It used to skip
// the login while a cached token looked fresh (based on a fixed expiry timer), but the gateway
// idle-expires tokens far sooner than any fixed timer, so for these one-shot CLI commands we just
// log in each time when auto-login is enabled. (The long-running `ha` daemon re-logs-in reactively
// on PERMISSION_DENIED instead.)
func autoLogin(cmd *cobra.Command, args []string) error {
	if !loginRequired(cmd.Use) {
		cli.Log.Debugf("Login not required. Don't need to auto login for this command.")
		return nil
	}

	if !viper.GetBool(ArgNameAutoLogin) {
		cli.Log.Debugf("Auto login is disabled.")
		return nil
	}

	cli.Log.Infof("Logging in...")
	if err := loginCmdFunc(); err != nil {
		return fmt.Errorf("failed to auto login. %v", err)
	}

	return nil
}
