package cli

import "flag"

// ServerFlags ..
type ServerFlags struct {
	ServerName string
	ServerEnv  string
}

var (
	serverFlags = &ServerFlags{}
)

func Init() {

	flag.StringVar(&serverFlags.ServerName, "server_name", "", "Name of the server.")
	flag.StringVar(&serverFlags.ServerEnv, "server_env", "", "env of the server. release,staging,dev,private")

	flag.Parse()
}

// GetServerName ..
func GetServerName() string {
	return serverFlags.ServerName
}

// GetServerEnv ..
func GetServerEnv() string {
	return serverFlags.ServerEnv
}
