package cmd

import (
	"context"
	"flag"
	"fmt"
	"halsecur/cli"
	"halsecur/cli/homeAssistant"
	"halsecur/cli/utils"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configFileName = "config.yaml"
)

func getIntSlice(key string) []int {
	raw := viper.Get(key)

	if slice, err := cast.ToIntSliceE(raw); err == nil && len(slice) > 0 {
		return slice
	}

	str := viper.GetString(key)
	if str == "" {
		return []int{}
	}
	var result []int
	for _, s := range strings.Split(str, ",") {
		if n, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
			result = append(result, n)
		}
	}
	return result
}

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
		mqttTelePeriodFast     time.Duration
		devicePorts            []int
		doorStatusSupported    bool
	)

	haCmd := &cobra.Command{
		Use:     HomeAssistantCmdName,
		Short:   "Start MQTT client compatible with Home assistant auto discovery",
		Long:    ``,
		PreRunE: preRunFuncs,
		Run: func(cmd *cobra.Command, args []string) {
			var ha *homeAssistant.HomeAssistanceMqttClient

			deviceMac = viper.GetString(ArgNameDeviceMac)
			host = viper.GetString(ArgNameHost)
			port = viper.GetInt(ArgNamePort)
			token = viper.GetUint32(ArgNameToken)

			lastLoginTime = viper.GetInt64(ArgNameLastLoginTimeStamp)
			lastLoginTimeDt := time.UnixMicro(lastLoginTime)

			username = viper.GetString(ArgNameUsername)
			password = viper.GetString(ArgNamePassword)

			mqttServerName = viper.GetString(ArgMqttServerName)
			mqttServerPort = viper.GetInt(ArgMqttPortName)
			mqttServerTls = viper.GetBool(ArgMqttTlsName)
			mqttServerTlsValidaton = viper.GetBool(ArgMqttStrictTlsValidationName)
			mqttBaseTopic = viper.GetString(ArgMqttBaseTopicName)
			mqttDeviceName = viper.GetString(ArgMqttDeviceNameName)
			mqttUserName = viper.GetString(ArgMqttUserNameName)
			mqttPassword = viper.GetString(ArgMqttPasswordName)
			mqttTelePeriod = viper.GetDuration(ArgMqttTelePeriodName)
			mqttTelePeriodFast = viper.GetDuration(ArgMqttTelePeriodFastName)
			devicePorts = getIntSlice(ArgDevicePortsName)
			doorStatusSupported = viper.GetBool(ArgDoorStatusSupported)

			autoTokenRefresh := viper.GetBool(ArgNameAutoLogin)
			if !autoTokenRefresh {
				cli.Log.Warningf("Auto re-login is disabled by '%s' option.", ArgNameAutoLogin)
			}

			if len(devicePorts) == 0 {
				cli.Log.Warnf("%s parameter empty. This might be wrong.", ArgDevicePortsName)
			}

			// Store token in persistent
			defer func() {
				viper.Set(ArgNameToken, token)
				viper.Set(ArgNameLastLoginTimeStamp, lastLoginTime)
				_, err := os.Stat("config.yaml")
				if os.IsNotExist(err) {
					err = viper.WriteConfigAs("config.yaml")
				} else {
					err = viper.WriteConfig()
				}
				if err != nil {
					cli.Log.Errorf("Failed to save new configuration. %v", err)
				}
			}()

			mqttClientId := fmt.Sprintf("clientId_%s", deviceMac)

			mac, err := utils.ParesMacString(deviceMac)
			if err != nil {
				cli.Log.Fatalf("%v", err)
				os.Exit(1)
			}

			// Store token in persistent config
			defer func() {
				newToken := ha.GetToken()
				newTs := ha.GetLastLoginTime()
				cli.Log.Debugf("Saving values of token (%v) and lastLoginTime (%v) to config file...", newToken, newTs)
				viper.Set(ArgNameToken, newToken)
				viper.Set(ArgNameLastLoginTimeStamp, newTs.UnixMicro())

				_, err := os.Stat(configFileName)
				if os.IsNotExist(err) {
					err = viper.WriteConfigAs(configFileName)
				} else {
					err = viper.WriteConfig()
				}
				if err != nil {
					cli.Log.Errorf("failed to save new configuration. %v", err)
				}
			}()

			// We need to terminate based on system signals. This is essential to save token when exiting.
			ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)

			ha, err = homeAssistant.NewHomeAssistanceMqttClient(
				ctx, cli.Log, localMac, mac, username, password, host, port, token, lastLoginTimeDt,
				mqttServerName, mqttClientId, mqttServerPort, mqttServerTls, mqttServerTlsValidaton,
				mqttBaseTopic, mqttDeviceName, mqttUserName, mqttPassword, mqttTelePeriod, mqttTelePeriodFast,
				utils.IntArrayToByteArray(devicePorts), doorStatusSupported, autoTokenRefresh,
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

			<-ctx.Done()
			cli.Log.Infof("Exiting...")
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
	haCmd.Flags().DurationVarP(&mqttTelePeriod, ArgMqttTelePeriodName, "e", 15*time.Second, "Frequency of device state publish")
	haCmd.Flags().DurationVarP(&mqttTelePeriodFast, ArgMqttTelePeriodFastName, "f", 5*time.Second, "Frequency of device state publish when door might be moving")
	haCmd.Flags().IntSliceVar(&devicePorts, ArgDevicePortsName, []int{}, "Port numbers of the doors")
	haCmd.Flags().BoolVar(&doorStatusSupported, ArgDoorStatusSupported, true, "Whether the controlled door supports door status (opening state) or not")
	flag.Parse()
	err := viper.BindPFlags(haCmd.Flags())
	if err != nil {
		cli.Log.Fatalf("failed to bind flags: %v", err)
		os.Exit(1)
	}

	// ENV support
	viper.SetEnvPrefix("HALSECUR")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

}
