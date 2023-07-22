package http_push

import (
	log "github.com/Miguel-Dorta/logolang"
	"github.com/Miguel-Dorta/rtsp-recorder/pkg/utils"
	"net/http"
	"os"
	"strings"
)

var (
	Url string
	c   = new(http.Client)
)

func Report(ok bool, msg string) {
	resp, err := c.Post(utils.TernaryOperator(ok, Url, Url+"/fail"), "text/plain", strings.NewReader(msg))
	if err != nil {
		log.Criticalf("error sending HTTP push: %s", err.Error())
		os.Exit(1)
	}
	if err = resp.Body.Close(); err != nil {
		log.Errorf("error closing HTTP push response body: %s", err.Error())
		os.Exit(1)
	}
	if resp.StatusCode >= 400 {
		log.Criticalf("error status (%d) returned in HTTP push: %s", resp.StatusCode, err.Error())
		os.Exit(1)
	}
}
