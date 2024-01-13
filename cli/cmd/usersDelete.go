package cmd

import (
	"github.com/spf13/cobra"
)

var usersDeleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete a gateway user",
	Long:    `Delete a gateway user`,
	PreRunE: preRunFuncs,
	Run: func(cmd *cobra.Command, args []string) {
		log.Fatalf("Not implemented yet")
	},
}

func init() {
	usersCmd.AddCommand(usersDeleteCmd)
}
