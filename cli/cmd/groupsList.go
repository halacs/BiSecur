package cmd

import (
	"halsecur/cli"
	"halsecur/cli/bisecur/groups"
	"halsecur/cli/utils"
	"os"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var groupsListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List current gateway groups",
	Long:    `List current gateway groups`,
	PreRunE: preRunFuncs,
	Run: func(cmd *cobra.Command, args []string) {
		deviceMac := viper.GetString(ArgNameDeviceMac)
		host := viper.GetString(ArgNameHost)
		port := viper.GetInt(ArgNamePort)
		token := viper.GetUint32(ArgNameToken)

		mac, err := utils.ParesMacString(deviceMac)
		if err != nil {
			cli.Log.Fatalf("%v", err)
			os.Exit(1)
		}

		grps, err := groups.ListGroups(localMac, mac, host, port, token)
		if err != nil {
			cli.Log.Fatalf("%v", err)
			os.Exit(2)
		}

		cli.Log.Infof("%s", grps.String())
	},
}

func init() {
	groupsCmd.AddCommand(groupsListCmd)
}
