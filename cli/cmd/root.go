package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	host      string
	port      int
	username  string
	password  string
	deviceMac string
	localMac  = [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x09}
	debug     bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "halsecur",
	Short: "Application to manage your Hörmann BiSecur gateway without the central cloud directly on your LAN.",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cli.yaml)")

	rootCmd.PersistentFlags().StringVar(&username, "username", "", "Valid username")
	rootCmd.PersistentFlags().StringVar(&password, "password", "", "Valid password belongs to the given username")
	rootCmd.PersistentFlags().StringVar(&host, "host", "", "IP or host name or the Hörmann BiSecure gateway")
	rootCmd.PersistentFlags().IntVar(&port, "port", 4000, "")
	rootCmd.PersistentFlags().StringVar(&deviceMac, "mac", "", "MAC address of the Hörmann BiSecur gateway")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug log level")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().String("username", "", "")
}