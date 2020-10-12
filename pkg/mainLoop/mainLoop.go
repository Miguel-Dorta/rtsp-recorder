package mainLoop

import (
	log "github.com/Miguel-Dorta/logolang"
	"github.com/Miguel-Dorta/rtsp-recorder/pkg/ffmpeg"
	"github.com/Miguel-Dorta/rtsp-recorder/pkg/utils"
	"golang.org/x/sys/unix"
	"os"
	"os/signal"
	"time"
)

// Start does the main logic of the application
func Start(alias, savingPath, url string, recDuration, timeout time.Duration, verbose bool) int {
	quit := make(chan os.Signal, 2)
	signal.Notify(quit, unix.SIGINT, unix.SIGTERM)

	for {
		path, err := utils.NewRecordFilepath(savingPath, alias)
		if err != nil {
			log.Critical(err.Error())
			return 1
		}
		i, err := ffmpeg.NewInstance(url, path, recDuration, timeout, verbose)
		if err != nil {
			log.Criticalf("error creating an ffmpeg instance: %s", err)
			return 1
		}

		select {
		case err := <- i.ExitChan:
			if err != nil {
				log.Criticalf("ffmpeg instance returned an error: %s", err)
				return 1
			}
		case <-utils.NewChannelWithTimeout(recDuration):
			stop(i)
		case <-quit:
			stop(i)
			return 0
		}
	}
}

func stop(i *ffmpeg.Instance) {
	i.Stop()
	if err := <- i.ExitChan; err != nil {
		log.Errorf("error stopping ffmpeg instance: %s", err)
	}
}
