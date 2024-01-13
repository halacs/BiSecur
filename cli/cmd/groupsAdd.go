package cmd

import (
	"github.com/spf13/cobra"
)

var groupsCreateCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a new gateway group",
	Long:    `Create a new gateway group`,
	PreRunE: preRunFuncs,
	Run: func(cmd *cobra.Command, args []string) {
		log.Fatalf("Not implemented yet")
	},
}

func init() {
	groupsCmd.AddCommand(groupsCreateCmd)
}
