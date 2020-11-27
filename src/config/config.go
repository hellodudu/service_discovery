package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"github.com/rs/zerolog/log"
)

type CombinedServices map[int32]*ServiceConfig

type BaseConfig interface {
	GetID() int32
}

type StartSceneConfig struct {
	ID        int32  `json:"_id"`
	Process   int32  `json:"Process"`
	Zone      int32  `json:"Zone"`
	SceneType string `json:"SceneType"`
	Name      string `json:"Name"`
	OuterPort int32  `json:"OuterPort"`
}

func (c *StartSceneConfig) GetID() int32 {
	return c.ID
}

type StartProcessConfig struct {
	ID        int32  `json:"_id"`
	MachineId int32  `json:"MachineId"`
	InnerPort string `json:"InnerPort"`
}

func (c *StartProcessConfig) GetID() int32 {
	return c.ID
}

type StartMachineConfig struct {
	ID       int32  `json:"_id"`
	InnerIP  string `json:"InnerIP"`
	OuterIP  string `json:"OuterIP"`
	HttpPort string `json:"HttpPort"`
}

func (c *StartMachineConfig) GetID() int32 {
	return c.ID
}

// 最终要导入consul中的服务配置
type ServiceConfig struct {
	*StartSceneConfig `json:",inline"`
	InnerIP           string `json:"InnerIP"`
	OuterIP           string `json:"OuterIP"`
	InnerPort         string `json:"InnerPort"`
	HttpPort          string `json:"HttpPort"`
}

type ConfigManager struct {
	mapSceneConfig   map[int32]BaseConfig
	mapProcessConfig map[int32]BaseConfig
	mapMachineConfig map[int32]BaseConfig

	mapCombinedService CombinedServices
}

func extractFromFile(path string, rtype reflect.Type) (values []BaseConfig, err error) {
	data, e := ioutil.ReadFile(path)
	if e != nil {
		dir, _ := os.Getwd()
		log.Error().Str("path", path).Str("dir", dir).Err(e).Msg("read file failed")
		err = e
		return
	}

	// remove last comma
	s := string(data)[:]
	n := strings.LastIndex(s, ",")
	news := strings.Join([]string{s[:n], s[n+1:]}, "")

	var datas [][]interface{}
	e = json.Unmarshal([]byte(news), &datas)
	if e != nil {
		log.Error().Err(e).Msg("json unmarshal failed")
		err = e
		return
	}

	for _, v := range datas {
		jsonBody := v[1]
		d, e := json.Marshal(jsonBody)
		if e != nil {
			log.Error().Err(e).Msg("marshal json failed")
			err = e
			return
		}

		config := reflect.New(rtype.Elem()).Interface().(BaseConfig)
		e = json.Unmarshal(d, config)
		if e != nil {
			log.Error().Err(e).Msg("unmarshal json failed")
			err = e
			return
		}

		values = append(values, config)
	}

	return
}

func NewConfigManager() *ConfigManager {
	return &ConfigManager{
		mapSceneConfig:     make(map[int32]BaseConfig),
		mapProcessConfig:   make(map[int32]BaseConfig),
		mapMachineConfig:   make(map[int32]BaseConfig),
		mapCombinedService: make(CombinedServices),
	}
}

func (cm *ConfigManager) GetCombinedService() CombinedServices {
	return cm.mapCombinedService
}

func (cm *ConfigManager) LoadFromFile(path string) {
	// load config txt
	listSceneConfigs, err := extractFromFile(fmt.Sprintf("%s/StartSceneConfig.txt", path), reflect.TypeOf(&StartSceneConfig{}))
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	listProcessConfigs, err := extractFromFile(fmt.Sprintf("%s/StartProcessConfig.txt", path), reflect.TypeOf(&StartProcessConfig{}))
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	listMachineConfigs, err := extractFromFile(fmt.Sprintf("%s/StartMachineConfig.txt", path), reflect.TypeOf(&StartMachineConfig{}))
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	decoratorFunc := func(mapConfig map[int32]BaseConfig, lst []BaseConfig) {
		for _, v := range lst {
			mapConfig[v.GetID()] = v
		}
	}

	decoratorFunc(cm.mapSceneConfig, listSceneConfigs)
	decoratorFunc(cm.mapProcessConfig, listProcessConfigs)
	decoratorFunc(cm.mapMachineConfig, listMachineConfigs)
}

func (cm *ConfigManager) CombineService() {
	for _, v := range cm.mapSceneConfig {
		service := &ServiceConfig{
			StartSceneConfig: v.(*StartSceneConfig),
		}

		processBaseConfig, ok := cm.mapProcessConfig[service.Process]
		if !ok {
			log.Fatal().Int32("process_id", service.Process).Msg("combine process config failed")
		}

		processConfig := processBaseConfig.(*StartProcessConfig)
		machineBaseConfig, ok := cm.mapMachineConfig[processConfig.MachineId]
		if !ok {
			log.Fatal().Int32("machine_id", processConfig.MachineId).Msg("combine machine config failed")
		}

		machineConfig := machineBaseConfig.(*StartMachineConfig)

		service.InnerIP = machineConfig.InnerIP
		service.OuterIP = machineConfig.OuterIP
		service.HttpPort = machineConfig.HttpPort
		service.InnerPort = processConfig.InnerPort

		cm.mapCombinedService[service.ID] = service
	}
}
