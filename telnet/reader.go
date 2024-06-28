package telnet

import (
	"fmt"
	"log"
	"net"
	"strings"
	"telnet/json"
)

type Reader interface {
	Read([]byte) (int, error)
}

// readLine 从Reader中读取一行数据，直到遇到 '\n'
func ReadLine(c net.Conn, r Reader) (string, error) {

	res := ""
	for {
		var buffer [1]byte
		p := buffer[:]
		n, err := r.Read(p)
		if err != nil {
			return "", err
		}

		tmp := string(p[:n])

		res += tmp

		log.Println("tmp:", buffer)

		if buffer[0] == 3 {
			json.GlobalLog.HoneyLog(c, "close", nil)
			return "close", fmt.Errorf("exit")
		}

		if tmp == "\n" {
			log.Println("n")
			break
		}
		if tmp == "\r" {
			log.Println("r")
			break
		}
		if tmp == "\r\n" {
			log.Println("rn")
			break
		}
	}
	res = strings.TrimRight(res, "\r\n")

	return res, nil
}
