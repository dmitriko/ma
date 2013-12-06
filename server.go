package ma

import (
	"fmt"
	"code.google.com/p/go.net/websocket"
	"net/http"
	"log"
)

const (
	RemoteUrlPath = "/remote"
	DefaultPort = 8000
)

// handle websocket connection from remote hots
func remoteHandler(ws *websocket.Conn) {
	log.Printf("got websocket conn from %v\n", ws.Request().RemoteAddr)
}

func StartServer(host string) {
	addr := fmt.Sprintf("%s:%d", host, DefaultPort)
	http.Handle(RemoteUrlPath, websocket.Handler(remoteHandler))
	log.Printf("listen on %s\n", addr)
	err := http.ListenAndServeTLS(addr, "cert.pem", "key.pem", nil)
	if err != nil {
		panic("ListentAndServer: " + err.Error())
	}
}
