// test config processing

package ma

import "testing"

const (
	JSON_CONFIG = `
{"Name": "TestCluster",
"Description": "Test Cluster",
"Hosts": [{"RemoteIp":"127.0.0.1"}]
}`
)

func TestConfigNew(t *testing.T) {
	err, config := NewClusterConfig(JSON_CONFIG)
	if err != nil {
		t.Error(err)
	}
	if config.Name != "TestCluster" {
		t.Errorf("expect %s for config.Name, got `%s`", "TestCluster", config.Name)
	}
	if config.Hosts[0].RemoteIp != "127.0.0.1" {
		t.Errorf("expect %s for config.Hosts[0].RemoteIp, got '%s'", "127.0.0.1", config.Hosts[0].RemoteIp)
	}

}
