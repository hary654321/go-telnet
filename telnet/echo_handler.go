package telnet

import (
	"log"
	"telnet/cmd"
)

// EchoHandler is a simple TELNET server which "echos" back to the client any (non-command)
// data back to the TELNET client, it received from the TELNET client.
var EchoHandler Handler = internalEchoHandler{}

type internalEchoHandler struct{}

func (handler internalEchoHandler) ServeTELNET(ctx Context, w Writer, r Reader) error {

	cmdstr, err := ReadLine(r)

	if err != nil {
		log.Println("读取命令失败:", err)
		return err
	}

	if cmdstr != "" {
		w.Write([]byte(cmd.Cmd(cmdstr)))
	}

}
