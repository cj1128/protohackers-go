package main

import (
	"io"
	"net"

	"cjting.me.protohackers/util"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:8888")

	if err != nil {
		util.Fatalln("Error listening: ", err.Error())
	}

	defer l.Close()

	util.Infoln("Server started on :8888")

	for {
		conn, err := l.Accept()

		if err != nil {
			util.Errorln("Failed to accept: ", err.Error())
			continue
		}

		util.Infoln("Client connected: ", conn.RemoteAddr().String())

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)

		if n > 0 {
			conn.Write(buffer[:n])
		}

		if err != nil {
			if err != io.EOF {
				util.Errorln("Error reading: ", err.Error())
			}
			break
		}
	}
}
