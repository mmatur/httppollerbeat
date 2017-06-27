// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import (
	"time"
)

type ConfigSettings struct {
	Httppollerbeat HostsConfig
}

type HostsConfig struct {
	Hosts []HostConfig
}

const (
	DefaultSchedule     time.Duration = 1 * time.Minute
	DefaultTimeout      time.Duration = 60 * time.Second
	DefaultDocumentType string        = "httppollerbeat"
)

type HostConfig struct {
	Schedule     time.Duration     `config:"schedule"`
	Url          string            `config:"url"`
	Method       string            `config:"method"`
	Body         string            `config:"body"`
	Headers      map[string]string `config:"headers"`
	Timeout      time.Duration     `config:"timeout"`
	DocumentType string            `config:"document_type"`
	Fields       map[string]string `config:"fields"`
	Datas        []DataConfig      `config:"datas"`
}

type DataConfig struct {
	Name     string `config:"name"`
	JsonPath string `config:"json_path"`
	JsonType string `config:"json_type"`
}
