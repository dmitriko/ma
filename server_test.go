// Test manager http server

package ma

import (
	"code.google.com/p/go.net/websocket"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
	"net/http"
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

func TestConfigHttpGet(t *testing.T) {
	once.Do(startServer)
	var client = NewHttpClient()
	url := server.GetBaseUrl() + "/config"
	for _, accept_header := range []string{"json", "toml", "yaml"} {
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("Accept", "application/"+accept_header)
		resp, err := client.Do(req)
		if err != nil {
			t.Error(err)
		}
		if resp.StatusCode != 200 {
			t.Errorf("GET /config %s", resp.Status)
		}
		content_type := resp.Header["Content-Type"][0]
		if !strings.Contains(content_type, "application/" + accept_header) {
			t.Errorf("expect application/%s from GET /config, got %s", 
				accept_header, content_type)
		}

	}

}
