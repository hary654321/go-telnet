package telnet

import (
	"crypto/tls"
	"log"
	"net"
	"strconv"
	"telnet/json"
)

func NewServer(addr string, handler Handler) *Server {
	return &Server{
		Addr:    addr,
		Handler: handler,
		UserDB:  make(UserDatabase),
	}
}

// UserDatabase 存储用户名和密码的结构
type UserDatabase map[string]string
type Server struct {
	Addr    string  // TCP address to listen on; ":telnet" or ":telnets" if empty (when used with ListenAndServe or ListenAndServeTLS respectively).
	Handler Handler // handler to invoke; telnet.EchoServer if nil

	TLSConfig *tls.Config // optional TLS configuration; used by ListenAndServeTLS.

	Logger Logger
	UserDB UserDatabase // 用户数据库
}

// AddUser 向用户数据库添加用户
func (s *Server) AddUser(username, password string) {
	s.UserDB[username] = password
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

	var ctx Context = NewContext().InjectLogger(logger)

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
		Protocol: "telnet",
	}

	jsonH.Log(log)

	// 创建一个新的认证处理器
	authHandler := &AuthHandler{next: handler, users: server.UserDB}
	// 使用认证处理器代替原始处理器
	authHandler.ServeTELNET(ctx, w, r)
	c.Close()
}

// AuthHandler 负责处理认证
type AuthHandler struct {
	next     Handler      // 认证成功后的处理器
	users    UserDatabase // 用户数据库
	Name     string
	Pwd      string
	FailTime int
	Succ     bool
	PwdTip   bool
}

func (h *AuthHandler) ServeTELNET(ctx Context, w Writer, r Reader) {

	log.Println("进来了")
	w.Write([]byte("User Access Verification\r\n"))
	w.Write([]byte("Username:"))

	for {

		if h.Name == "" {
			// 读取用户名
			h.Name, _ = ReadLine(r)

		}

		if h.Name != "" && h.Pwd == "" {

			if !h.PwdTip {
				w.Write([]byte("Password:"))
			}
			h.PwdTip = true

			// 读取用户名
			h.Pwd, _ = ReadLine(r)
		}

		if h.Name != "" && h.Pwd != "" {
			log.Println("账密", h.Name, h.Pwd)
			// 验证用户名和密码
			if storedPassword, ok := h.users[h.Name]; ok && storedPassword == h.Pwd {
				// 认证成功，调用下一个处理器

				if !h.Succ {
					w.Write([]byte("Authentication succeed.\r\n"))
					h.Succ = true
				}

				err := h.next.ServeTELNET(ctx, w, r)

				if err != nil {
					log.Println(err)
					return
				}
			} else {
				// 认证失败，发送错误消息并关闭连接
				w.Write([]byte("Authentication failed.\r\n"))
				h.Name = ""
				h.Pwd = ""
				h.PwdTip = false
				w.Write([]byte("Username:"))

				h.FailTime++
			}
		}
	}
}

func (server *Server) logger() Logger {
	logger := server.Logger
	if nil == logger {
		logger = internalDiscardLogger{}
	}

	return logger
}
