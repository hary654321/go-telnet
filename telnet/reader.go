package telnet

import (
	"log"
	"strings"
)

type Reader interface {
	Read([]byte) (int, error)
}

// readLine 从Reader中读取一行数据，直到遇到 '\n'
func ReadLine(r Reader) (string, error) {

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
			return "close", nil
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
