package netdata

import (
	"fmt"
	"sync"
	"time"

	"systemResources/floatutils"

	"github.com/shirou/gopsutil/v3/net"
)

type Data struct {
	UploadKBs   float64 `json:"uploadKBs"`
	DownloadKBs float64 `json:"downloadKBs"`
}

var (
	data = Data{}
	mu   = &sync.Mutex{}
)

func Get() Data {
	mu.Lock()
	defer mu.Unlock()

	return data
}

func Maintain() {
	ticker := time.NewTicker(1 * time.Second) // need a slightly larger interval as 750ms is causing strange updating values
	defer ticker.Stop()

	var prevStats []net.IOCountersStat

	for {
		select {
		case <-ticker.C:
			mu.Lock()

			stats, err := net.IOCounters(true)
			if err != nil {
				fmt.Printf("error getting network stats: %v\n", err)
			}

			if len(stats) > 0 {
				var totalUpload, totalDownload uint64

				for _, networkStat := range stats {
					totalUpload += networkStat.BytesSent
					totalDownload += networkStat.BytesRecv
				}

				if len(prevStats) > 0 {
					var prevUpload, prevDownload uint64
					for _, prevStat := range prevStats {
						prevUpload += prevStat.BytesSent
						prevDownload += prevStat.BytesRecv
					}

					uploadDiff := floatutils.BytesToKB(totalUpload - prevUpload)
					downloadDiff := floatutils.BytesToKB(totalDownload - prevDownload)

					data.UploadKBs = floatutils.ToFixed(uploadDiff, 2)
					data.DownloadKBs = floatutils.ToFixed(downloadDiff, 2)
				}

				prevStats = stats
			}

			mu.Unlock()
		}
	}
}
