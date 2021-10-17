package utilities

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func GetContextForSigHandling() context.Context {
	sigs := make(chan os.Signal, 1)
	ctx, cancel := context.WithCancel(context.Background())

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		log.Println("Signal caught, shutdown!")
		cancel()
	}()

	return ctx
}
