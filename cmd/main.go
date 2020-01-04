package main

import (
	"fmt"
	"github.com/mzelenkin/go-calendar/internal/config"
	"os"
)

const ConfigFilename = "./configs/config.yaml"

func main() {
	cfg, err := config.LoadConfig(ConfigFilename)
	exitIfError(err)

	log, err := config.ConfigureLogging(&cfg.Log)
	exitIfError(err)

	fmt.Printf("Configuration loaded: %+v\n", *cfg)
	log.Info("Hello")
}

func exitIfError(err error) {
	if err != nil {
		println("Error: ", err.Error())
		os.Exit(1)
	}
}
