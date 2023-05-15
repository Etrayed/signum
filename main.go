package main

import (
	"errors"
	"io"
	"log"
	"net"
	"os"
	"signum/config"
	"signum/processor"
	"time"
)

func main() {
	server, err := net.Listen("tcp", config.GetAddress())

	if err != nil {
		log.Fatal(err)
	}

	defer server.Close()

	for {
		con, err := server.Accept()

		if err != nil {
			log.Fatal(err)
		}

		go processConnection(con)
	}
}

func processConnection(connection net.Conn) {
	connection.SetDeadline(time.Now().Add(3 * time.Second))

	defer connection.Close()

	err := processor.Process(connection)

	if errors.Is(err, io.EOF) || os.IsTimeout(err) {
		return
	}

	if err != nil {
		println("Failed to process "+connection.RemoteAddr().String()+":", err.Error())
	}
}
