package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	yaml "gopkg.in/yaml.v2"

	"go.yhsif.com/vanity"
)

const configFile = "config.yaml"

// AppEngine log will auto add date and time, so there's no need to double log
// them in our own logger.
var logger = log.New(os.Stdout, "", log.Lshortfile)

func main() {
	cfg := loadConfig(configFile)

	http.Handle("/", vanity.Handler(vanity.Args{
		Config: cfg,
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		logger.Printf("Defaulting to port %s", port)
	}
	logger.Printf("Listening on port %s", port)

	logger.Print(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func loadConfig(path string) vanity.Config {
	f, err := os.Open(path)
	if err != nil {
		logger.Panic(err)
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	decoder.SetStrict(true)
	var cfg vanity.Config
	if err := decoder.Decode(&cfg); err != nil {
		logger.Panic(err)
	}
	logger.Printf("Config: %#v", cfg)
	return cfg
}
