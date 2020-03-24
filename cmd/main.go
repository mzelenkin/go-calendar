package main

import (
	"fmt"
	"github.com/mzelenkin/go-calendar/internal/config"
	"log"
	"os"
)

const ConfigFilename = "./configs/config.yaml"

func main() {
	cfg, err := config.LoadConfig(ConfigFilename)
	exitIfError(err)

	logger, err := config.ConfigureLogging(&cfg.Log)
	exitIfError(err)

	fmt.Printf("Configuration loaded: %+v\n", *cfg)
	logger.Info("Hello")
}

func exitIfError(err error) {
	if err != nil {
		log.Fatal("Error: ", err.Error())
		os.Exit(1)
	}
}
