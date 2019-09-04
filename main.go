package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {
	w := newWorker()
	w.Start()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	w.Stop()
}
