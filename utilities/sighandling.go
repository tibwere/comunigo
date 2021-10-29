package utilities

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Funzione che permette di creare un contesto da usare
// per cancellare le attività nel momento in cui
// un segnale fra SIGINT e SIGTERM viene catturato
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
