package cmd

import (
	"bisecur/cli"
	"bisecur/cli/homeAssistant"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"time"
)

func init() {
	var (
		mqttServerName         string
		mqttServerPort         int
		mqttServerTls          bool
		mqttServerTlsValidaton bool
		mqttBaseTopic          string
		mqttDeviceName         string
		mqttUserName           string
		mqttPassword           string
		mqttTelePeriod         time.Duration
		devicePort             int
	)

	haCmd := &cobra.Command{
		Use:     HomeAssistantCmdName,
		Short:   "Start MQTT client compatible with Home assistant auto discovery",
		Long:    ``,
		PreRunE: preRunFuncs,
		Run: func(cmd *cobra.Command, args []string) {
			deviceMac := viper.GetString(ArgNameDeviceMac)
			host := viper.GetString(ArgNameHost)
			port := viper.GetInt(ArgNamePort)
			token := viper.GetUint32(ArgNameToken)
			username := viper.GetString(ArgNameUsername)
			password := viper.GetString(ArgNamePassword)

			mqttClientId := fmt.Sprintf("clientId_%s", deviceMac)

			mac, err := cli.ParesMacString(deviceMac)
			if err != nil {
				cli.Log.Fatalf("%v", err)
				os.Exit(1)
			}

			ha, err := homeAssistant.NewHomeAssistanceMqttClient(
				cli.Log, localMac, mac, username, password, host, port, token,
				mqttServerName, mqttClientId, mqttServerPort, mqttServerTls, mqttServerTlsValidaton,
				mqttBaseTopic, mqttDeviceName, mqttUserName, mqttPassword, mqttTelePeriod,
				byte(devicePort),
			)
			if err != nil {
				cli.Log.Fatalf("%v", err)
				os.Exit(2)
			}

			err = ha.Start()
			if err != nil {
				cli.Log.Fatalf("%v", err)
				os.Exit(3)
			}
		},
	}
	rootCmd.AddCommand(haCmd)

	haCmd.Flags().StringVarP(&mqttServerName, "mqtt", "H", "192.168.0.31", "MQTT server name or IP") // TODO change default to something public
	haCmd.Flags().StringVarP(&mqttUserName, "mqttUserName", "u", "", "MQTT server username")         // TODO change default to something public
	haCmd.Flags().StringVarP(&mqttPassword, "mqttPassword", "p", "", "MQTT server password")         // TODO change default to something public
	haCmd.Flags().IntVarP(&mqttServerPort, "mqttPort", "P", 1883, "MQTT server port")
	haCmd.Flags().BoolVarP(&mqttServerTls, "mqttTls", "s", false, "use TLS to connect MQTT server")
	haCmd.Flags().BoolVarP(&mqttServerTlsValidaton, "mqttStrictTlsValidation", "i", true, "if false, skip server certificate validation")
	haCmd.Flags().StringVarP(&mqttBaseTopic, "mqttBaseTopic", "b", "halsecur", "MQTT topic")
	haCmd.Flags().StringVarP(&mqttDeviceName, "name", "n", "garage", "Name of the local device in MQTT messages")
	haCmd.Flags().DurationVarP(&mqttTelePeriod, "mqttTelePeriod", "e", 1*time.Second, "Frequency of device state publish")
	haCmd.Flags().IntVar(&devicePort, devicePortName, 0, "Port number of the door")
}
