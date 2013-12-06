package ma

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"sync"
	"testing"
)

var (
	origin = "http://localhost"
	url    = fmt.Sprintf("wss://localhost:%d%s", DefaultPort, RemoteUrlPath)
	once   sync.Once
)

func startServer() {
	go StartServer("127.0.0.1")
}

func TestWebsocketAllowed(t *testing.T) {
	once.Do(startServer)
	AllowedIPs["127.0.0.1"] = true
	ws, err := WebsocketDial(url, origin)
	if err != nil {
		t.Error(err)
	}
	var msg string
	err = websocket.Message.Receive(ws, &msg)
	if err != nil {
		t.Error(err)
	}
	if msg != "ok" {
		t.Error("did't get ok from server")
	}
}

func TestWebsocketNotAllowed(t *testing.T) {
	once.Do(startServer)
	AllowedIPs["127.0.0.1"] = false
	ws, err := WebsocketDial(url, origin)
	if err != nil {
		t.Error(err)
	}
	var msg string
	err = websocket.Message.Receive(ws, &msg)

	if err == nil || msg == "ok" {
		t.Error("get ok from server, but should not")
	}
}
