package cron

import "github.com/robfig/cron"

// C 可调用的调度把
var C *cron.Cron

func init() {
	C = cron.New()
}
