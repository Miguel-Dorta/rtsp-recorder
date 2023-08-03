package main

import (
	"github.com/Miguel-Dorta/rtsp-recorder/pkg/log"
	"github.com/Miguel-Dorta/rtsp-recorder/pkg/recorder"
	"time"
)

func main() {
	ParseArgs()
	quit := StartSignalListener()

	for {
		rec, err := recorder.Start(args.SavingPath, args.Alias, args.URL, args.Proto,
			time.Minute*time.Duration(args.RecTimeMin),
			time.Second*time.Duration(args.RecTimeoutSec),
			args.Verbose)
		if err != nil {
			log.Criticalf("error starting recorder: %s", err)
		}

		if rec.Wait(quit) {
			return
		}
	}
}
