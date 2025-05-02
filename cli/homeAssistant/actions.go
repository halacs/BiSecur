package homeAssistant

import (
	"bisecur/cli"
	"bisecur/cli/bisecur"
	"bisecur/cli/homeAssistant/mockDoor"
	"bisecur/cli/utils"
	"fmt"
	"time"
)

func (ha *HomeAssistanceMqttClient) autoLoginBisecur() error {
	if ha.tokenCreated.Add(bisecur.TokenExpirationTime).Before(time.Now()) {
		cli.Log.Infof("Token expired. Logging in...")

		var err error

		err = bisecur.Logout(ha.localMac, ha.deviceMac, ha.host, ha.port, ha.token)
		if err != nil {
			ha.log.Warnf("failed to logout. %v", err)
		}

		ha.token, err = bisecur.Login(ha.localMac, ha.deviceMac, ha.host, ha.port, ha.deviceUsername, ha.devicePassword)
		if err != nil {
			return fmt.Errorf("failed to auto login. %v", err)
		}

		ha.tokenCreated = time.Now() // note when token was received
	}

	return nil
}

func (ha *HomeAssistanceMqttClient) setStateBisecurMultiCall(count int) error {
	for i := 0; i < count; i++ {
		ha.log.Debugf("Setting door state %d/%d", i+1, count)

		err := ha.autoLoginBisecur()
		if err != nil {
			return fmt.Errorf("failed to auto login. %v", err)
		}

		err = bisecur.SetState(ha.localMac, ha.deviceMac, ha.host, ha.port, ha.devicePort, ha.token)
		if err != nil {
			return fmt.Errorf("failed to get door status. %v", err)
		}

		if i < count-1 {
			const delayDuration = 1000 * time.Millisecond
			ha.log.Debugf("Waiting for %d ms before the next call...", delayDuration)
			time.Sleep(delayDuration) // wait for 1 second before the next door call to avoid overloading the Hormann bisecur gateway
		}
	}
	return nil
}

func (ha *HomeAssistanceMqttClient) openDoor() error {
	direction := mockDoor.GetDirection()
	lastDirection := mockDoor.GetLastDirection()
	position := mockDoor.GetPosition()

	switch direction {
	case utils.CLOSING:
		mockDoor.SetPosition()
		mockDoor.SetPosition()
	case utils.STOPPED:
		if position < 100 { // check if door is not already open
			mockDoor.SetPosition()
			if lastDirection == utils.OPENING { // check if door needs to be reversed
				mockDoor.SetPosition() // stop moving
				mockDoor.SetPosition() // reverse
			}
		}
	case utils.OPENING:
		// nothing to do, it's already opening
	}
	return nil
}

func (ha *HomeAssistanceMqttClient) closeDoor() error {
	direction := mockDoor.GetDirection()
	lastDirection := mockDoor.GetLastDirection()
	position := mockDoor.GetPosition()

	switch direction {
	case utils.OPENING:
		mockDoor.SetPosition()
		mockDoor.SetPosition()
	case utils.STOPPED:
		if position > 0 { // check if door is not already closed
			mockDoor.SetPosition()
			if lastDirection == utils.CLOSING { // check if door needs to be reversed
				mockDoor.SetPosition() // stop moving
				mockDoor.SetPosition() // reverse
			}
		}
	case utils.CLOSING:
		// nothing to do, it's already closing
	}

	return nil
}

func (ha *HomeAssistanceMqttClient) stopDoor() error {
	direction := mockDoor.GetDirection()

	if direction != utils.STOPPED {
		mockDoor.SetPosition()
	}

	return nil
}

func (ha *HomeAssistanceMqttClient) getDoorStatus() (direction string, position int, err error) {
	/*
		status="{\"StateInPercent\":0,\"DesiredStateInPercent\":0,\"Error\":false,\"AutoClose\":false,\"DriveTime\":0,\"Gk\":257,\"Hcp\":{\"PositionOpen\":false,\"PositionClose\":true,\"OptionRelais\":false,\"LightBarrier\":false,\"Error\":false,\"DrivingToClose\":false,\"Driving\":false,\"HalfOpened\":false,\"ForecastLeadTime\":false,\"Learned\":true,\"NotReferenced\":false},\"Exst\":\"AAAAAAAAAAA=\",\"Time\":\"2025-04-28T17:41:02.979836814+02:00\"}"
	*/
	/*
		position = mockDoor.GetPosition()
		direction = mockDoor.GetDirection()
		return direction, position, nil
	*/

	err = ha.autoLoginBisecur()
	if err != nil {
		return utils.UNKNOWN, 0, fmt.Errorf("failed to auto login. %v", err)
	}

	status, err := bisecur.GetStatus(ha.localMac, ha.deviceMac, ha.host, ha.port, ha.devicePort, ha.token)
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
