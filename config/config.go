package config

import "os"

func GetAddress() string {
	host := os.Getenv("SIGNUM_HOST")
	port := os.Getenv("SIGNUM_PORT")

	if len(port) == 0 {
		port = "25565"
	}

	return host + ":" + port
}

func GetStatus() string {
	return os.Getenv("SIGNUM_STATUS")
}

func GetKickMessage() string {
	return os.Getenv("SIGNUM_KICK_MESSAGE")
}
