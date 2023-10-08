package processor

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"signum/config"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
)

func Process(con net.Conn) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Failed to handle "+con.RemoteAddr().String()+":", r)
		}
	}()

	return expect(con, HANDSHAKE_HANDSHAKE, handshake)
}

var utf16beEncoder *encoding.Encoder

func doLegacyPing(con net.Conn) error {
	buf := &bytes.Buffer{}
	_, err := buf.Write([]byte{0xFF})

	if err != nil {
		return err
	}

	legacyStatus, err := config.GetLegacyStatus()

	if err != nil {
		return err
	}

	if utf16beEncoder == nil {
		utf16beEncoder = unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewEncoder()
	}

	responseStr := "\u00a71\u0000" // prefix

	responseStr += strconv.FormatInt(int64(legacyStatus.Version.Protocol), 10) + "\u0000" // protocol version
	responseStr += legacyStatus.Version.Name + "\u0000"                                   // protocol name
	responseStr += legacyStatus.Motd + "\u0000"                                           // motd
	responseStr += strconv.FormatInt(int64(legacyStatus.CurrentPlayers), 10) + "\u0000"   // current players
	responseStr += strconv.FormatInt(int64(legacyStatus.MaxPlayers), 10)                  // max players

	strBuf := &bytes.Buffer{}

	if _, err := utf16beEncoder.Writer(strBuf).Write([]byte(responseStr)); err != nil {
		return err
	}

	err = binary.Write(buf, binary.BigEndian, uint16(strBuf.Len()/2)) // write string length

	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	_, err = strBuf.WriteTo(buf)

	if err != nil {
		return err
	}

	_, err = io.Copy(con, buf)

	return err
}

func expect(con io.Reader, id int, handler func(io.Reader) error) error {
	packetLen, err := readVarInt(con) // packet length

	if err != nil {
		return err
	}

	actualId, err := readVarInt(con) // packet id

	if err != nil {
		return err
	}

	if id == HANDSHAKE_HANDSHAKE && config.IsLegacyPingEnabled() && packetLen == 0xFE && actualId == 0x01 { // check for legacy ping
		return doLegacyPing(con.(net.Conn))
	}

	if actualId == id {
		return handler(con)
	} else {
		return errors.New("Expected 0x" + strconv.FormatInt(int64(id), 16) + " but got 0x" + strconv.FormatInt(int64(actualId), 16))
	}
}

func handshake(con io.Reader) error {
	_, err := readVarInt(con) // protocol version, irrelevant

	if err != nil {
		return err
	}

	err = skipLengthed(con) // hostname

	if err != nil {
		return err
	}

	err = skip(con, 2) // port

	if err != nil {
		return err
	}

	nextState, err := readVarInt(con)

	if err != nil {
		return err
	}

	if nextState == HANDSHAKE_STATE_STATUS {
		return expect(con, STATUS_REQUEST, statusRequest)
	} else if nextState == HANDSHAKE_STATE_LOGIN {
		return expect(con, LOGIN_START, loginStart)
	}

	return errors.New("Unknown state provided: " + strconv.FormatInt(int64(nextState), 10))
}

func loginStart(con io.Reader) error {
	playerName, err := readString(con)

	if err != nil { // in this case, we actually dont care about the error since a player name is not necessary
		playerName = "???"
	}

	err = writePacket(con.(io.Writer), LOGIN_KICK, func(w io.Writer) error {
		return writeString(w, strings.ReplaceAll(config.GetKickMessage(), "%player%", playerName))
	})

	if err == nil {
		time.Sleep(500 * time.Millisecond) // wait to make sure the packet reaches the client before the connection is being closed
	}

	return err
}

func statusRequest(con io.Reader) error {
	err := writePacket(con.(io.Writer), STATUS_RESPONSE, func(w io.Writer) error {
		return writeString(w, config.GetStatus())
	})

	if err != nil {
		return err
	}

	return expect(con, STATUS_PING, func(r io.Reader) error {
		return writePacket(r.(io.Writer), STATUS_PONG, func(w io.Writer) error {
			buffer := make([]byte, 8)
			_, err := r.Read(buffer) // read timestamp sent by client

			if err != nil {
				return err
			}

			_, err = w.Write(buffer) // and just send it back

			return err
		})
	})
}
