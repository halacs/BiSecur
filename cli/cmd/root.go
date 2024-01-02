package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	localMac = [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x09}
	log      *logrus.Logger
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "halsecur",
	Short:   "Application to manage your Hörmann BiSecur gateway without the central cloud directly on your LAN.",
	Long:    ``,
	PreRunE: preRunFuncs,
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

	var (
		host      string
		port      int
		token     uint32
		username  string
		password  string
		deviceMac string
		debug     bool
		autoLogin bool
	)

	rootCmd.PersistentFlags().Uint32Var(&token, ArgNameToken, 0x0, "Valid authentication token")
	rootCmd.PersistentFlags().StringVar(&username, ArgNameUsername, "", "Valid username")
	rootCmd.PersistentFlags().StringVar(&password, ArgNamePassword, "", "Valid password belongs to the given username")
	rootCmd.PersistentFlags().StringVar(&host, ArgNameHost, "", "IP or host name or the Hörmann BiSecure gateway")
	rootCmd.PersistentFlags().IntVar(&port, ArgNamePort, 4000, "")
	rootCmd.PersistentFlags().StringVar(&deviceMac, ArgNameDeviceMac, "", "MAC address of the Hörmann BiSecur gateway")
	rootCmd.PersistentFlags().BoolVar(&debug, ArgNameDebug, true, "debug log level")
	rootCmd.PersistentFlags().BoolVar(&autoLogin, ArgNameAutoLogin, true, "login automatically on demand")

	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	//viper.AddConfigPath("/etc/appname/")  // path to look for the config file in
	//viper.AddConfigPath("$HOME/.appname") // call multiple times to add many search paths
	viper.AddConfigPath(".") // optionally look for config in the working directory
	err := viper.BindPFlags(rootCmd.PersistentFlags())
	if err != nil {
		log.Fatalf("failed to bind PFlags. %v", err)
		os.Exit(1)
	}
	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if ok {
			log.Warning("Config file not found. Most probably you can ignore this message.")
		} else {
			log.Errorf("Failed to parse config file. %v", err)
		}
	}
}

func preRunFuncs(cmd *cobra.Command, args []string) error {
	toggleDebug(cmd, args)

	err := autoLogin(cmd, args)
	if err != nil {
		return fmt.Errorf("failed to auto login. %v", err)
	}

	return nil
}

func toggleDebug(cmd *cobra.Command, args []string) {
	debug := viper.GetBool(ArgNameDebug)
	if debug {
		log.SetLevel(logrus.DebugLevel)
	} else {
		log.SetLevel(logrus.InfoLevel)
	}
}

func newLogger() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		//ForceColors:   true,
		DisableColors: false,
		FullTimestamp: true,
	})
	log.SetReportCaller(false)

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
	return log
}
