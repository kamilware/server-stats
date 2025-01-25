package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"systemResources/cpudata"
	"systemResources/diskdata"
	"systemResources/memdata"
	"systemResources/netdata"
)

var (
	shutdownCh = make(chan os.Signal, 1)
)

func main() {
	signal.Notify(shutdownCh, syscall.SIGTERM, syscall.SIGINT)

	go cpudata.Maintain()
	go memdata.Maintain()
	go diskdata.Maintain()
	go netdata.Maintain()

	http.HandleFunc("/stats", statsHandler)

	go func() {
		log.Println("Starting server on :8080...")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	<-shutdownCh
	fmt.Println("Shutting down after a delay...")
	time.Sleep(5 * time.Second)

	os.Exit(0)
}
