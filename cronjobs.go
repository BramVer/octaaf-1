package main

import (
	"octaaf/models"

	"github.com/jasonlvhit/gocron"
)

func initCrons() {
	gocron.Every(1).Day().At("13:37").Do(sendGlobal, "1337")
	gocron.Every(1).Day().At("14:14").Do(sendGlobal, "🐔 Daar is hij weer. Het kip. Mooi. \nhttps://youtu.be/r_qOBZcQqWo")
	gocron.Every(1).Day().At("16:20").Do(sendGlobal, "420")
	gocron.Every(1).Day().At("23:59").Do(setKaliCount)

	<-gocron.Start()
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
