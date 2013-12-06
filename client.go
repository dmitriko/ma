// Web socket client 
package ma

import (
	"code.google.com/p/go.net/websocket"
	"crypto/tls"
)

func WebsocketDial(url_, origin string) (ws *websocket.Conn, err error) {
	config, err := websocket.NewConfig(url_, origin)
	if err != nil {
		return nil, err
	}
	config.TlsConfig = &tls.Config{InsecureSkipVerify: true}
	return websocket.DialConfig(config)

}
