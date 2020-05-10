package cron

import "github.com/robfig/cron"

var c *cron.Cron

func init(){
	c = cron.New()
}
