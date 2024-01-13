package cmd

import (
	"github.com/spf13/cobra"
)

// passwordChangeCmd represents the passwordChange command
var passwordChangeCmd = &cobra.Command{
	Use:   "passwordChange",
	Short: "Change password of a gateway user",
	Long:  `Change password of a gateway user`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Fatalf("Not implemented yet")
	},
}

func init() {
	usersCmd.AddCommand(passwordChangeCmd)
}
