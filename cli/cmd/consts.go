package cmd

import "time"

const (
	ArgNameUserId             = "uid"
	ArgNameNewPassword        = "newpassword"
	ArgNameToken              = "token"
	ArgNameUsername           = "username"
	ArgNamePassword           = "password"
	ArgNameHost               = "host"
	ArgNamePort               = "port"
	ArgNameDeviceMac          = "mac"
	ArgNameDebug              = "debug"
	ArgNameJsonOutput         = "json"
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
