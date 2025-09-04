package service

import "fmt"

func RecoverFunc() {
	if err := recover(); err != nil {
		ziLog.Error(fmt.Sprintf("recover %v", err), debug)
	}
}
