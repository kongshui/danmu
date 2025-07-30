package common

import (
	"fmt"
	"os"
)

func GetSysLoad() (load1, load5, load15 float64) {
	content, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		fmt.Println("Error reading /proc/loadavg:", err)
		return
	}
	fmt.Sscanf(string(content), "%f %f %f", &load1, &load5, &load15)
	return
}
