package telnet

import (
	"telnet/cmd"

	"github.com/reiver/go-oi"
)

// EchoHandler is a simple TELNET server which "echos" back to the client any (non-command)
// data back to the TELNET client, it received from the TELNET client.
var EchoHandler Handler = internalEchoHandler{}

type internalEchoHandler struct{}

func (handler internalEchoHandler) ServeTELNET(ctx Context, w Writer, r Reader) {

	var buffer [512]byte // 定义一个足够大的缓冲区
	p := buffer[:]

	for {
		n, err := r.Read(p)

		if n > 0 {
			// 检查是否包含回车字符
			if containsReturn(p[:n]) {
				cmdstr := string(p[:n])
				response := cmd.Cmd(cmdstr)       // 处理命令
				oi.LongWrite(w, []byte(response)) // 回显结果
			}
		}

		if nil != err {
			break
		}
	}
}

func containsReturn(data []byte) bool {
	for _, b := range data {
		if b == '\n' || b == '\r' {
			return true
		}
	}
	return false
}
