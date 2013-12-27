// Test manager http server

package ma

import (
	"code.google.com/p/go.net/websocket"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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

func testConfigHttpGet(accept_header, config_data string) error {
	var client = NewHttpClient()
	url := server.GetBaseUrl() + "/config"
	_, _ = NewClusterConfig(config_data)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/"+accept_header)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("GET /config %s", resp.Status))
	}
	content_type := resp.Header["Content-Type"][0]
	if !strings.Contains(content_type, "application/"+accept_header) {
		return errors.New(fmt.Sprintf(
			"expect application/%s from GET /config, got %s",
			accept_header, content_type))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	config, err := NewClusterConfig(string(body))
	if err != nil {
		return err
	}
	if len(config.Hosts) != 1 || config.Hosts[0].RemoteIp != "127.0.0.1" {
		return errors.New(
			fmt.Sprintf("got wrong response from /config `%s` for %s", 
				body, accept_header))
	}
	return nil
}

func TestConfigHttpGet(t *testing.T) {
	once.Do(startServer)
	m := make(map[string]string)
	m["json"] = JSON_CONFIG
	m["toml"] = TOML_CONFIG
	m["yaml"] = YAML_CONFIG
	for key, value := range m {
		err := testConfigHttpGet(key, value)
		if err != nil {
			t.Error(err)
		}
	}
}
