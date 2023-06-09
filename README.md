# Signum: A Standalone Application for Simulating the Minecraft Server Status and Login Protocol

Signum is a standalone application written in Golang that simulates the status and login protocol of a Minecraft server. It allows you to fake the presence of a Minecraft server, and reject client connections with a custom message.

## Minecraft Versions

- Supports 1.7.x-1.19.x out of the box
- Legacy-Support (1.4.x-1.6.x) can be enabled using the `SIGNUM_LEGACY` and `SIGNUM_STATUS_LEGACY` environment variables.

## Requirements

- Golang version 1.20 or later

## Installation

1. Clone the repository: `git clone https://github.com/Etrayed/signum`
2. Build the application: `go build`

## Usage

To build the Signum server, run the following command:

```
go build
```

The server will be built and ready to run. You must specify `SIGNUM_STATUS` and `SIGNUM_KICK_MESSAGE` as environment variables.

### Configuring the server

The Signum server can be configured using environment variables. The following variables are available:

- `SIGNUM_STATUS`: The status json sent to clients. (See [here](https://wiki.vg/Server_List_Ping#Status_Response))
- `SIGNUM_KICK_MESSAGE`: The message json sent to clients when they are kicked from the server. (See [here](https://wiki.vg/Chat) &rarr; If you just want to send a simple message, use something like this: `{"text":"Your Message here."}`)
- `SIGNUM_HOST`: The host address the server listens on. Default is "0.0.0.0".
- `SIGNUM_PORT`: The port the server listens on. Default is "25565".
- `SIGNUM_LEGACY`: A boolean indicating whether legacy pings should be supported.
- `SIGNUM_STATUS_LEGACY`: The legacy status response as JSON. Use the following format: `{"version": {"protocol": 127,"name": "Unknown"}, "motd": "A Minecraft Server", "currentPlayers": 0, "maxPlayers": 20}` (Note that the legacy format does not support special characters.)

To change the status or kick message, modify the corresponding environment variable and restart the server.
