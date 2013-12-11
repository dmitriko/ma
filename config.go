//Place for ClusterConfig, HostConfig

package ma

import (
	"encoding/json"
//	"errors"
	"github.com/BurntSushi/toml"
	"strings"
)

type ClusterConfig struct {
	Name        string
	Description string
	Hosts       []HostConfig
}

type HostConfig struct {
	RemoteIp string `toml:"remote_ip"`
}

func NewClusterConfig(data string) (error, *ClusterConfig) {
	data = strings.TrimSpace(data)
	var config *ClusterConfig
	if strings.HasPrefix(data, "{") {
		err := json.Unmarshal([]byte(data), &config)
		if err != nil {
			return err, nil
		}
		return nil, config
	}
	_, err := toml.Decode(data, &config)
	if err != nil {
		return err, nil
	}

	return nil, config
}
