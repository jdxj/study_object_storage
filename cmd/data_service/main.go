package main

import (
	"log"

	"github.com/jdxj/study_object_storage/pkg/config"
	"github.com/jdxj/study_object_storage/pkg/logger"
	"github.com/jdxj/study_object_storage/pkg/rabbit"
)

func main() {
	conf, err := config.New("./conf.yaml")
	if err != nil {
		log.Fatalln(err)
	}

	loggerCfg := conf.Logger
	logger.Init(loggerCfg.FileName, loggerCfg.AppName, loggerCfg.MaxSize, loggerCfg.MaxAge, loggerCfg.MaxBackups,
		loggerCfg.Level, loggerCfg.LocalTime, loggerCfg.Compress)

	rabbitCfg := conf.Rabbit
	err = rabbit.Init(rabbitCfg.User, rabbitCfg.Pass, rabbitCfg.Host, rabbitCfg.Port)
	if err != nil {
		log.Fatalln(err)
	}

	webCfg := conf.Web
	storage := NewStorage(webCfg.Host, webCfg.Port)
	err = storage.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
