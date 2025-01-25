package cpudata

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"

	"systemResources/floatutils"
)

type Data struct {
	Model       string  `json:"model"`
	Cores       int     `json:"cores"`
	Threads     int     `json:"threads"`
	UsedPercent float64 `json:"usedPercent"`
	Temperature float64 `json:"temperature"`
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

			cpuInfo, err := cpu.Info()
			if err != nil || len(cpuInfo) == 0 {
				fmt.Printf("error getting CPU info: %v\n", err)
			}

			data.Model = cpuInfo[0].ModelName
			data.Cores = len(cpuInfo) / 2
			data.Threads = len(cpuInfo)

			usage, err := cpu.Percent(time.Second, false)
			if err != nil {
				fmt.Printf("Error getting CPU usage: %v\n", err)
			}

			data.UsedPercent = floatutils.ToFixed(usage[0], 2)

			temperature, err := getTemperature()
			if err != nil {
				fmt.Printf("error getting CPU temperature: %v\n", err)
			}

			data.Temperature = floatutils.ToFixed(temperature, 2)

			mu.Unlock()
		}
	}
}

func getTemperature() (float64, error) {
	cmd := exec.Command("sensors")
	out, err := cmd.Output()
	if err != nil {
		return 0.0, fmt.Errorf("error running sensors: %v", err)
	}

	re := regexp.MustCompile(`Tctl:\s*\+?(\d+\.\d+)°C`)
	matches := re.FindStringSubmatch(string(out))
	if len(matches) > 1 {
		temp, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			return 0.0, fmt.Errorf("error parsing CPU temperature: %v", err)
		}

		return temp, nil
	}

	return 0.0, fmt.Errorf("could not find CPU temperature in the output")
}
