//Place for ClusterConfig, HostConfig

package ma

import (
	"encoding/json"
	//	"errors"
	yaml "github.com/gonuts/yaml"
	"strings"
	"log"
//	"fmt"
)

type ClusterConfig struct {
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description" json:"description"`
	Hosts       []HostConfig `yaml: "hosts" json:"hosts"`
	Raw string `yaml:"-" json:"-"`
}

func (c *ClusterConfig) Yaml () ([]byte, error) {
	b, err := yaml.Marshal(c)
	if err != nil {
		return b, err
	}
	r := append([]byte("---\n"), b...)
	return r, nil

}

func (c *ClusterConfig) Json () ([]byte, error) {
	return json.MarshalIndent(c, "", " ")
}

var clusterConfig *ClusterConfig

type HostConfig struct {
	RemoteIp string `yaml:"remote_ip" json:"remote_ip"`
}

func NewClusterConfig(data string) (config *ClusterConfig, err error) {
	data = strings.TrimSpace(data)
	config = &ClusterConfig{Raw:data}
	clusterConfig = config  // set global var
	if strings.HasPrefix(data, "{") {
		err = json.Unmarshal([]byte(data), &config)
		if err != nil {
			return nil, err
		}
		log.Print("loaded JSON config")
		return 
	}

	err = yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		return nil, err
	}
	log.Print("loaded YAML config")
	return 
}

