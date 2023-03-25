package constants

const APPLICATION_NAME string = "QR-AUTH"

const DEFAULT_PORT string = "1515"

var ENV_NAMES = EnvNames{
	ENV:  "ENV",
	PORT: "PORT",
}

var EVENTS = Events{
	RegisterConnection: "register-connection",
	RegisterUser:       "register-user",
}
