package main

import (
	"time"

	"github.com/jdxj/study_object_storage/pkg/logger"
	"github.com/jdxj/study_object_storage/pkg/rabbit"
)

func NewAPI(host string, port int) *API {
	api := &API{
		host: host,
		port: port,
		sm:   &StorageManager{addresses: make(map[string]time.Time)},
	}
	return api
}

type API struct {
	host string
	port int

	sm *StorageManager
}

func (api *API) Run() error {
	err := api.SubscribeStorage()
	if err != nil {
		return err
	}
	api.RemoveExpiredStorage()

	// todo: 阻塞
	time.Sleep(10 * time.Minute)
	return nil
}

func (api *API) SubscribeStorage() error {
	h := func(msg *rabbit.Message) (*rabbit.Message, error) {
		api.sm.AddAddress(string(msg.Body))
		logger.Debugf("SubscribeStorage: %s", msg.Body)
		return nil, nil
	}
	return rabbit.Subscribe("data_service", h)
}

func (api *API) RemoveExpiredStorage() {
	interval := 10 * time.Second
	h := func(addr string, t time.Time) bool {
		if t.Add(interval).Before(time.Now()) {
			logger.Debugf("RemoveExpiredStorage: %s", addr)
			return true
		}
		return false
	}

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
			}

			api.sm.DelRange(h)
		}
	}()
}
