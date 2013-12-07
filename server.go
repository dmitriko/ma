package ma

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	//	"github.com/codegangsta/martini"
	"log"
	"net/http"
	"strings"
)

var (
	//	m *martini.ClassicMartini
	Sessions   = make(map[string]bool)
	AllowedIPs = make(map[string]bool)
	closeCh = make(chan int)
)

const (
	RemoteUrlPath = "/remote"
	DefaultPort   = 8000
)

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
	log.Printf("got websocket conn from %v\n", ws.Request().RemoteAddr)
	if !isWebsocketAllowed(ws) {
		log.Printf("WARN: websocket from %s is not allowed", ws.Request().RemoteAddr)
		ws.Close()
		return
	}
	msg := "ok"
	websocket.Message.Send(ws, msg)
	<- closeCh
}

func StartServer(host string) {
	/*	m = martini.Classic()
		m.Get("/api/v1.0/", func() string {
			return "yo there"
		})
	*/
	addr := fmt.Sprintf("%s:%d", host, DefaultPort)
	http.Handle(RemoteUrlPath, websocket.Handler(remoteHandler))
	//	http.Handle("/", m)
	log.Printf("listen on %s\n", addr)
	err := http.ListenAndServeTLS(addr, "cert.pem", "key.pem", nil)
	if err != nil {
		panic("ListentAndServer: " + err.Error())
	}
}
