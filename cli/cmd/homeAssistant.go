package cmd

import (
	"bisecur/cli"
	"bisecur/cli/homeAssistant"
	"flag"
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

			mqttServerName = viper.GetString(ArgMqttServerName)
			mqttServerPort = viper.GetInt(ArgMqttPortName)
			mqttServerTls = viper.GetBool(ArgMqttTlsName)
			mqttServerTlsValidaton = viper.GetBool(ArgMqttStrictTlsValidationName)
			mqttBaseTopic = viper.GetString(ArgMqttBaseTopicName)
			mqttDeviceName = viper.GetString(ArgMqttDeviceNameName)
			mqttUserName = viper.GetString(ArgMqttUserNameName)
			mqttPassword = viper.GetString(ArgMqttPasswordName)
			mqttTelePeriod = viper.GetDuration(ArgMqttTelePeriodName)
			devicePort = viper.GetInt(ArgDevicePortName)

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

	haCmd.Flags().StringVarP(&mqttServerName, ArgMqttServerName, "H", "test.mosquitto.org", "MQTT server name or IP")
	haCmd.Flags().StringVarP(&mqttUserName, ArgMqttUserNameName, "u", "", "MQTT server username")
	haCmd.Flags().StringVarP(&mqttPassword, ArgMqttPasswordName, "p", "", "MQTT server password")
	haCmd.Flags().IntVarP(&mqttServerPort, ArgMqttPortName, "P", 1883, "MQTT server port")
	haCmd.Flags().BoolVarP(&mqttServerTls, ArgMqttTlsName, "s", false, "use TLS to connect MQTT server")
	haCmd.Flags().BoolVarP(&mqttServerTlsValidaton, ArgMqttStrictTlsValidationName, "i", true, "if false, skip server certificate validation")
	haCmd.Flags().StringVarP(&mqttBaseTopic, ArgMqttBaseTopicName, "b", "halsecur", "MQTT topic")
	haCmd.Flags().StringVarP(&mqttDeviceName, ArgMqttDeviceNameName, "n", "garage", "Name of the local device in MQTT messages")
	haCmd.Flags().DurationVarP(&mqttTelePeriod, ArgMqttTelePeriodName, "e", 1*time.Second, "Frequency of device state publish")
	haCmd.Flags().IntVar(&devicePort, ArgDevicePortName, 0, "Port number of the door")
	flag.Parse()
	err := viper.BindPFlags(haCmd.Flags())
	if err != nil {
		cli.Log.Fatalf("failed to bind flags: %v", err)
		os.Exit(1)
	}
}
