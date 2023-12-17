package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
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
	log       *logrus.Logger
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
	log = newLogger()

	rootCmd.PersistentFlags().StringVar(&username, "username", "", "Valid username")
	rootCmd.PersistentFlags().StringVar(&password, "password", "", "Valid password belongs to the given username")
	rootCmd.PersistentFlags().StringVar(&host, "host", "", "IP or host name or the Hörmann BiSecure gateway")
	rootCmd.PersistentFlags().IntVar(&port, "port", 4000, "")
	rootCmd.PersistentFlags().StringVar(&deviceMac, "mac", "", "MAC address of the Hörmann BiSecur gateway")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug log level")
}

func newLogger() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		//ForceColors:   true,
		DisableColors: false,
		FullTimestamp: true,
	})
	log.SetReportCaller(false)
	if debug {
		log.SetLevel(logrus.DebugLevel)
	} else {
		log.SetLevel(logrus.InfoLevel)
	}
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
	return log
}
