package cmd

import (
	"fmt"
	"halsecur/cli"
	"halsecur/cli/bisecur"
	"halsecur/cli/utils"
	"halsecur/sdk"
	"os"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var portTypeNames = map[byte]string{
	sdk.PORT_TYPE_NONE:       "NONE",
	sdk.PORT_TYPE_IMPULS:     "IMPULS",
	sdk.PORT_TYPE_AUTO_CLOSE: "AUTO_CLOSE",
	sdk.PORT_TYPE_ON_OFF:     "ON_OFF",
	sdk.PORT_TYPE_UP:         "UP",
	sdk.PORT_TYPE_DOWN:       "DOWN",
	sdk.PORT_TYPE_HALF:       "HALF",
	sdk.PORT_TYPE_WALK:       "WALK",
	sdk.PORT_TYPE_LIGHT:      "LIGHT",
	sdk.PORT_TYPE_ON:         "ON",
	sdk.PORT_TYPE_OFF:        "OFF",
	sdk.PORT_TYPE_LOCK:       "LOCK",
	sdk.PORT_TYPE_UNLOCK:     "UNLOCK",
	sdk.PORT_TYPE_OPEN_DOOR:  "OPEN_DOOR",
	sdk.PORT_TYPE_LIFT:       "LIFT",
	sdk.PORT_TYPE_SINK:       "SINK",
}

var portsListCmd = &cobra.Command{
	Use:     PortsListCmdName,
	Short:   "List all ports and their types on the gateway",
	Long:    `List all configured ports on the gateway with their type (IMPULS, AUTO_CLOSE, etc).`,
	PreRunE: preRunFuncs,
	Run: func(cmd *cobra.Command, args []string) {
		deviceMac := viper.GetString(ArgNameDeviceMac)
		host := viper.GetString(ArgNameHost)
		port := viper.GetInt(ArgNamePort)
		token := viper.GetUint32(ArgNameToken)

		mac, err := utils.ParesMacString(deviceMac)
		if err != nil {
			cli.Log.Fatalf("%v", err)
			os.Exit(1)
		}

		portIds, err := bisecur.GetPorts(localMac, mac, host, port, token)
		if err != nil {
			cli.Log.Fatalf("%v", err)
			os.Exit(2)
		}

		if len(portIds) == 0 {
			cli.Log.Infof("No ports configured on this gateway")
			return
		}

		for _, portId := range portIds {
			portType, err := bisecur.GetType(localMac, mac, host, port, portId, token)
			if err != nil {
				cli.Log.Warnf("Port %d: failed to get type: %v", portId, err)
				continue
			}

			typeName, ok := portTypeNames[portType]
			if !ok {
				typeName = fmt.Sprintf("UNKNOWN(%d)", portType)
			}

			cli.Log.Infof("Port %d: %s (%d)", portId, typeName, portType)
		}
	},
}

func init() {
	rootCmd.AddCommand(portsListCmd)
}
