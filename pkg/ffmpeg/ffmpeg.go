package ffmpeg

import (
	"context"
	"fmt"
	log "github.com/Miguel-Dorta/logolang"
	"io"
	"os"
	"os/exec"
	"time"
)

type Instance struct {
	ExitChan chan error
	cmd      *exec.Cmd
	stdin    io.WriteCloser
}

var ffmpegPath string

func init() {
	var err error
	ffmpegPath, err = exec.LookPath("ffmpeg")
	if err != nil {
		log.DefaultLogger.Color = false
		log.Critical("dependency ffmpeg not found")
		os.Exit(1)
	}
}

// NewInstance starts a new instance of FFMPEG with the args provided
// Args:
//     - url: RTSP URL
//     - path: path of the file to save the recording
//     - recDuration: duration of the recording
//     - timeout: the extra time for finishing and saving. The instance will be killed at recDuration + timeout
func NewInstance(url, path string, recDuration, timeout time.Duration, verbose bool) (*Instance, error) {
	ctx, _ := context.WithTimeout(context.Background(), recDuration+timeout)
	cmd := exec.CommandContext(ctx, ffmpegPath, "-rtsp_transport", "tcp", "-i", url, "-c", "copy", path)
	if verbose {
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("error creating stdin pipe in ffmpeg instance: %w", err)
	}

	exitChannel := make(chan error)
	go func() {
		if err := cmd.Run(); err != nil {
			exitChannel <- err
		}
		close(exitChannel)
	}()

	return &Instance{
		cmd:      cmd,
		stdin:    stdin,
		ExitChan: exitChannel,
	}, nil
}

// Stop tries to stop the instance gratefully. It is NON blocking. If it fails, it kills the instance.
func (i *Instance) Stop() {
	_, err := io.WriteString(i.stdin, "q")
	if err != nil {
		i.cmd.Process.Kill()
	}
}
