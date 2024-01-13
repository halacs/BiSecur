package cmd

import (
	"github.com/spf13/cobra"
)

var usersCreateCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a new gateway user",
	Long:    `Create a new gateway user`,
	PreRunE: preRunFuncs,
	Run: func(cmd *cobra.Command, args []string) {
		log.Fatalf("Not implemented yet")
	},
}

func init() {
	usersCmd.AddCommand(usersCreateCmd)
}
