//Place for ClusterConfig, HostConfig

package ma

import (
	"encoding/json"
	//	"errors"
	"github.com/BurntSushi/toml"
	yaml "github.com/gonuts/yaml"
	"strings"
	"log"
//	"fmt"
)

type ClusterConfig struct {
	Name        string `yaml:"name" json:"name" toml:"name"`
	Description string `yaml:"description" json:"description" toml:"description"`
	Hosts       []HostConfig `yaml: "hosts" json:"hosts" toml:"hosts"`
}

type HostConfig struct {
	RemoteIp string `toml:"remote_ip" yaml:"remote_ip" json:"remote_ip"`
}

func NewClusterConfig(data string) (err error, config *ClusterConfig) {
	data = strings.TrimSpace(data)

	if strings.HasPrefix(data, "{") {
		err = json.Unmarshal([]byte(data), &config)
		if err != nil {
			return err, nil
		}
		log.Print("loaded JSON config")
		return 
	}

	if strings.HasPrefix(data, "---") {
		err = yaml.Unmarshal([]byte(data), &config)
		if err != nil {
			return err, nil
		}
		log.Print("loaded YAML config")
		return 
	}

	_, err = toml.Decode(data, &config)
	if err != nil {
		return err, nil
	}
	log.Print("loaded TOML config")

	return
}
