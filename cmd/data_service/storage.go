package main

import (
	"fmt"
	"time"

	"github.com/jdxj/study_object_storage/pkg/logger"
	"github.com/jdxj/study_object_storage/pkg/rabbit"
)

func NewStorage(host string, port int) *Storage {
	s := &Storage{
		host: host,
		port: port,
	}
	return s
}

type Storage struct {
	host string
	port int
}

func (s *Storage) Run() error {
	s.Heartbeat()

	// todo: 阻塞
	time.Sleep(10 * time.Minute)
	return nil
}

func (s *Storage) Heartbeat() {
	msg := &rabbit.Message{
		Headers: nil,
		Body:    []byte(fmt.Sprintf("%s:%d", s.host, s.port)),
	}

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
			}

			err := rabbit.Publish("data_service", msg)
			if err != nil {
				logger.Errorf("Publish: %s", err)
			}
			logger.Debugf("send heartbeat: %s", msg.Body)
		}
	}()
}
