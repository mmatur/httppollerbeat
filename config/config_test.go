// +build !integration

package config

import (
	"path/filepath"
	"testing"

	"github.com/elastic/beats/libbeat/cfgfile"
	"github.com/stretchr/testify/assert"
	"time"
)

func TestReadConfig(t *testing.T) {
	absPath, err := filepath.Abs("../tests/files/")

	assert.NotNil(t, absPath)
	assert.Nil(t, err)

	config := &ConfigSettings{}

	err = cfgfile.Read(config, absPath+"/config.yml")
	assert.Nil(t, err)

	hosts := config.Httppollerbeat.Hosts
	assert.Equal(t, 1, len(hosts))

	assert.Equal(t, "https://url.com", hosts[0].Url)
	assert.Equal(t, "get", hosts[0].Method)
	assert.Equal(t, 10*time.Second, hosts[0].Schedule)
	assert.Equal(t, "body", hosts[0].Body)
	assert.Equal(t, 10*time.Second, hosts[0].Timeout)
	assert.Equal(t, 1, len(hosts[0].Headers))
	assert.Equal(t, 2, len(hosts[0].Fields))
	assert.Equal(t, "httppollerbeat", hosts[0].DocumentType)
	assert.Equal(t, 2, len(hosts[0].Datas))
	assert.Equal(t, "name1", hosts[0].Datas[0].Name)
	assert.Equal(t, "{.name1}", hosts[0].Datas[0].JsonPath)
	assert.Equal(t, "int", hosts[0].Datas[0].JsonType)

}
