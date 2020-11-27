package main

import (
	"flag"
	"fmt"
	"io"
	"main/src/config"
	"main/src/export"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	rotate "gopkg.in/natefinch/lumberjack.v2"
)

var importerPath string
var exporterPath string

func init() {
	flag.StringVar(&importerPath, "import_path", "../test/scene", "待转换txt路径")
	flag.StringVar(&exporterPath, "export_path", "../test/consul", "输出service.json路径")
	initLogger("consul_convertor")
}

func initLogger(appName string) {
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

func main() {
	flag.Parse()

	// load config from file
	cm := config.NewConfigManager()
	cm.LoadFromFile(importerPath)

	// combine to service config
	cm.CombineService()
	log.Info().Interface("services", cm.GetCombinedService()).Msg("combine service success")

	// generate consul's service.json
	ce := export.NewConsulExporter()
	destPath := fmt.Sprintf("%s/service.json", exporterPath)
	ce.WriteServicesToFile(cm.GetCombinedService(), destPath)
}
