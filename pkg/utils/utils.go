package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// NewRecordFilepath returns a new filepath for a record with the savingPath and alias provided. It creates its parent dirs.
func NewRecordFilepath(savingPath, alias string) (string, error) {
	now := time.Now()
	path := filepath.Join(savingPath, alias, strconv.Itoa(now.Year()), itoaTwoChars(int(now.Month())), itoaTwoChars(now.Day()))
	if err := os.MkdirAll(path, 0770); err != nil {
		return "", fmt.Errorf("error creating dirs (%s): %w", path, err)
	}
	return filepath.Join(path, fmt.Sprintf("%d-%02d-%02d_%02d-%02d-%02d",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())), nil
}

// NewChannelWithTimeout returns a channel that will close in the time specified.
func NewChannelWithTimeout(d time.Duration) chan struct{} {
	c := make(chan struct{})
	go func() {
		time.Sleep(d)
		close(c)
	}()
	return c
}

func TernaryOperator[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}

func itoaTwoChars(i int) string {
	if i > 9 {
		return strconv.Itoa(i)
	}
	return "0" + strconv.Itoa(i)
}
