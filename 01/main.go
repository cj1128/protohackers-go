package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/big"
	"net"

	"cjting.me.protohackers/util"
	"github.com/k0kubun/pp/v3"
)

// NOTE: Any JSON number is a valid number, including floating-point values.
type Request struct {
	Method string   `json:"method"`
	Number *float64 `json:"number"`
}

type Response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

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

		util.Infoln("Client connected:", conn.RemoteAddr().String())

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		bytes := scanner.Bytes()
		util.Debugln("JSON payload:", string(bytes))

		req := parseReq(bytes)

		if req == nil {
			util.Errorln("Invalid json encountered")
			conn.Write([]byte("invalid json request\n"))
			break
		}

		res := Response{
			Method: "isPrime",
			Prime:  isPrime(*req.Number),
		}

		util.Infoln("Request received:", pp.Sprint(req))

		b, _ := json.Marshal(res)
		conn.Write(b)
		conn.Write([]byte("\n"))
	}

	if scanner.Err() != nil {
		fmt.Println("Error reading:", scanner.Err())
	}
}

// return nil if request is malformed
func parseReq(input []byte) *Request {
	req := &Request{}

	if err := json.Unmarshal(input, req); err != nil {
		return nil
	}

	if req.Method != "isPrime" || req.Number == nil {
		return nil
	}

	return req
}

// NOTE: float64 can store large integers that int64 can not fit
func isPrime(f float64) bool {
	bf := big.NewFloat(f)

	if !bf.IsInt() {
		return false
	}

	n, _ := bf.Int(nil)

	return n.ProbablyPrime(20)
}
