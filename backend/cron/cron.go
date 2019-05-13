package cron

import (
	"github.com/fancxxy/gocomic/backend/config"
	"github.com/robfig/cron"
)

// Init function create and start cron job
func Init() error {
	c := cron.New()
	if err := c.AddFunc(config.Config.Cron.Update, func() { Update() }); err != nil {
		return err
	}
	if err := c.AddFunc(config.Config.Cron.Clean, func() { Clean() }); err != nil {
		return err
	}

	c.Start()
	return nil
}
