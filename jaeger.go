package main

import (
	"io"

	// tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	opentracing "github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	jaegerClient "github.com/uber/jaeger-client-go"
	config "github.com/uber/jaeger-client-go/config"
)

var Tracer opentracing.Tracer

func initJaeger(service string) io.Closer {
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}

	// opentracing.

	var err error
	var closer io.Closer
	Tracer, closer, err = cfg.New(service, config.Logger(jaegerClient.StdLogger))

	if err != nil {
		log.Panicf("Jaeger init error: %v", err)
	}

	return closer
}

// func getSpanCtx(message *tgbotapi.Message) {
// 	message.Chat.ID
// 	message.MessageID
// }
