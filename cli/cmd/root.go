package cmd

import (
	"bisecur/cli"
	"bisecur/logger"
	"fmt"
	"github.com/spf13/viper"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	localMac = [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x09}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "halsecur",
	Short:   "Application to manage your Hörmann BiSecur gateway without the central cloud directly on your LAN.",
	Long:    ``,
	PreRunE: preRunFuncs,
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
	cli.Log = logger.NewLogger()

	var (
		host      string
		port      int
		token     uint32
		username  string
		password  string
		deviceMac string
	)

	rootCmd.PersistentFlags().Uint32Var(&token, ArgNameToken, 0x0, "Valid authentication token")
	rootCmd.PersistentFlags().StringVar(&username, ArgNameUsername, "", "Valid username")
	rootCmd.PersistentFlags().StringVar(&password, ArgNamePassword, "", "Valid password belongs to the given username")
	rootCmd.PersistentFlags().StringVar(&host, ArgNameHost, "", "IP or host name or the Hörmann BiSecure gateway")
	rootCmd.PersistentFlags().IntVar(&port, ArgNamePort, 4000, "")
	rootCmd.PersistentFlags().StringVar(&deviceMac, ArgNameDeviceMac, "", "MAC address of the Hörmann BiSecur gateway")
	rootCmd.PersistentFlags().Bool(ArgNameDebug, true, "debug Log level")
	rootCmd.PersistentFlags().Bool(ArgNameJsonOutput, false, "use json logging format instead of human readable")
	rootCmd.PersistentFlags().Bool(ArgNameAutoLogin, true, "login automatically on demand")

	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	//viper.AddConfigPath("/etc/appname/")  // path to look for the config file in
	//viper.AddConfigPath("$HOME/.appname") // call multiple times to add many search paths
	viper.AddConfigPath(".") // optionally look for config in the working directory
	err := viper.BindPFlags(rootCmd.PersistentFlags())
	if err != nil {
		cli.Log.Fatalf("failed to bind PFlags. %v", err)
		os.Exit(1)
	}
	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if ok {
			cli.Log.Warning("Config file not found. Most probably you can ignore this message.")
		} else {
			cli.Log.Errorf("Failed to parse config file. %v", err)
		}
	}
}

func preRunFuncs(cmd *cobra.Command, args []string) error {
	jsonOutput(cmd, args)
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
		cli.Log.SetLevel(logrus.DebugLevel)
	} else {
		cli.Log.SetLevel(logrus.InfoLevel)
	}
}

func jsonOutput(cmd *cobra.Command, args []string) {
	jsonOutput := viper.GetBool(ArgNameJsonOutput)
	if jsonOutput {
		jsonFormatter := &logrus.JSONFormatter{
			PrettyPrint: true,
		}

		cli.Log.SetFormatter(jsonFormatter)
	}
}
