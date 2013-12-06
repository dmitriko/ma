package ma

import (
	"fmt"
	"code.google.com/p/go.net/websocket"
	"net/http"
)

const (
	RemoteUrlPath = "/remote"
	DefaultPort = 8000
)

// handle websocket connection from remote hots
func remoteHandler(ws *websocket.Conn) {
	fmt.Printf("got websocket conn from %v\n", ws.Request().RemoteAddr)
}

func StartServer(host string) {
	addr := fmt.Sprintf("%s:%d", host, DefaultPort)
	http.Handle(RemoteUrlPath, websocket.Handler(remoteHandler))
	fmt.Printf("listen on %s\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic("ListentAndServer: " + err.Error())
	}
}
