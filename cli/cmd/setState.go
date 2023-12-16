package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// setStateCmd represents the setState command
var setStateCmd = &cobra.Command{
	Use:   "set-state",
	Short: "Open or close a door connected to your Hörmann BiSecur gateway.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("setState called")
	},
}

func init() {
	rootCmd.AddCommand(setStateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setStateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setStateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
