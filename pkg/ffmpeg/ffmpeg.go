package ffmpeg

import (
	"context"
	"fmt"
	"github.com/Miguel-Dorta/rtsp-recorder/pkg/log"
	"golang.org/x/sys/unix"
	"io"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type Instance struct {
	ExitChannel chan error
	cmd         *exec.Cmd
	stdin       io.WriteCloser
	isStopped   bool
}

var ffmpegPath string

func init() {
	var err error
	ffmpegPath, err = exec.LookPath("ffmpeg")
	if err != nil {
		log.Critical("dependency ffmpeg not found")
	}
}

func StartRecording(rtspUrl, rtspProto, recordPath string, recDuration, timeout time.Duration, verbose bool) (*Instance, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), recDuration+timeout)
	cmd := exec.CommandContext(ctx, ffmpegPath,
		"-nostdin",
		"-rtsp_transport", rtspProto,
		"-t", strconv.Itoa(int(recDuration.Seconds())),
		"-i", rtspUrl,
		"-c", "copy",
		recordPath)
	if verbose {
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		cancelCtx()
		return nil, fmt.Errorf("cannot create stding pipe in ffmpeg instance: %w", err)
	}

	instance := &Instance{
		ExitChannel: make(chan error),
		cmd:         cmd,
		stdin:       stdin,
		isStopped:   false,
	}

	go func() {
		instance.ExitChannel <- cmd.Run()
		instance.isStopped = true
		cancelCtx()
	}()
	go func() {
		time.Sleep(recDuration)
		if instance.isStopped {
			return
		}
		instance.Stop()
	}()

	return instance, nil
}

func (instance *Instance) Stop() {
	_, err := instance.stdin.Write([]byte{'q'})
	if err != nil {
		log.Errorf("error writing exit command to ffmpeg instance: %s", err)
		_ = instance.cmd.Process.Kill()
		return
	}

	time.Sleep(time.Second)
	if instance.isStopped {
		return
	}
	if err = instance.stdin.Close(); err != nil {
		log.Errorf("error closing stdin pipe of ffmpeg instance: %s", err)
		_ = instance.cmd.Process.Kill()
		return
	}

	time.Sleep(time.Second)
	for i := 0; i < 2; i++ {
		if instance.isStopped {
			return
		}
		if err = instance.cmd.Process.Signal(unix.SIGINT); err != nil {
			log.Errorf("error sending SIGINT to ffmpeg instance: %s", err)
			_ = instance.cmd.Process.Kill()
			return
		}
		time.Sleep(time.Second * 5)
	}

	if instance.isStopped {
		return
	}
	log.Errorf("instance killed due to not stopping after multiple safe stop mechanisms")
	_ = instance.cmd.Process.Kill()
}
