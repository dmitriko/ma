// test config processing

package ma

import (
	"testing"
)

const (
	JSON_CONFIG = `
{"Name": "TestCluster",
"Description": "Test Cluster",
"Hosts": [{"RemoteIp":"127.0.0.1"}]
}`
	TOML_CONFIG = `
Name = "TestCluster"
Description = "Test Cluster"
[[hosts]]
remote_ip = "127.0.0.1"
`
)

func TestConfigNew(t *testing.T) {
	for _, data := range []string{JSON_CONFIG, TOML_CONFIG} {
		err, config := NewClusterConfig(data)
		if err != nil {
			t.Error(err)
			continue
		}
		if config.Name != "TestCluster" {
			t.Errorf("expect %s for config.Name, got `%s`", "TestCluster", config.Name)
			continue
		}
		if config.Hosts[0].RemoteIp != "127.0.0.1" {
			t.Errorf("expect %s for config.Hosts[0].RemoteIp, got '%s'", "127.0.0.1", config.Hosts[0].RemoteIp)
		}
	}
}
