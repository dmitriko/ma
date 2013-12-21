// Test manager http server

package ma

import (
	"code.google.com/p/go.net/websocket"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

var (
	url  = fmt.Sprintf("wss://localhost:%d%s", DefaultPort, RemoteUrlPath)
	once sync.Once
	wg   sync.WaitGroup
)

func startServer() {
	_, s := NewServer(&ServerConfig{})
	go s.Start()
	time.Sleep(100 * time.Millisecond)
}

func testWebsocketAllowed() error {
	if !AllowedIPs["127.0.0.1"] {
		return errors.New("127.0.0.1 is not in AllowedIPs")
	}
	ws, err := WebsocketDial(url)
	if err != nil {
		return err
	}
	var msg string
	err = websocket.Message.Receive(ws, &msg)
	if err != nil {
		return err
	}
	if msg != OK {
		return errors.New("did't get ok from server")
	}
	return nil
}

func TestWebsocketAllowed(t *testing.T) {
	once.Do(startServer)
	AllowedIPs["127.0.0.1"] = true
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			err := testWebsocketAllowed()
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
	ws, err := WebsocketDial(url)
	if err != nil {
		t.Error(err)
	}
	var msg string
	err = websocket.Message.Receive(ws, &msg)
	if msg == OK {
		t.Error("got ok from server, but should not")
	}
}
