package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/rs/zerolog/log"
)

type ConsulServiceCheckEntry struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name"`
	Http     string                 `json:"http"`
	Method   string                 `json:"method"`
	Header   map[string]interface{} `json:"header"`
	Body     string                 `json:"body"`
	Interval string                 `json:"interval"`
	Timeout  string                 `json:"timeout"`
}

// service.json single service
type ConsulServiceEntry struct {
	ID      string                   `json:"id"`
	Name    string                   `json:"name"`
	Tags    []string                 `json:"tags"`
	Address string                   `json:"address"`
	Port    int                      `json:"port"`
	Check   *ConsulServiceCheckEntry `json:"check"`
}

// service.json single watch
type ConsulWatchEntry struct {
	Type              string                             `json:"type"`
	Key               string                             `json:"key,omitempty"`
	Service           string                             `json:"service,omitempty"`
	PassingOnly       bool                               `json:"passingonly,omitempty"`
	HandlerType       string                             `json:"handler_type"`
	HttpHandlerConfig *ConsulWatchHttpHandlerConfigEntry `json:"http_handler_config"`
}

// service.json single http_handler_config
type ConsulWatchHttpHandlerConfigEntry struct {
	Path          string `json:"path"`
	Method        string `json:"method"`
	Timeout       string `json:"timeout"`
	TlsSkipVerify bool   `json:"tls_skip_verify"`
}

type ConsulServiceExporter struct {
	Services []*ConsulServiceEntry `json:"services"`
	Watches  []*ConsulWatchEntry   `json:"watches"`
}

type ConsulExporter struct {
	cse *ConsulServiceExporter
}

func NewConsulExporter() *ConsulExporter {
	return &ConsulExporter{
		cse: &ConsulServiceExporter{
			Services: make([]*ConsulServiceEntry, 0),
			Watches:  make([]*ConsulWatchEntry, 0),
		},
	}
}

func (ce *ConsulExporter) WriteServicesToFile(configs map[int32]*ServiceConfig, path string) {
	for _, config := range configs {
		service := &ConsulServiceEntry{
			ID:      config.Name,
			Name:    config.SceneType,
			Tags:    []string{config.SceneType},
			Address: config.InnerIP,
		}

		var err error
		service.Port, err = strconv.Atoi(config.InnerPort)
		if err != nil {
			log.Fatal().Err(err).Msg("convert port to int failed")
		}

		service.Check = &ConsulServiceCheckEntry{
			ID:     fmt.Sprintf("http_check_%s", config.Name),
			Name:   fmt.Sprintf("http_check_%s", config.SceneType),
			Http:   fmt.Sprintf("http://%s:%s/%s", config.InnerIP, config.InnerPort, "health_check"),
			Method: "GET",
			Header: map[string]interface{}{
				"Content-Type": []string{"application/json"},
			},
			Body:     fmt.Sprintf("{\"service_id\":\"%s\"}", config.Name),
			Interval: "10s",
			Timeout:  "3s",
		}

		watchKey := &ConsulWatchEntry{
			Type:        "key",
			Key:         "service_config",
			HandlerType: "http",
			HttpHandlerConfig: &ConsulWatchHttpHandlerConfigEntry{
				Path:          fmt.Sprintf("http://%s:%s/watch_key", config.InnerIP, config.InnerPort),
				Method:        "GET",
				Timeout:       "10s",
				TlsSkipVerify: false,
			},
		}

		watchService := &ConsulWatchEntry{
			Type:        "service",
			Service:     config.SceneType,
			PassingOnly: false,
			HandlerType: "http",
			HttpHandlerConfig: &ConsulWatchHttpHandlerConfigEntry{
				Path:          fmt.Sprintf("http://%s:%s/watch_key", config.InnerIP, config.InnerPort),
				Method:        "GET",
				Timeout:       "10s",
				TlsSkipVerify: false,
			},
		}

		ce.cse.Services = append(ce.cse.Services, service)
		ce.cse.Watches = append(ce.cse.Watches, watchKey)
		ce.cse.Watches = append(ce.cse.Watches, watchService)
	}

	data, err := json.Marshal(ce.cse)
	if err != nil {
		log.Fatal().Err(err).Msg("json marshal failed")
	}

	err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
		log.Fatal().Err(err).Msg("write service.json failed")
	}
}

// test
func (ce *ConsulExporter) UnmarshalToStruct() {
	data, err := ioutil.ReadFile("../config/consul/service.json")
	if err != nil {
		log.Fatal().Err(err).Msg("read service.json failed")
	}

	err = json.Unmarshal(data, ce.cse)
	if err != nil {
		log.Fatal().Err(err).Msg("unmarshal json failed")
	}

	log.Info().Interface("consul service", ce.cse).Msg("unmarshal success")
}
