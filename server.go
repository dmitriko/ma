package ma

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/codegangsta/martini"
//	"errors"
	"log"
	"net/http"
	"strings"
)

var (
	m *martini.ClassicMartini
	Sessions   = make(map[string]bool)
	AllowedIPs = make(map[string]bool)
	closeCh    = make(chan int)
	server     *Server
)

const (
	RemoteUrlPath = "/remote"
	DefaultPort   = 8000
)

type Server struct {
	Host      string
	Port      int
	IsTls     bool
	IsRunning bool
}

type ServerConfig struct {
	Host string
	Port int
	NoTls bool
}

func NewServer(c *ServerConfig) (error, *Server) {
	s := &Server{c.Host, c.Port, true, false}
	if c.NoTls {
		s.IsTls = false
	}
	if s.Host == "" {
		s.Host = "localhost"
	}
	if s.Port == 0 {
		s.Port = DefaultPort
	}
	server = s
	return nil, server
}

// return absolute url of base path, like
// https://10.12.196.11:8080/
func (s *Server) GetBaseUrl() (string) {
	if s.Host == "" || s.Port == 0 {
		panic("server is not setup")
	}
	proto := "https"
	if !s.IsTls {
		proto = "http"
	}
	return fmt.Sprintf("%s://%s:%d/", proto, s.Host, s.Port)
}

//set defaults for server

func (s *Server) Start() {
	m = martini.Classic()
	m.Get("/config", func() string {
			return "foo"
		})
	defer func() { s.IsRunning = false }()
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	http.Handle(RemoteUrlPath, websocket.Handler(remoteHandler))
	http.Handle("/", m)
	log.Printf("listen on %s\n", addr)
	err := http.ListenAndServeTLS(addr, "cert.pem", "key.pem", nil)
	if err != nil {
		panic("ListentAndServer: " + err.Error())
	}
}

func isWebsocketAllowed(ws *websocket.Conn) bool {
	addr := ws.Request().RemoteAddr
	t := strings.Split(addr, ":")
	if len(t) != 2 {
		log.Printf("ERROR: %s is wrong formated", addr)
		return false
	}
	return AllowedIPs[t[0]]

}

// handle websocket connection from remote host
func remoteHandler(ws *websocket.Conn) {
	// please note, ws will be closed after return in this func
	log.Printf("got websocket conn from %v\n", ws.Request().RemoteAddr)
	if !isWebsocketAllowed(ws) {
		log.Printf("WARN: websocket from %s is not allowed", ws.Request().RemoteAddr)
		websocket.Message.Send(ws, NOTOK)
		return
	}
	websocket.Message.Send(ws, OK)
	<-closeCh
}
