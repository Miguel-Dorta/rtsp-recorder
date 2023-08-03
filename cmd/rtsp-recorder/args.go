package main

import (
	"fmt"
	"github.com/Miguel-Dorta/rtsp-recorder/pkg"
	"github.com/Miguel-Dorta/rtsp-recorder/pkg/http_push"
	"github.com/Miguel-Dorta/rtsp-recorder/pkg/log"
	"github.com/jessevdk/go-flags"
	"os"
	"strings"
	"time"
)

var args struct {
	Alias            string `short:"a" long:"alias" env:"RECORDING_ALIAS" description:"Alias of the recording"`
	Proto            string `long:"protocol" env:"RTSP_PROTOCOL" description:"Protocol to use for RTSP" default:"TCP"`
	RecTimeMin       int    `short:"t" long:"time" env:"RECORDING_TIME" description:"Duration of each recording (in minutes)" default:"10"`
	RecTimeoutSec    int    `long:"timeout" env:"RECORDING_TIMEOUT" description:"Time before killing recording process (in seconds)" default:"60"`
	ReportMethod     string `long:"report-method" env:"REPORT_METHOD" description:"Method to do HTTP reports" default:"POST"`
	ReportTimeoutSec int    `long:"report-timeout" env:"REPORT_TIMEOUT" description:"Timeout for HTTP report requests (in seconds)" default:"5"`
	ReportUrlErr     string `long:"report-url-error" env:"REPORT_URL_ERROR" description:"Do HTTP request to this URL to report errors"`
	ReportUrlOK      string `long:"report-url-ok" env:"REPORT_URL_OK" description:"Do HTTP request to this URL to confirm successful recordings"`
	SavingPath       string `short:"p" long:"path" env:"SAVING_PATH" description:"Path to save the recordings"`
	URL              string `short:"u" long:"url" env:"RTSP_URL" description:"RTSP URL. It must start with rtsp://"`
	Verbose          bool   `short:"v" long:"verbose" env:"VERBOSE" description:"Verbose output"`
	Version          bool   `short:"V" long:"version" description:"Print version and exit"`
}

func ParseArgs() {
	if _, err := flags.Parse(&args); err != nil {
		if !flags.WroteHelp(err) {
			log.Criticalf("error parsing args: %s", err)
		}
		os.Exit(0)
	}
	if args.Version {
		fmt.Println(pkg.Version)
		os.Exit(0)
	}

	if args.Alias == "" {
		log.Critical("alias cannot be empty")
	}
	if args.Proto != "TCP" && args.Proto != "UDP" {
		log.Critical("protocol must be 'TCP' or 'UDP'")
	}
	args.Proto = strings.ToLower(args.Proto)
	if args.RecTimeMin < 1 {
		log.Critical("recording time cannot be less than a minute")
	}
	if args.RecTimeoutSec < 15 {
		log.Critical("stop recording timeout cannot be less than 15 seconds")
	}
	if args.ReportUrlOK != "" && args.ReportUrlErr != "" {
		http_push.InitReports(args.ReportUrlOK, args.ReportUrlErr, args.ReportMethod, time.Second*time.Duration(args.ReportTimeoutSec))
	}
	if args.SavingPath == "" {
		log.Critical("path cannot be empty")
	}
	if !strings.HasPrefix(args.URL, "rtsp://") {
		log.Critical("invalid url")
	}
}
