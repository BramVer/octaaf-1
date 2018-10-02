package main

import (
	"io"

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

	var err error
	var closer io.Closer
	Tracer, closer, err = cfg.New(service, config.Logger(jaegerClient.StdLogger))

	if err != nil {
		log.Panicf("Jaeger init error: %v", err)
	}

	return closer
}
