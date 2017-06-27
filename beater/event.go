package beater

import (
	"github.com/elastic/beats/libbeat/common"
	"time"
)

type Event struct {
	ReadTime     time.Time
	DocumentType string
	Fields       map[string]string
	Datas        map[string]interface{}
	Url          string
}

func (h *Event) ToMapStr() common.MapStr {
	if h.Fields != nil {
		for i, f := range h.Fields {
			h.Datas[i] = f
		}
	}
	event := common.MapStr{
		"@timestamp": common.Time(h.ReadTime),
		"type":       h.DocumentType,
		"url":        h.Url,
	}

	for k, v := range h.Datas {
		event.Put(k, v)
	}

	return event
}
