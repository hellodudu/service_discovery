package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	rotate "gopkg.in/natefinch/lumberjack.v2"
)

type StartSceneConfig struct {
	ID        int32  `json:"_id"`
	Process   int32  `json:"Process"`
	Zone      int32  `json:"Zone"`
	SceneType string `json:"SceneType"`
	Name      string `json:"Name"`
	OuterPort int32  `json:"OuterPort"`
}

type StartProcessConfig struct {
	ID        int32  `json:"_id"`
	MachineId int32  `json:"MachineId"`
	InnerPort string `json:"InnerPort"`
}

type StartMachineConfig struct {
	ID      int32  `json:"_id"`
	InnerIP string `json:"InnerIP"`
	OuterIP string `json:"OuterIP"`
}

func InitLogger(appName string) {
	// log file name
	t := time.Now()
	fileTime := fmt.Sprintf("%d-%d-%d %d-%d-%d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	logFn := fmt.Sprintf("../data/log/%s_%s.log", appName, fileTime)

	// set console writer and file rotate writer
	log.Logger = log.Output(io.MultiWriter(zerolog.ConsoleWriter{Out: os.Stdout}, &rotate.Logger{
		Filename:   logFn,
		MaxSize:    200, // megabytes
		MaxBackups: 3,
		MaxAge:     15, //days
	})).With().Caller().Logger()
}

func ExtractFromFile(path string, values []interface{}) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error().Str("path", path).Err(err).Msg("read file failed")
		return err
	}

	// remove last comma
	s := string(data)[:]
	n := strings.LastIndex(s, ",")
	news := strings.Join([]string{s[:n], s[n+1:]}, "")

	var datas [][]interface{}
	err = json.Unmarshal([]byte(news), &datas)
	if err != nil {
		log.Error().Err(err).Msg("json unmarshal failed")
		return err
	}

	for _, v := range datas {
		jsonBody := v[1]
		d, err := json.Marshal(jsonBody)
		if err != nil {
			log.Error().Err(err).Msg("marshal json failed")
			return err
		}

		var config interface{}
		err = json.Unmarshal(d, &config)
		if err != nil {
			log.Error().Err(err).Msg("unmarshal json failed")
			return err
		}

		log.Info().Interface("config", config).Msg("unmarshal success")
		values = append(values, config)
	}

	return nil
}

func main() {
	InitLogger("importer")

	listSceneConfigs := make([]interface{}, 0)
	listProcessConfigs := make([]interface{}, 0)
	listMachineConfigs := make([]interface{}, 0)
	ExtractFromFile("../StartSceneConfig.txt", listSceneConfigs)
	ExtractFromFile("../StartProcessConfig.txt", listProcessConfigs)
	ExtractFromFile("../StartMachineConfig.txt", listMachineConfigs)
}
