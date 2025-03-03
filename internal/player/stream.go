package player

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"os/exec"

	"golang.org/x/sync/errgroup"
)

const (
	SAMPLE_RATE = "48000"
	CHANNELS    = "2"
)

type Controller struct {
	isPlaying bool
	fcmd      *exec.Cmd
	dcmd      *exec.Cmd
}

func (c *Controller) ffmpeg(src string, dst io.WriteCloser) error {
	ffmpegArgs := []string{
		"-i", src,
		"-ar", SAMPLE_RATE,
		"-ac", CHANNELS,
		"-af", "volume=0.5",
		"-f", "s16le",
		"pipe:1",
	}
	c.fcmd = exec.Command("ffmpeg", ffmpegArgs...)
	c.fcmd.Stdout = dst
	defer dst.Close()

	err := c.fcmd.Start()
	if err != nil {
		return err
	}

	err = c.fcmd.Wait()
	if err != nil {
		return err
	}

	if c.fcmd.Err != nil {
		return c.fcmd.Err
	}

	return nil
}

func (c *Controller) dca(src io.Reader, dst io.WriteCloser) error {
	dcaArgs := []string{}
	c.dcmd = exec.Command("dca", dcaArgs...)
	c.dcmd.Stdin = src
	c.dcmd.Stdout = dst
	c.dcmd.Stderr = os.Stderr
	defer dst.Close()

	err := c.dcmd.Start()
	if err != nil {
		return err
	}

	err = c.dcmd.Wait()
	if err != nil {
		return err
	}

	if c.dcmd.Err != nil {
		return c.dcmd.Err
	}

	return nil
}

func (c *Controller) transfer(src io.Reader, output chan<- []byte) error {
	var opuslen int16

	for {
		if c.isPlaying {
			err := binary.Read(src, binary.LittleEndian, &opuslen)

			// If this is the end of the file, just return.
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				return fmt.Errorf("eof")
			} else if err != nil {
				return fmt.Errorf("error reading from dca file: %v", err)
			}

			// Read encoded pcm from dca file.
			buf := make([]byte, opuslen)
			binary.Read(src, binary.LittleEndian, &buf)

			output <- buf
		}
	}
}

func (c *Controller) Stream(source string, output chan<- []byte, next chan<- bool) error {
	dr, fw := io.Pipe()
	br, dw := io.Pipe()

	errs, ctx := errgroup.WithContext(context.Background())
	defer ctx.Done()

	errs.Go(func() error { return c.ffmpeg(source, fw) })
	errs.Go(func() error { return c.dca(dr, dw) })
	errs.Go(func() error { return c.transfer(br, output) })

	c.isPlaying = true

	err := errs.Wait()
	next <- true

	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) StopStream() error {
	errs, ctx := errgroup.WithContext(context.Background())
	defer ctx.Done()

	errs.Go(func() error { return c.fcmd.Process.Kill() })
	errs.Go(func() error { return c.fcmd.Process.Kill() })

	err := errs.Wait()
	c.isPlaying = false

	if err != nil {
		return err
	}

	return nil
}
