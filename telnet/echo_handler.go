package telnet

import (
	"log"
	"net"
	"telnet/cmd"
	"telnet/json"
)

// EchoHandler is a simple TELNET server which "echos" back to the client any (non-command)
// data back to the TELNET client, it received from the TELNET client.
var EchoHandler Handler = internalEchoHandler{}

type internalEchoHandler struct{}

func (handler internalEchoHandler) ServeTELNET(c net.Conn, w Writer, r Reader) error {

	cmdstr, err := ReadLine(c, r)

	if err != nil {
		log.Println("读取命令失败:", err)
		return err
	}

	if cmdstr != "" {

		output := cmd.Cmd(cmdstr)
		extend := make(map[string]any)
		extend["cmd"] = cmdstr
		extend["output"] = output

		json.GlobalLog.HoneyLog(c, "op", extend)

		w.Write([]byte(output + "\r\n"))
	}

	return nil
}
