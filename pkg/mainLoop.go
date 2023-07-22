package pkg

import (
	"fmt"
	"github.com/Miguel-Dorta/rtsp-recorder/pkg/ffmpeg"
	"github.com/Miguel-Dorta/rtsp-recorder/pkg/http_push"
	"github.com/Miguel-Dorta/rtsp-recorder/pkg/log"
	"github.com/Miguel-Dorta/rtsp-recorder/pkg/utils"
	"golang.org/x/sys/unix"
	"os"
	"os/signal"
	"time"
)

type config struct {
	Alias             string
	ClosingTimeout    time.Duration
	Path              string
	RecordingDuration time.Duration
	Url               string
	Verbose           bool
}

var Conf config

// Main does the main logic of the application
func Main() {
	quit := make(chan os.Signal, 2)
	signal.Notify(quit, unix.SIGINT, unix.SIGTERM)

	var (
		instance          *ffmpeg.Instance
		currentRecordPath string
		err               error
	)
	instance, currentRecordPath, err = start()
	if err != nil {
		log.Criticalf(err.Error())
	}

	for {
		select {
		// Case: ffmpeg ends before the timeout
		case err := <-instance.ExitChan:
			postRecordOps(currentRecordPath, false)
			if err != nil {
				log.Criticalf("ffmpeg instance returned an error: %s", err)
			}
			instance, currentRecordPath, err = start()
			if err != nil {
				log.Critical(err.Error())
			}
		// Case: end of the time of the recording
		case <-utils.NewChannelWithTimeout(Conf.RecordingDuration):
			// Start new recording before closing new one to have overlapping videos
			newInstance, newRecordPath, err := start()
			if err != nil {
				stop(instance)
				postRecordOps(currentRecordPath, false)
				log.Critical(err.Error())
			}
			stop(instance)
			postRecordOps(currentRecordPath, true)
			instance = newInstance
			currentRecordPath = newRecordPath
		// Case: User or system quit
		case <-quit:
			stop(instance)
			postRecordOps(currentRecordPath, false)
			os.Exit(0)
		}
	}
}

func start() (*ffmpeg.Instance, string, error) {
	path, err := utils.NewRecordFilepath(Conf.Path, Conf.Alias)
	if err != nil {
		return nil, "", err
	}
	i, err := ffmpeg.NewInstance(Conf.Url, path+".tmp.mkv", Conf.RecordingDuration, Conf.ClosingTimeout, Conf.Verbose)
	if err != nil {
		return nil, "", fmt.Errorf("error creating an ffmpeg instance: %w", err)
	}
	return i, path, nil
}

func stop(i *ffmpeg.Instance) {
	i.Stop()
	if err := <-i.ExitChan; err != nil {
		log.Errorf("error stopping ffmpeg instance: %s", err)
	}
}

func postRecordOps(recordPath string, reportOk bool) {
	if err := os.Rename(recordPath+".tmp.mkv", recordPath+".mkv"); err != nil {
		log.Errorf("error removing tmp extension from file '%s': %s", recordPath, err.Error())
	}
	if !reportOk {
		return
	}
	http_push.Report(true, "")
}
