package bisecur

import "time"

const (
	/*
		TODO
		My gateway drops my out after the first relogin attempt and became crazy for a while.
		Don't know yet why. Would be nice to find out the right value for this const if needed at all.
	*/
	TokenExpirationTime = 6 * time.Hour
)
