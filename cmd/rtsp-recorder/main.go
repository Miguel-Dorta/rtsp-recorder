package main

import (
	"fmt"
	log "github.com/Miguel-Dorta/logolang"
	"github.com/Miguel-Dorta/rtsp-recorder/pkg"
	"github.com/Miguel-Dorta/rtsp-recorder/pkg/mainLoop"
	"github.com/jessevdk/go-flags"
	"os"
	"strings"
	"time"
)

var args struct {
	Alias         string `short:"a" long:"alias" env:"CAMERA_ALIAS" description:"Camera alias"`
	RecTimeMin    int    `short:"t" long:"time" env:"RECORDING_TIME" description:"Duration of each recording (in minutes)" default:"10"`
	RecTimeoutSec int    `long:"timeout" env:"RECORDING_TIMEOUT" description:"Time before killing recording process (in seconds)" default:"60"`
	SavingPath    string `short:"p" long:"path" env:"SAVING_PATH" description:"Path to save the recordings"`
	URL           string `short:"u" long:"url" env:"RTSP_URL" description:"RTSP URL. It must start with rtsp://"`
	Verbose       bool   `short:"v" long:"verbose" env:"VERBOSE" description:"Verbose output"`
	Version       bool   `short:"V" long:"version" description:"Print version and exit"`
}

func init() {
	log.DefaultLogger.Color = false

	if _, err := flags.Parse(&args); err != nil {
		logCriticalf("Error parsing args: %s", err)
	}
	if args.Version {
		fmt.Println(pkg.Version)
		os.Exit(0)
	}
	if args.Verbose {
		log.DefaultLogger.Level = log.LevelDebug
	}

	// Verify args
	if args.Alias == "" {
		logCriticalf("alias cannot be empty")
	}
	if args.SavingPath == "" {
		logCriticalf("path cannot be empty")
	}
	if !strings.HasPrefix(args.URL, "rtsp://") {
		logCriticalf("invalid url")
	}
}

func logCriticalf(format string, v ...interface{}) {
	log.Criticalf(format, v...)
	os.Exit(1)
}

func main() {
	recDuration := time.Minute * time.Duration(args.RecTimeMin)
	timeout := time.Second * time.Duration(args.RecTimeoutSec)
	os.Exit(mainLoop.Start(args.Alias, args.SavingPath, args.URL, recDuration, timeout, args.Verbose))
}
