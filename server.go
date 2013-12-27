package ma

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"log"
	"net/http"
	"strings"
	"io/ioutil"
)

var (
	Sessions   = make(map[string]bool)
	AllowedIPs = make(map[string]bool)
	closeCh    = make(chan int)
	server     *Server
)

const (
	RemoteUrlPath = "/remote"
	DefaultPort   = 8000
	ConfigUrlPath = "/config"
)

type Server struct {
	Host      string
	Port      int
	IsTls     bool
	IsRunning bool
}

type ServerConfig struct {
	Host  string
	Port  int
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
// https://10.12.196.11:8080
func (s *Server) GetBaseUrl() string {
	if s.Host == "" || s.Port == 0 {
		panic("server is not setup")
	}
	proto := "https"
	if !s.IsTls {
		proto = "http"
	}
	return fmt.Sprintf("%s://%s:%d", proto, s.Host, s.Port)
}

//extract config data sent via POST as body or a parameter
func extractConfigValue(req *http.Request) ([]byte, error) {
	c := req.PostFormValue("config")
	if len(c) > 0 {
		return []byte(c), nil
	}
	
	data, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		return nil, err
	}
	return data, nil
}	

func configPOST(rw http.ResponseWriter, req *http.Request) {
	data, err := extractConfigValue(req)
	if err == nil {
		_, err = NewClusterConfig(string(data))
		if err == nil {
			rw.Write([]byte(OK))
			return
		}
	}
	log.Printf("could not handle POST to /config, got error %s", err)
	rw.WriteHeader(400)

}

func configGET(rw http.ResponseWriter, req *http.Request) {
	accept_header := req.Header["Accept"][0]
	var d []byte
	var e error
	switch accept_header {
	case "application/yaml":
		rw.Header().Set("Content-Type", "application/yaml")
		d, e = clusterConfig.Yaml()
	default:
		rw.Header().Set("Content-Type", "application/json")
		d, e = clusterConfig.Json()
	}
	if e != nil {
		log.Printf("ERROR: %s", e)
		rw.Write([]byte(""))
		return
	}
	rw.Write(d)

}

func ConfigHandler(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET": configGET(rw, req)
	case "POST": configPOST(rw, req)
	default: 
		rw.WriteHeader(400)
		rw.Write([]byte("Not Implemented"))
	}
}

func (s *Server) Start() {
	defer func() { s.IsRunning = false }()
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	http.Handle(RemoteUrlPath, websocket.Handler(remoteHandler))
	http.HandleFunc(ConfigUrlPath, ConfigHandler)
	log.Printf("listen on %s\n", addr)
	err := http.ListenAndServeTLS(addr, "cert.pem", "key.pem", nil)
	if err != nil {
		panic("ListentAndServer: " + err.Error())
	}
}

func isAgentAllowed(ws *websocket.Conn) bool {
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
	if !isAgentAllowed(ws) {
		log.Printf("WARN: websocket from %s is not allowed", ws.Request().RemoteAddr)
		websocket.Message.Send(ws, NOTOK)
		return
	}
	websocket.Message.Send(ws, OK)
	<-closeCh
}
