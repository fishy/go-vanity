package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/blendle/zapdriver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	yaml "gopkg.in/yaml.v2"

	"go.yhsif.com/vanity"
)

const configFile = "config.yaml"

type config struct {
	Config vanity.Config `yaml:",inline"`

	IndexTemplate string `yaml:"index,omitempty"`
}

func main() {
	zapcfg := zapdriver.NewProductionConfig()
	zapcfg.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	logger, err := zapcfg.Build(zapdriver.WrapCore())
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)

	cfg := loadConfig(configFile)

	if cfg.IndexTemplate != "" {
		var err error
		vanity.IndexTmpl, err = template.New("index").Parse(cfg.IndexTemplate)
		if err != nil {
			zap.S().Fatalw(
				"Invalid index template",
				"err", err,
			)
		}
	}

	http.HandleFunc(
		"/_ah/health",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "healthy")
		},
	)
	http.Handle("/", vanity.Handler(vanity.Args{
		Config: cfg.Config,
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

func loadConfig(path string) config {
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
	var cfg config
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
