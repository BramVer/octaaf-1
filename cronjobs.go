package main

import (
	"octaaf/models"

	"gopkg.in/robfig/cron.v2"
)

func initCrons() {
	Cron = cron.New()
	Cron.AddFunc("37 13 * * *", func() { sendGlobal("1337") })
	Cron.AddFunc("20 16 * * *", func() { sendGlobal("420") })
	Cron.AddFunc("59 23 * * *", setKaliCount)
	Cron.Start()
}

func setKaliCount() {
	lastCount := models.MessageCount{}

	err := DB.Last(&lastCount)

	count := models.MessageCount{
		Count: KaliCount,
		Diff:  0,
	}

	if err == nil && lastCount.Count > 0 {
		count.Diff = (KaliCount - lastCount.Count)
	}

	DB.Save(&count)
}
