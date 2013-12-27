// test config processing

package ma

import (
	"testing"
	"strings"
)

const (
	JSON_CONFIG = `
{"name": "TestCluster",
"description": "Test Cluster",
"hosts": [{"remote_ip":"127.0.0.1"}]
}`
	TOML_CONFIG = `
name = "TestCluster"
description = "Test Cluster"
[[hosts]]
remote_ip = "127.0.0.1"
`

	YAML_CONFIG = `
---
name: TestCluster
description: Test Cluster
hosts:
- remote_ip: 127.0.0.1
`
)

func TestConfigNew(t *testing.T) {
	for _, data := range []string{JSON_CONFIG, TOML_CONFIG, YAML_CONFIG} {
		config, err := NewClusterConfig(data)
		if err != nil {
			t.Error(err)
			continue
		}
		if config.Name != "TestCluster" {
			t.Errorf("expect %s for config.Name, got `%s` from %s", 
				"TestCluster", config.Name, data)
			t.Errorf("%v", config)
			continue
		}
		if config.Hosts[0].RemoteIp != "127.0.0.1" {
			t.Errorf("expect %s for config.Hosts[0].RemoteIp, got '%s' from %s", 
				"127.0.0.1", config.Hosts[0].RemoteIp, data)
			continue
		}
		if (clusterConfig == nil) {
			t.Error("clusterConfig is not set")
			println(clusterConfig)
			continue
		}
		if clusterConfig.Raw != strings.TrimSpace(data) {
			t.Error("wrong configCluster.Raw", clusterConfig.Raw)
		}
	}
}
