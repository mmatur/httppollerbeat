package beater

import (
	"fmt"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	cfg "github.com/mmatur/httppollerbeat/config"
)

type Httppollerbeat struct {
	done   chan struct{}
	config cfg.HostsConfig
	client publisher.Client
}

// Creates beater
func New(b *beat.Beat, rawConfig *common.Config) (beat.Beater, error) {

	config := cfg.HostsConfig{}
	if err := rawConfig.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Httppollerbeat{
		done:   make(chan struct{}),
		config: config,
	}
	return bt, nil
}

func (bt *Httppollerbeat) Run(b *beat.Beat) error {
	var poller *Poller

	logp.Info("httppollerbeat is running! Hit CTRL-C to stop it.")

	bt.client = b.Publisher.Connect()

	for i, h := range bt.config.Hosts {
		logp.Debug("httppollerbeat", "Creating poller #%v with URL: %v", i, h.Url)
		poller = CreatePoller(bt, h)
		go poller.Run()
	}

	for {
		select {
		case <-bt.done:
			return nil
		}
	}
}

func (bt *Httppollerbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
