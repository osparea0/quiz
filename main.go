package main

import (
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	hostname, err := os.Hostname()
	if err != nil {
		logger.Error("error getting hostname for local system", hostname, err.Error())
		return
	}

	logger.Info("starting quiz application on", "hostname:", hostname)

}
