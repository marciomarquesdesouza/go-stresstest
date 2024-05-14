package internal

import (
	"errors"
	"fmt"
	"time"
)

func logMessage(format string, a ...any) (n int, err error) {
	timeString := time.Now().UTC().Format("2006-01-02 15:04:05")
	args := []any{timeString, APP_NAME}
	args = append(args, a...)
	return fmt.Printf("%s [%s] "+format+"\n", args...)
}

func logMessageWithPrefix(prefix string, format string, a ...any) (n int, err error) {
	timeString := time.Now().UTC().Format("2006-01-02 15:04:05")
	args := []any{timeString, APP_NAME, prefix}
	args = append(args, a...)
	return fmt.Printf("%s [%s]%s "+format+"\n", args...)
}

func getLoadPerRunner(requests int64, runners int64) (*[]int64, error) {
	if requests < 0 {
		return nil, errors.New("valor inválido para requisições")
	}

	if runners < 0 {
		return nil, errors.New("valor inválido para runners")
	}

	baseLoad := requests / runners
	remainder := requests % runners

	loadPerRunner := make([]int64, runners)
	for i := int64(0); i < runners; i++ {
		loadPerRunner[i] = baseLoad
		if i < remainder {
			loadPerRunner[i]++
		}
	}

	return &loadPerRunner, nil
}
