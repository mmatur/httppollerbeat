package beater

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/mmatur/httppollerbeat/config"
	"k8s.io/client-go/pkg/util/jsonpath"
	"net/http"
	"strconv"
	"time"
)

type Poller struct {
	HttpPollerBeat *Httppollerbeat
	config         config.HostConfig
	schedule       time.Duration
	documentType   string
	timeout        time.Duration
}

func CreatePoller(hb *Httppollerbeat, config config.HostConfig) *Poller {
	poller := &Poller{
		HttpPollerBeat: hb,
		config:         config,
	}

	return poller
}

func (p *Poller) Run() {
	// setup default config
	p.schedule = config.DefaultSchedule
	p.documentType = config.DefaultDocumentType
	p.timeout = config.DefaultTimeout

	if p.config.DocumentType != "" {
		p.documentType = p.config.DocumentType
	}

	ticker := time.NewTicker(p.config.Schedule)
	for {
		select {
		case <-p.HttpPollerBeat.done:
			break
		case <-ticker.C:
			p.makeHttpCall()
		}

	}
}

func (p *Poller) makeHttpCall() error {
	url := p.config.Url
	method := p.config.Method
	client := http.Client{
		Timeout: p.config.Timeout,
	}
	body := p.config.Body
	headers := p.config.Headers

	request, err := p.newRequest(method, url, body)

	if err != nil || request == nil {
		return fmt.Errorf("Failed to create HTTP request [%s] [%s]: %s", method, url, err)
	}

	for k, v := range headers {
		request.Header.Set(k, v)
	}

	resp, err := client.Do(request)

	if err == nil {
		defer resp.Body.Close()
	}

	var jsonBody map[string]interface{}

	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	errs := decoder.Decode(&jsonBody)
	if errs != nil {
		jsonBody = nil
		logp.Err("httppollerbeat", "An error occurred while marshalling response to JSON: %w", errs)
	}

	var datas map[string]interface{}
	datas = make(map[string]interface{})

	for _, e := range p.config.Datas {
		value, err := parseJSONPath(jsonBody, e.Name, e.JsonPath, e.JsonType)
		logp.Info("value name : %s json path: %s value %s", e.Name, e.JsonPath, value)
		if err != nil {
			return fmt.Errorf("error parsing %s: %v", e.Name, err)
		}
		datas[e.Name] = value
	}

	event := Event{
		ReadTime:     time.Now(),
		DocumentType: p.documentType,
		Fields:       p.config.Fields,
		Datas:        datas,
		Url:          p.config.Url,
	}
	p.HttpPollerBeat.client.PublishEvent(event.ToMapStr())

	return nil
}

func parseJSONPath(input interface{}, name, template string, jsonType string) (interface{}, error) {
	j := jsonpath.New(name)
	buf := new(bytes.Buffer)
	if err := j.Parse(template); err != nil {
		return "", err
	}
	if err := j.Execute(buf, input); err != nil {
		return "", err
	}
	switch jsonType {
	case "int":
		parseResult, parseError := strconv.ParseInt(buf.String(), 10, 64)
		if parseError != nil {
			logp.Err("parseError %v", parseError)
			return "", parseError
		}
		return parseResult, nil
	case "float":
		parseResult, parseError := strconv.ParseFloat(buf.String(), 64)
		if parseError != nil {
			return "", parseError
		}
		return parseResult, nil
	case "boolean":
		parseResult, parseError := strconv.ParseBool(buf.String())
		if parseError != nil {
			return "", parseError
		}
		return parseResult, nil
	default:
		return buf.String(), nil
	}
}

func (p *Poller) newRequest(method string, url string, body string) (*http.Request, error) {
	switch method {
	case "get":
		return http.NewRequest("GET", url, nil)
	case "delete":
		return http.NewRequest("DELETE", url, nil)
	case "head":
		return http.NewRequest("HEAD", url, nil)
	case "patch":
		return http.NewRequest("PATCH", url, bytes.NewBufferString(body))
	case "post":
		return http.NewRequest("POST", url, bytes.NewBufferString(body))
	case "put":
		return http.NewRequest("PUT", url, bytes.NewBufferString(body))
	default:
		return nil, fmt.Errorf("Unsupported HTTP method %s", method)
	}
}
