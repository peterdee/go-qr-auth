package constants

const APPLICATION_NAME string = "QR-AUTH"

const CONNECTION_TIMEOUT int64 = 4 * 30 * 1000 // 2 minutes

const DEFAULT_PORT string = "1515"

var ENV_NAMES = EnvNames{
	ENV_TYPE: "ENV_TYPE",
	PORT:     "PORT",
}

var EVENTS = Events{
	AuthenticateTarget: "authenticate-target",
	ClientDisconnect:   "client-disconnect",
	InvalidTarget:      "invalid-target",
	Ping:               "ping",
	PingResponse:       "pong",
	RegisterConnection: "register-connection",
	RegisterUser:       "register-user",
	ServerDisconnect:   "server-disconnect",
	Unauthorized:       "unauthorized",
}
