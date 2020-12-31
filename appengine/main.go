package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/blendle/zapdriver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	yaml "gopkg.in/yaml.v2"

	"go.yhsif.com/vanity"
)

const configFile = "config.yaml"

func main() {
	zapcfg := zapdriver.NewProductionConfig()
	zapcfg.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	logger, err := zapcfg.Build(zapdriver.WrapCore())
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)

	cfg := loadConfig(configFile)

	http.Handle("/", vanity.Handler(vanity.Args{
		Config: cfg,
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		zap.S().Warnw(
			"Using default port",
			"port", port,
		)
	}
	zap.S().Infow(
		"Started listening",
		"port", port,
	)

	zap.S().Errorw(
		"HTTP server returned",
		"err", http.ListenAndServe(fmt.Sprintf(":%s", port), nil),
	)
}

func loadConfig(path string) vanity.Config {
	f, err := os.Open(path)
	if err != nil {
		zap.S().Fatalw(
			"Unable to open config file",
			"path", path,
			"err", err,
		)
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	decoder.SetStrict(true)
	var cfg vanity.Config
	if err := decoder.Decode(&cfg); err != nil {
		zap.S().Fatalw(
			"Unable to decode config file",
			"err", err,
		)
	}

	zap.S().Infow(
		"Loaded config",
		"config", cfg,
	)
	return cfg
}
