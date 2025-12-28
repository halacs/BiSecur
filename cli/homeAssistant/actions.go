package homeAssistant

import (
	"bisecur/cli"
	"bisecur/cli/bisecur"
	"bisecur/cli/utils"
	"bisecur/sdk/payload"
	"fmt"
	"time"
)

func (ha *HomeAssistanceMqttClient) autoLoginBisecur() error {
	if ha.tokenCreated.Add(bisecur.TokenExpirationTime).Before(time.Now()) {
		ha.log.Info("Token expired.")
		return ha.forceReLogin()
	}

	return nil
}

func (ha *HomeAssistanceMqttClient) LogoutBisecur() error {
	cli.Log.Infof("Logging out from the gateway...")

	err := bisecur.Logout(ha.localMac, ha.deviceMac, ha.host, ha.port, ha.token)
	if err != nil {
		return err
	}

	// Note that the token has been invalidated
	ha.token = 0
	ha.tokenCreated = time.UnixMicro(0)

	return nil
}

func (ha *HomeAssistanceMqttClient) forceReLogin() error {
	cli.Log.Infof("Logging in to the gateway...")

	var err error

	err = bisecur.Logout(ha.localMac, ha.deviceMac, ha.host, ha.port, ha.token)
	if err != nil {
		ha.log.Errorf("failed to logout. %v", err)
	}
	// clear token and the timestamp of the token after the successful logout
	ha.token = 0
	ha.tokenCreated = time.UnixMicro(0)

	// Not sure, this is really needed, but since I know Hormann BS gateway can become crazy if it gets overloaded...
	time.Sleep(2 * time.Second)

	ha.token, err = bisecur.Login(ha.localMac, ha.deviceMac, ha.host, ha.port, ha.deviceUsername, ha.devicePassword)
	if err != nil {
		return fmt.Errorf("login failed. %v", err)
	}

	ha.tokenCreated = time.Now() // note when token was received

	return nil
}

func (ha *HomeAssistanceMqttClient) setStateMultiCall(count int, devicePort byte) error {
	return ha.setStateBisecurMultiCall(count, devicePort)
	//return mockDoor.SetStateMockMultiCall(count)
}

func (ha *HomeAssistanceMqttClient) setStateBisecurMultiCall(count int, devicePort byte) error {
	const delayDuration = 1500 * time.Millisecond

	for i := 0; i < count; i++ {
		ha.log.Debugf("Setting door state %d/%d", i+1, count)

		err := ha.autoLoginBisecur()
		if err != nil {
			return fmt.Errorf("auto login failed. %v", err)
		}

		err = bisecur.SetState(ha.localMac, ha.deviceMac, ha.host, ha.port, devicePort, ha.token)
		if err != nil {
			return fmt.Errorf("failed to get door status. %v", err)
		}

		if i < count-1 {
			ha.log.Debugf("Waiting for %d ms before the next call...", delayDuration)
			time.Sleep(delayDuration) // wait for 1 second before the next door call to avoid overloading the Hormann bisecur gateway
		}
	}
	return nil
}

func (ha *HomeAssistanceMqttClient) impuls(devicePort byte) any {
	ha.log.Info("Sending impuls...")

	err := ha.setStateMultiCall(1, devicePort)
	if err != nil {
		return fmt.Errorf("failed to send impuls. %v", err)
	}

	return nil
}

func (ha *HomeAssistanceMqttClient) openDoor(devicePort byte) error {
	ha.log.Info("Opening door...")

	direction, position, err := ha.getDoorStatus(devicePort)
	if err != nil {
		return fmt.Errorf("failed to get door status. %v", err)
	}

	switch direction {
	case utils.CLOSING:
		err := ha.setStateMultiCall(2, devicePort)
		if err != nil {
			return fmt.Errorf("failed to set state. %v", err)
		}
	case utils.STOPPED, utils.OPEN, utils.CLOSED:
		if position < 100 { // check if door is not already fully open
			err := ha.setStateMultiCall(1, devicePort)
			if err != nil {
				return fmt.Errorf("failed to set state. %v", err)
			}

			newDirection, _, err := ha.getDoorStatus(devicePort)
			if err != nil {
				return fmt.Errorf("failed to get door status to confirm it is moving into the right direction. %v", err)
			}

			if newDirection != utils.OPENING && newDirection != utils.OPEN { // check if door needs to be reversed
				err := ha.setStateMultiCall(2, devicePort)
				if err != nil {
					return fmt.Errorf("failed to reverse moving direction. %v", err)
				}
			}
		}
	case utils.OPENING:
		ha.log.Infof("Door is already opening. Nothing to do.")
	default:
		ha.log.Infof("Unhandled direction in openDoor action: %s. Doing nothing.", direction)
	}

	ha.setRequestFastDootStatus()

	return nil
}

func (ha *HomeAssistanceMqttClient) closeDoor(devicePort byte) error {
	ha.log.Info("Closing door...")

	direction, position, err := ha.getDoorStatus(devicePort)
	if err != nil {
		return fmt.Errorf("failed to get door status. %v", err)
	}

	switch direction {
	case utils.OPENING:
		err := ha.setStateMultiCall(2, devicePort) // stop then reverse
		if err != nil {
			return fmt.Errorf("failed to set state. %v", err)
		}
	case utils.STOPPED, utils.OPEN:
		if position > 0 { // check if door is not already fully closed
			err := ha.setStateMultiCall(1, devicePort)
			if err != nil {
				return fmt.Errorf("failed to set state. %v", err)
			}

			newDirection, _, err := ha.getDoorStatus(devicePort)
			if err != nil {
				return fmt.Errorf("failed to get door status to confirm it is moving into the right direction. %v", err)
			}

			if newDirection != utils.CLOSING && newDirection != utils.CLOSED { // check if door needs to be reversed
				err := ha.setStateMultiCall(2, devicePort)
				if err != nil {
					return fmt.Errorf("failed to reverse moving direction. %v", err)
				}
			}
		}
	case utils.CLOSING:
		ha.log.Infof("Door is already closing. Nothing to do.")
	default:
		ha.log.Infof("Unhandled direction in closeDoor action: %s. Doing nothing.", direction)
	}

	ha.setRequestFastDootStatus()

	return nil
}

func (ha *HomeAssistanceMqttClient) stopDoor(devicePort byte) error {
	ha.log.Infof("Stopping door...")

	direction, _, err := ha.getDoorStatus(devicePort)
	if err != nil {
		return fmt.Errorf("failed to get door status. %v", err)
	}

	if direction == utils.OPENING || direction == utils.CLOSING { // anything which means moving door
		err := ha.setStateMultiCall(1, devicePort)
		if err != nil {
			return fmt.Errorf("failed to stop the door. %v", err)
		}
	} else {
		ha.log.Infof("Door is not moving. Nothing to stop. Current direction: %s", direction)
	}

	ha.setRequestFastDootStatus()

	return nil
}

func (ha *HomeAssistanceMqttClient) getDoorStatus(devicePort byte) (direction string, position int, err error) {
	/*
		status="{\"StateInPercent\":0,\"DesiredStateInPercent\":0,\"ErrorResponse\":false,\"AutoClose\":false,\"DriveTime\":0,
		\"Gk\":257,\"Hcp\":{\"PositionOpen\":false,\"PositionClose\":true,\"OptionRelais\":false,\"LightBarrier\":false,
		\"ErrorResponse\":false,\"DrivingToClose\":false,\"Driving\":false,\"HalfOpened\":false,\"ForecastLeadTime\":false,
		\"Learned\":true,\"NotReferenced\":false},\"Exst\":\"AAAAAAAAAAA=\",\"Time\":\"2025-04-28T17:41:02.979836814+02:00\"}"
	*/
	/*
		position = mockDoor.GetPosition()
		direction = mockDoor.GetDirection()
		return direction, position, nil
	*/
	err = ha.autoLoginBisecur()
	if err != nil {
		return utils.UNKNOWN, 0, fmt.Errorf("auto login failed. %v", err)
	}

	var status *payload.HmGetTransitionResponse
	err = utils.RetryAlways(utils.RetryCount, func() error {
		var err2 error
		status, err2 = bisecur.GetStatus(ha.localMac, ha.deviceMac, ha.host, ha.port, devicePort, ha.token)

		if err2.Error() == "PERMISSION_DENIED" { // TODO don't like string comparisons so should be refactored somehow while relogin also should be make more generic (think of other commands)
			// Does it make sense to force relogin after a PERMISSION_DENIED error?
			time.Sleep(2 * time.Second)
			err3 := ha.forceReLogin()
			if err3 != nil {
				return fmt.Errorf("error while re-login after a PERMISSION_DENIED error. %v. %v", err2, err3)
			}
		}

		return err2
	})

	if err != nil {
		return utils.UNKNOWN, 0, fmt.Errorf("failed to get door status. %v", err)
	}

	position = status.StateInPercent

	if status.Hcp.Driving {
		// Door is moving
		if status.Hcp.DrivingToClose {
			direction = utils.CLOSING
		} else {
			direction = utils.OPENING
		}
	} else {
		// Door is not moving
		if status.Hcp.PositionOpen {
			direction = utils.OPEN
		} else if status.Hcp.PositionClose {
			direction = utils.CLOSED
		} else if status.Hcp.HalfOpened {
			ha.log.Warnf("Door is half opened. This is not supported by Home Assistant so set it as OPEN.")
			direction = utils.OPEN
		} else {
			// not fully open or closed, neither half-open

			//direction = utils.STOPPED
			if position > 0 && position <= 100 {
				direction = utils.OPEN
			} else if position == 0 {
				direction = utils.CLOSED
			} else {
				ha.log.Errorf("Door position is %d but direction is unknown. This should not happen.", position)
				direction = utils.UNKNOWN
			}
		}
	}

	ha.log.Infof("Door position: %d, direction: %s", position, direction)

	// Some sanity checks...
	if position > 0 && direction != utils.OPEN && direction != utils.OPENING && direction != utils.CLOSING {
		ha.log.Errorf("Door position is %d but direction is %s. This should not happen.", position, direction)
	}

	return direction, position, nil
}
