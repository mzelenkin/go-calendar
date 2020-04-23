package main

import (
	"github.com/mzelenkin/go-calendar/cmd"
	"github.com/mzelenkin/go-calendar/internal/logging"
)

// main точка входа в приложение
func main() {
	logging.Exit(cmd.Execute())
}
