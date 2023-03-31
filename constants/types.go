package constants

type EnvNames struct {
	ENV_TYPE string
	PORT     string
}

type Events struct {
	AuthenticateTarget string
	ClientDisconnect   string
	InvalidTarget      string
	Ping               string
	PingResponse       string
	RegisterConnection string
	RegisterUser       string
	ServerDisconnect   string
	SignOut            string
	Unauthorized       string
}
