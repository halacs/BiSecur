package cmd

import (
	"github.com/spf13/cobra"
)

var usersCmd = &cobra.Command{
	Use:     UsersCmdName,
	Short:   "Manages users defined in your Hörmann BiSecur gateway.",
	Long:    ``,
	PreRunE: preRunFuncs,
	Run:     nil,
}

func init() {
	rootCmd.AddCommand(usersCmd)
}
