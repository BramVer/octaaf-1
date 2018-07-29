package main

import (
	"gopkg.in/robfig/cron.v2"
)

// Cron executes all the cron jobs
var Cron *cron.Cron

func initCrons() {
	Cron = cron.New()
	// Cron func: ss mm hh
	Cron.AddFunc("01 37 13 * * *", func() { sendGlobal("1337") })
	Cron.AddFunc("30 38 13 * * *", func() { getLeetBlazers("1337") })
	Cron.AddFunc("01 20 16 * * *", func() { sendGlobal("420") })
	Cron.AddFunc("30 21 16 * * *", func() { getLeetBlazers("420") })
	Cron.AddFunc("45 59 23 * * *", setKaliCount)
	Cron.Start()
}
