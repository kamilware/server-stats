package diskdata

import (
	"fmt"
	"sync"
	"time"

	"systemResources/floatutils"

	"github.com/shirou/gopsutil/v3/disk"
)

type Data struct {
	TotalAvailableGB float64 `json:"totalAvailableGB"`
	UsedPercent      float64 `json:"usedPercent"`
	UsedGB           float64 `json:"usedGB"`
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
	ticker := time.NewTicker(750 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mu.Lock()

			diskInfo, err := disk.Usage("/")
			if err != nil {
				fmt.Printf("error getting diskdata info: %v\n", err)
			}

			data.TotalAvailableGB = floatutils.ToFixed(floatutils.BytesToGB(diskInfo.Total), 2)
			data.UsedPercent = floatutils.ToFixed(diskInfo.UsedPercent, 2)
			data.UsedGB = floatutils.ToFixed(floatutils.BytesToGB(diskInfo.Used), 2)

			mu.Unlock()
		}
	}
}
