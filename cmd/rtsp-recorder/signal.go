package main

import (
	"golang.org/x/sys/unix"
	"os"
	"os/signal"
)

func StartSignalListener() chan os.Signal {
	c := make(chan os.Signal, 2)
	signal.Notify(c, unix.SIGINT, unix.SIGTERM)
	return c
}
