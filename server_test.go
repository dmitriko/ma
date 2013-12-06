package ma

import (
	"testing"
	"fmt"
)

var (
	origin = "http://localhost"
	url    = fmt.Sprintf("wss://localhost:%d%s", DefaultPort, RemoteUrlPath)
)


func TestWebsocketConn(t *testing.T) {
	go StartServer("127.0.0.1")
	_, err := WebsocketDial(url, origin)
	if err != nil {
		t.Error(err)
	}
}
