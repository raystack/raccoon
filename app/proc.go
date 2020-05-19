package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var stopchan chan bool
var wg sync.WaitGroup

func Run() error {

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	stopchan = make(chan bool, 1)

	StartHTTPServer()
	waitForServerShutdown()
	return nil
}

func waitForServerShutdown() {
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for {
		sig := <-signalChan
		switch sig {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			log.Print(fmt.Sprintf("Received a signal %s", sig))
			log.Println("waiting for grace period before shutdown")
			time.Sleep(3 * time.Second)
			os.Exit(0)
		default:
			log.Print(fmt.Sprintf("Received a unexpected signal %s", sig))
		}
	}
}