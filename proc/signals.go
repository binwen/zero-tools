package proc

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGTERM)

		for {
			v := <-signals
			switch v {
			case syscall.SIGTERM:
				GracefulStop(signals)
			default:
				log.Print("Got unregistered signal:", v)
			}
		}
	}()
}
