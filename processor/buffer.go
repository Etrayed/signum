package processor

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

const SEGMENT_BITS = 0x7F
const CONTINUE_BIT = 0x80

func readVarInt(con io.Reader) (int, error) {
	var val int32 = 0
	var pos int = 0
	var currentByte byte

	for {
		err := binary.Read(con, binary.BigEndian, &currentByte)

		if err != nil {
			return -1, err
		}

		val |= int32(currentByte&SEGMENT_BITS) << pos

		if (currentByte & CONTINUE_BIT) == 0 {
			break
		}

		pos += 7

		if pos >= 32 {
			return -1, errors.New("VarInt too big")
		}
	}

	return int(val), nil
}

func writeVarInt(w io.Writer, val int) error {
	for {
		if (val & SEGMENT_BITS) == 0 {
			_, err := w.Write([]byte{byte(val)})

			return err
		}

		_, err := w.Write([]byte{byte((val & SEGMENT_BITS) | CONTINUE_BIT)})

		if err != nil {
			return err
		}

		val = int(uint(val) >> 7)
	}
}

func writeString(w io.Writer, val string) error {
	data := []byte(val)

	err := writeVarInt(w, len(data))

	if err != nil {
		return err
	}

	_, err = w.Write(data)

	return err
}

func readString(r io.Reader) (string, error) {
	length, err := readVarInt(r)

	if err != nil {
		return "", err
	}

	data := make([]byte, length)

	_, err = r.Read(data)

	if err != nil {
		return "", err
	}

	return string(data), nil
}

func writePacket(w io.Writer, id int, writer func(w io.Writer) error) error {
	buf := &bytes.Buffer{}

	err := writeVarInt(buf, id) // id

	if err != nil {
		return err
	}

	err = writer(buf) // data

	if err != nil {
		return err
	}

	err = writeVarInt(w, buf.Len()) // write length

	if err != nil {
		return err
	}

	_, err = io.Copy(w, buf) // write contents

	return err
}

var SKIP_BUF = make([]byte, 1)

func skipLengthed(con io.Reader) error {
	length, err := readVarInt(con)

	if err != nil {
		return err
	}

	return skip(con, length)
}

func skip(con io.Reader, amount int) error {
	for i := 0; i < amount; i++ {
		_, err := con.Read(SKIP_BUF)

		if err != nil {
			return err
		}
	}

	return nil
}
