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

/**

`{
		"version": {
			"name": "1.19.4",
			"protocol": 762
		},
		"players": {
			"max": 100,
			"online": 5,
			"sample": [
				{
					"name": "thinkofdeath",
					"id": "4566e69f-c907-48ee-8d71-d7ba5aa00d20"
				}
			]
		},
		"description": {
			"text": "Hello world"
		}
	}`
	"{\"text\":\"§cNot today, %player%!\"}"
*/

func GetKickMessage() string {
	return os.Getenv("SIGNUM_KICK_MESSAGE")
}
