package homeAssistant

import (
	"bisecur/cli/homeAssistant/mockDoor"
	"bisecur/cli/utils"
)

func openDoor() error {
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

func closeDoor() error {
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

func stopDoor() error {
	direction := mockDoor.GetDirection()

	if direction != utils.STOPPED {
		mockDoor.SetPosition()
	}

	return nil
}
