// Test manager http server

package ma

import (
	"code.google.com/p/go.net/websocket"
	"errors"
	"fmt"
	"sync"
	"testing"
)

var (
	origin = "http://localhost"
	url    = fmt.Sprintf("wss://localhost:%d%s", DefaultPort, RemoteUrlPath)
	once   sync.Once
	wg sync.WaitGroup
)

func startServer() {
	go StartServer("127.0.0.1")
}

func _testWebsocketAllowed() error {
	if !AllowedIPs["127.0.0.1"] {
		panic("ouch")
	}
	ws, err := WebsocketDial(url, origin)
	if err != nil {
		return err
	}
	var msg string
	err = websocket.Message.Receive(ws, &msg)
	if err != nil {
		return err
	}
	if msg != "ok" {
		return errors.New("did't get ok from server")
	}
	fmt.Println("OK")
	return nil
}

func TestWebsocketAllowed(t *testing.T) {
	once.Do(startServer)
	AllowedIPs["127.0.0.1"] = true
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			err := _testWebsocketAllowed()
			if err != nil {
				t.Error(err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
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
