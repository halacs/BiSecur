package cmd

import (
	"bisecur/cli"
	"github.com/spf13/cobra"
)

var groupsDeleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete a gateway group",
	Long:    `Delete a gateway group`,
	PreRunE: preRunFuncs,
	Run: func(cmd *cobra.Command, args []string) {
		cli.Log.Fatalf("Not implemented yet")
	},
}

func init() {
	groupsCmd.AddCommand(groupsDeleteCmd)
}
