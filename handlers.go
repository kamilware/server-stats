package server_stats

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"systemResources/cpudata"
	"systemResources/diskdata"
	"systemResources/memdata"
	"systemResources/netdata"
)

type Response struct {
	CPU     cpudata.Data  `json:"cpu"`
	RAM     memdata.Data  `json:"ram"`
	Disk    diskdata.Data `json:"disk"`
	Network netdata.Data  `json:"network"`
}

func statsHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)

	_, ok := w.(http.Flusher)
	if !ok {
		fmt.Println("Flusher not supported")
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)

		return
	}

	for {
		response := Response{
			CPU:     cpudata.Get(),
			RAM:     memdata.Get(),
			Disk:    diskdata.Get(),
			Network: netdata.Get(),
		}

		data, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error marshaling response: %v", err)
			break
		}

		_, _ = fmt.Fprintf(w, "data: %s\n\n", string(data))
		flusher, ok := w.(http.Flusher)
		if ok {
			flusher.Flush()
		}

		time.Sleep(1 * time.Second)
	}
}
