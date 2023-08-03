package recorder

import (
	"fmt"
	"github.com/Miguel-Dorta/rtsp-recorder/pkg/ffmpeg"
	"github.com/Miguel-Dorta/rtsp-recorder/pkg/http_push"
	"github.com/Miguel-Dorta/rtsp-recorder/pkg/log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	fileExt    = ".mkv"
	tmpFileExt = ".tmp" + fileExt
)

type Recorder struct {
	instance      *ffmpeg.Instance
	filePathNoExt string
}

func Start(savingPath, alias, url, proto string, recDuration, timeout time.Duration, verbose bool) (*Recorder, error) {
	now := time.Now()
	parentDir := filepath.Join(savingPath, alias, strconv.Itoa(now.Year()), fmt.Sprintf("%02d", now.Month()), fmt.Sprintf("%02d", now.Day()))
	if err := os.MkdirAll(parentDir, 0770); err != nil {
		return nil, fmt.Errorf("error creating recording parent directory '%s': %w", parentDir, err)
	}
	filePathNoExt := filepath.Join(parentDir, fmt.Sprintf("%d-%02d-%02d_%02d-%02d-%02d",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second()))

	instance, err := ffmpeg.StartRecording(url, proto, filePathNoExt+tmpFileExt, recDuration, timeout, verbose)
	if err != nil {
		return nil, fmt.Errorf("error starting ffmpeg recording: %w", err)
	}

	return &Recorder{
		instance:      instance,
		filePathNoExt: filePathNoExt,
	}, nil
}

func (r *Recorder) Wait(quit <-chan os.Signal) (exit bool) {
	var err error
	select {
	case <-quit:
		exit = true
		r.instance.Stop()
		err = <-r.instance.ExitChannel
	case err = <-r.instance.ExitChannel:
	}
	if err != nil && !strings.HasPrefix(err.Error(), "exit status") {
		log.Errorf("ffmpeg instance returned an error: %s", err)
	}

	if err = os.Rename(r.filePathNoExt+tmpFileExt, r.filePathNoExt+fileExt); err != nil {
		log.Errorf("error removing tmp extension from file '%s': %s", r.filePathNoExt, err)
		return
	}
	http_push.Report(true, "")
	return
}
