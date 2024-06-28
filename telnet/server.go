package telnet

import (
	"crypto/tls"
	"log"
	"net"
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

	var w Writer = newDataWriter(c)
	var r Reader = newDataReader(c)

	log.Println("建立连接啦")

	json.GlobalLog.HoneyLog(c, "scan", nil)

	// 创建一个新的认证处理器
	authHandler := &AuthHandler{next: handler, users: server.UserDB}
	// 使用认证处理器代替原始处理器
	authHandler.ServeTELNET(c, w, r)
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

func (h *AuthHandler) ServeTELNET(c net.Conn, w Writer, r Reader) {

	log.Println("进来了")
	h.FailTime = 1
	w.Write([]byte("User Access Verification\r\n"))
	w.Write([]byte("Username:"))

	var err error
	for {

		//已登录
		if h.Succ {
			err = h.next.ServeTELNET(c, w, r)

			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}

			continue
		}

		// 读取用户名
		if h.Name == "" {

			h.Name, err = ReadLine(c, r)
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}

		}

		// 读取密码
		if h.Name != "" && h.Pwd == "" {

			if !h.PwdTip {
				w.Write([]byte("Password:"))
			}
			h.PwdTip = true

			h.Pwd, err = ReadLine(c, r)
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}

		}

		//登录校验
		if h.Name != "" && h.Pwd != "" {
			log.Println("账密", h.Name, h.Pwd)
			// 验证用户名和密码
			if storedPassword, ok := h.users[h.Name]; ok && storedPassword == h.Pwd {
				// 认证成功，调用下一个处理器

				if !h.Succ {
					w.Write([]byte("Authentication succeed.\r\n"))
					h.Succ = true
				}

				extend := make(map[string]any)

				extend["username"] = h.Name
				extend["password"] = h.Name
				extend["succ"] = true

				json.GlobalLog.HoneyLog(c, "login", extend)
			} else {

				extend := make(map[string]any)
				extend["username"] = h.Name
				extend["password"] = h.Name
				extend["succ"] = false
				extend["trytime"] = h.FailTime

				json.GlobalLog.HoneyLog(c, "login", extend)

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
