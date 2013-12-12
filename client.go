// Web socket client
package ma

import (
	"code.google.com/p/go.net/websocket"
	"crypto/tls"
	"os"
	"log"
	"fmt"
	"net/http"
)

func WebsocketDial(url_ string) (ws *websocket.Conn, err error) {
	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("WARN: could not get hostname %s", err)
		hostname = "localhost"
	}
	origin := fmt.Sprintf("http://%s", hostname)
	config, err := websocket.NewConfig(url_, origin)
	if err != nil {
		return nil, err
	}
	config.TlsConfig = &tls.Config{InsecureSkipVerify: true}
	return websocket.DialConfig(config)

}

func NewHttpClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: tr}
}
