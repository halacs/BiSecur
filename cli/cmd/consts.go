package cmd

import "time"

const (
	ArgNameToken              = "token"
	ArgNameUsername           = "username"
	ArgNamePassword           = "password"
	ArgNameHost               = "host"
	ArgNamePort               = "port"
	ArgNameDeviceMac          = "mac"
	ArgNameDebug              = "debug"
	ArgNameAutoLogin          = "autologin"
	ArgNameLastLoginTimeStamp = "lastLogin"
	devicePortName            = "devicePort"
	UsersCmdUse               = "users"
	StatusCmdUse              = "status"
	SetStateCmdUse            = "set-state"
	GroupsCmdName             = "groups"
	LoginCmdName              = "login"
	LogoutCmdName             = "logout"
	PingCmdName               = "ping"
	TokenExpirationTime       = time.Minute * 5
	retryCount                = 3
)
