package cmd

import (
	"github.com/spf13/cobra"
)

var groupsCmd = &cobra.Command{
	Use:     GroupsCmdName,
	Short:   "Manages doors defined in your Hörmann BiSecur gateway.",
	Long:    ``,
	PreRunE: preRunFuncs,
	Run:     nil,
}

func init() {
	rootCmd.AddCommand(groupsCmd)
}
