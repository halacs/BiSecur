package cmd

import (
	"github.com/spf13/cobra"
)

var portsCmd = &cobra.Command{
	Use:     PortsCmdName,
	Short:   "Manages ports (paired devices) on your Hörmann BiSecur gateway.",
	Long:    ``,
	PreRunE: preRunFuncs,
	Run:     nil,
}

func init() {
	rootCmd.AddCommand(portsCmd)
}
