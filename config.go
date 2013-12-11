//Place for ClusterConfig, HostConfig

package ma

import (
	"strings"
	"encoding/json"
)

type ClusterConfig struct {
	Name string
	Description string
	Hosts []HostConfig
}

type HostConfig struct {
	RemoteIp string
	HostName string
}


func NewClusterConfig(data string) (error, *ClusterConfig) {
	data = strings.TrimSpace(data)
	var config *ClusterConfig 
	if strings.HasPrefix(data, "{") {
		err := json.Unmarshal([]byte(data), &config)
		if err != nil {
			return err, nil
		}
	}
	return nil, config
}
