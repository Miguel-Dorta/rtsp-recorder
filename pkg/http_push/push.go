package http_push

import (
	log "github.com/Miguel-Dorta/logolang"
	"github.com/Miguel-Dorta/rtsp-recorder/pkg/utils"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	c                     *http.Client
	urlOk, urlErr, method string
	enabled               = false
)

func InitReports(reportOkUrl, reportErrUrl, reportMethod string, timeout time.Duration) {
	urlOk = reportOkUrl
	urlErr = reportErrUrl
	method = reportMethod
	c = &http.Client{Timeout: timeout}
	enabled = true
}

func Report(ok bool, msg string) {
	if !enabled {
		return
	}

	req, err := http.NewRequest(method, utils.TernaryOperator(ok, urlOk, urlErr), strings.NewReader(msg))
	if err != nil {
		log.Criticalf("error creating HTTP push request: %s", err)
		os.Exit(1)
	}

	resp, err := c.Do(req)
	if err != nil {
		log.Criticalf("error doing HTTP push request: %s", err)
		os.Exit(1)
	}
	if err = resp.Body.Close(); err != nil {
		log.Errorf("error closing HTTP push response body: %s", err)
		os.Exit(1)
	}
	if resp.StatusCode >= 400 {
		log.Criticalf("error status (%d) returned in HTTP push: %s", resp.StatusCode, err)
		os.Exit(1)
	}
}
