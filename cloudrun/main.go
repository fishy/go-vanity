package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"go.yhsif.com/ctxslog"
	"golang.org/x/exp/slog"
	yaml "gopkg.in/yaml.v2"

	"go.yhsif.com/vanity"
)

const configFile = "config.yaml"

type config struct {
	Config vanity.Config `yaml:",inline"`

	IndexTemplate string `yaml:"index,omitempty"`
}

func main() {
	slog.SetDefault(slog.New(ctxslog.ContextHandler(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
		ReplaceAttr: ctxslog.ChainReplaceAttr(
			ctxslog.GCPKeys,
			ctxslog.StringDuration,
		),
	}))))

	cfg := loadConfig(configFile)

	if cfg.IndexTemplate != "" {
		var err error
		vanity.IndexTmpl, err = template.New("index").Parse(cfg.IndexTemplate)
		if err != nil {
			slog.Error(
				"Invalid index template",
				"err", err,
			)
			os.Exit(1)
		}
	}

	http.HandleFunc(
		"/_ah/health",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "healthy")
		},
	)
	http.Handle("/", vanity.Handler(vanity.Args{
		Config: cfg.Config,
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		slog.Warn(
			"Using default port",
			"port", port,
		)
	}
	slog.Info(
		"Started listening",
		"port", port,
	)

	slog.Info(
		"HTTP server returned",
		"err", http.ListenAndServe(fmt.Sprintf(":%s", port), nil),
	)
}

func loadConfig(path string) config {
	f, err := os.Open(path)
	if err != nil {
		slog.Error(
			"Unable to open config file",
			"err", err,
			"path", path,
		)
		os.Exit(1)
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	decoder.SetStrict(true)
	var cfg config
	if err := decoder.Decode(&cfg); err != nil {
		slog.Error(
			"Unable to decode config file",
			"err", err,
		)
		os.Exit(1)
	}

	slog.Info(
		"Loaded config",
		"config", cfg,
	)
	return cfg
}
