package pkg

import (
	"fmt"
	log "github.com/Miguel-Dorta/logolang"
	"github.com/Miguel-Dorta/rtsp-recorder/pkg/ffmpeg"
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
func Main() int {
	quit := make(chan os.Signal, 2)
	signal.Notify(quit, unix.SIGINT, unix.SIGTERM)

	i, err := start()
	if err != nil {
		log.Criticalf(err.Error())
		return 1
	}

	for {
		select {
		// Case: ffmpeg ends before the timeout
		case err := <-i.ExitChan:
			if err != nil {
				log.Criticalf("ffmpeg instance returned an error: %s", err)
				return 1
			}
			i, err = start()
			if err != nil {
				log.Critical(err.Error())
				return 1
			}
		// Case: end of the time of the recording
		case <-utils.NewChannelWithTimeout(Conf.RecordingDuration):
			// Start new recording before closing new one to have overlapping videos
			i2, err := start()
			if err != nil {
				log.Critical(err.Error())
				stop(i)
				return 1
			}
			stop(i)
			i = i2
		// Case: User or system quit
		case <-quit:
			stop(i)
			return 0
		}
	}
}

func start() (*ffmpeg.Instance, error) {
	path, err := utils.NewRecordFilepath(Conf.Path, Conf.Alias)
	if err != nil {
		return nil, err
	}
	i, err := ffmpeg.NewInstance(Conf.Url, path, Conf.RecordingDuration, Conf.ClosingTimeout, Conf.Verbose)
	if err != nil {
		return nil, fmt.Errorf("error creating an ffmpeg instance: %w", err)
	}
	return i, nil
}

func stop(i *ffmpeg.Instance) {
	i.Stop()
	if err := <-i.ExitChan; err != nil {
		log.Errorf("error stopping ffmpeg instance: %s", err)
	}
}
