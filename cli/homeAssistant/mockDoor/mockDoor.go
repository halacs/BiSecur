package mockDoor

import (
	"bisecur/cli/utils"
	"time"
)

var (
	openPercent   = 0
	lastDirection = utils.CLOSING
	direction     = utils.STOPPED
)

func StartTicker() {
	ticker := time.NewTicker(500 * time.Millisecond)
	go func() {
		for range ticker.C {
			switch direction {
			case utils.OPENING:
				openPercent = openPercent + 1
				if openPercent == 100 {
					SetPosition() // make door state stopped
				}
			case utils.CLOSING:
				openPercent = openPercent - 1
				if openPercent == 0 {
					SetPosition() // make door state open
				}
			}
		}
	}()

}

func GetPosition() int {
	return openPercent
}

func GetDirection() string {
	return direction
}

func GetLastDirection() string {
	return lastDirection
}

/*
Start moving the door in the opposite direction or stop it if it was moving.
*/
func SetPosition() {
	switch direction {
	case utils.STOPPED:
		switch lastDirection {
		case utils.OPENING:
			lastDirection = direction
			direction = utils.CLOSING
		case utils.CLOSING:
			lastDirection = direction
			direction = utils.OPENING
		}
	case utils.OPENING:
		lastDirection = direction
		direction = utils.STOPPED
	case utils.CLOSING:
		lastDirection = direction
		direction = utils.STOPPED
	}
}

func SetStateMockMultiCall(count int) error {
	for i := 0; i < count; i++ {
		SetPosition()
	}

	return nil
}
