package ma

import (
	"code.google.com/p/go.net/websocket"
	"testing"
	"fmt"
)

var (
	origin = "http://localhost"
	url    = fmt.Sprintf("ws://localhost:%d%s", DefaultPort, RemoteUrlPath)
)


func TestWebsocketConn(t *testing.T) {
	go StartServer("127.0.0.1")
	_, err := websocket.Dial(url, "", origin)
	if err != nil {
		t.Error(err)
	}
}
