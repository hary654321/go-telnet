package telnet

import (
	"bufio"
	"context"
	"io"
	"log"
	"strings"
	"telnet/cmd"

	"github.com/reiver/go-oi"
)

// EchoHandler is a simple TELNET server which "echos" back to the client any (non-command)
// data back to the TELNET client, it received from the TELNET client.
var EchoHandler Handler = internalEchoHandler{}

type internalEchoHandler struct{}

func (handler internalEchoHandler) ServeTELNET(ctx context.Context, w io.Writer, r io.Reader) error {
	reader := bufio.NewReader(r) // 使用bufio.Reader以便更高效地读取数据

	for {
		line, err := reader.ReadString('\n') // 读取直到换行符
		log.Println(line)
		if err != nil {
			if err == io.EOF {
				// 客户端正常关闭连接
				log.Println("Client disconnected")
				return nil
			}
			// 发生其他错误，记录并返回
			log.Printf("Error reading from client: %v", err)
			return err
		}

		// 去除可能的回车符
		line = strings.TrimRight(line, "\r\n")

		// 检查是否为命令
		if cmd.IsCommand(line) { // 假设cmd包提供了IsCommand函数
			response := cmd.Cmd(line)                   // 处理命令
			_, err := oi.LongWrite(w, []byte(response)) // 回显结果
			if err != nil {
				log.Printf("Error writing to client: %v", err)
				return err
			}
		} else {
			// 非命令数据直接回显
			_, err := w.Write([]byte(line))
			if err != nil {
				log.Printf("Error writing to client: %v", err)
				return err
			}
			_, err = w.Write([]byte("\n")) // 添加换行符以保持格式
			if err != nil {
				log.Printf("Error writing to client: %v", err)
				return err
			}
		}

		// 检查上下文是否已被取消或设置了截止时间
		select {
		case <-ctx.Done():
			log.Println("Context cancelled, closing connection")
			return ctx.Err()
		default:
			// 继续处理下一个请求
		}
	}
}
