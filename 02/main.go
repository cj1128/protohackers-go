package main

import (
	"encoding/binary"
	"io"
	"math/big"
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
	prices := make(map[int32]int32)
	// var log = log.New(os.Stdout, conn.RemoteAddr().String()+" ", 0)

	buf := make([]byte, 9)

	for {
		_, err := io.ReadFull(conn, buf)

		if err != nil {
			if err != io.EOF {
				util.Errorln("Error reading: ", err.Error())
			}
			break
		}

		switch buf[0] {
		case 'I':
			{
				timestamp := int32(binary.BigEndian.Uint32(buf[1:5]))
				price := int32(binary.BigEndian.Uint32(buf[5:]))

				// undefined behaviour
				if _, ok := prices[timestamp]; ok {
					break
				}

				// log.Printf("Insert req: timestamp=%d, price=%d\n", timestamp, price)
				prices[timestamp] = price
			}
		case 'Q':
			{
				min := int32(binary.BigEndian.Uint32(buf[1:5]))
				max := int32(binary.BigEndian.Uint32(buf[5:]))

				mean := query(prices, min, max)
				// log.Printf("Query req: min=%d, max=%d, mean=%d\n", min, max, mean)

				buf := make([]byte, 4)
				binary.BigEndian.PutUint32(buf, uint32(mean))
				conn.Write(buf)
			}

		default:
			util.Errorln("Invaild msg:", buf)
			break
		}
	}
}

func query(prices map[int32]int32, min, max int32) int32 {
	if len(prices) == 0 || min > max {
		return 0
	}

	total := big.NewInt(0)
	count := big.NewInt(0)
	one := big.NewInt(1)

	// we don't care the performance here
	for ts, price := range prices {
		if ts >= min && ts <= max {
			total.Add(total, big.NewInt(int64(price)))
			count.Add(count, one)
		}
	}

	result := big.NewInt(0)

	if count.Int64() != 0 {
		result = result.Div(total, count)
	}

	return int32(result.Int64())
}
