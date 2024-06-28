package telnet

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"strconv"
	"telnet/json"
)

func ListenAndServe(addr string, handler Handler) error {
	server := &Server{Addr: addr, Handler: handler}
	return server.ListenAndServe()
}

// Serve accepts an incoming TELNET or TELNETS client connection on the net.Listener `listener`.
func Serve(listener net.Listener, handler Handler) error {

	server := &Server{Handler: handler}
	return server.Serve(listener)
}

type Server struct {
	Addr    string  // TCP address to listen on; ":telnet" or ":telnets" if empty (when used with ListenAndServe or ListenAndServeTLS respectively).
	Handler Handler // handler to invoke; telnet.EchoServer if nil

	TLSConfig *tls.Config // optional TLS configuration; used by ListenAndServeTLS.

	Logger Logger
}

func (server *Server) ListenAndServe() error {

	addr := server.Addr
	if "" == addr {
		addr = ":telnet"
	}

	listener, err := net.Listen("tcp", addr)
	if nil != err {
		return err
	}

	log.Println("启动服务")

	return server.Serve(listener)
}

// Serve accepts an incoming TELNET client connection on the net.Listener `listener`.
func (server *Server) Serve(listener net.Listener) error {

	defer listener.Close()

	logger := server.logger()

	handler := server.Handler
	if nil == handler {
		//@TODO: Should this be a "ShellHandler" instead, that gives a shell-like experience by default
		//       If this is changd, then need to change the comment in the "type Server struct" definition.
		logger.Debug("Defaulted handler to EchoHandler.")
		handler = EchoHandler
	}

	for {
		// Wait for a new TELNET client connection.
		logger.Debugf("Listening at %q.", listener.Addr())
		conn, err := listener.Accept()
		if err != nil {
			//@TODO: Could try to recover from certain kinds of errors. Maybe waiting a while before trying again.
			return err
		}
		logger.Debugf("Received new connection from %q.", conn.RemoteAddr())

		// Handle the new TELNET client connection by spawning
		// a new goroutine.
		go server.handle(conn, handler)
		logger.Debugf("Spawned handler to handle connection from %q.", conn.RemoteAddr())
	}
}

func (server *Server) handle(c net.Conn, handler Handler) {
	defer c.Close()

	logger := server.logger()

	defer func() {
		if r := recover(); nil != r {
			if nil != logger {
				logger.Errorf("Recovered from: (%T) %v", r, r)
			}
		}
	}()

	var ctx = context.Background()

	var w Writer = newDataWriter(c)
	var r Reader = newDataReader(c)

	log.Println("建立连接啦")

	jsonH := &json.Logger{LogFile: "log/log.json"}

	DestIP, DestPort, _ := net.SplitHostPort(c.LocalAddr().String())
	DestPortInt, _ := strconv.Atoi(DestPort)

	clientIP, clientPort, _ := net.SplitHostPort(c.RemoteAddr().String())
	portInt, _ := strconv.Atoi(clientPort)

	log := json.LogEntry{
		Type:     "scan",
		DestIP:   DestIP,
		DestPort: DestPortInt,
		SrcIP:    clientIP,
		SrcPort:  portInt,
		// Extend:   extend,
		UUID:     "<UUID>",
		App:      "telnet",
		Name:     "telnet",
		Protocol: "HTTP",
	}

	jsonH.Log(log)

	handler.ServeTELNET(ctx, w, r)
	c.Close()
}

func (server *Server) logger() Logger {
	logger := server.Logger
	if nil == logger {
		logger = internalDiscardLogger{}
	}

	return logger
}
