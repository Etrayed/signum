package config

import (
	"encoding/json"
	"os"
	"strconv"
)

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

func IsLegacyPingEnabled() bool {
	b, err := strconv.ParseBool(os.Getenv("SIGNUM_LEGACY"))

	return err == nil && b
}

type Version struct {
	Protocol int
	Name     string
}

type LegacyStatus struct {
	Version        Version
	Motd           string
	CurrentPlayers int
	MaxPlayers     int
}

var _legacyStatus LegacyStatus

func GetLegacyStatus() (LegacyStatus, error) {
	if (_legacyStatus == LegacyStatus{}) {
		err := json.Unmarshal([]byte(os.Getenv("SIGNUM_STATUS_LEGACY")), &_legacyStatus)

		if err != nil {
			return LegacyStatus{}, err
		}
	}

	return _legacyStatus, nil
}

func GetKickMessage() string {
	return os.Getenv("SIGNUM_KICK_MESSAGE")
}
