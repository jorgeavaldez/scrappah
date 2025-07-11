package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"scrappah/pkg"
)

func main() {
	s := make(chan os.Signal, 1)

	signal.Notify(s, syscall.SIGINT, syscall.SIGQUIT)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-s
		cancel()
	}()

	// TODO: use a reader here from the db
	conf, err := pkg.ParseConfig("")
}
