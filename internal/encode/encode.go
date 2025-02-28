package encode

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"os/exec"

	"golang.org/x/sync/errgroup"
)

// Remember to install ffmpeg
func ffmpeg(w io.WriteCloser, path string) error {
	ffmpegArgs := []string{
		"-i", path,
		"-ar", "48000",
		"-ac", "2",
		"-af", "volume=0.5",
		"-f", "s16le",
		"pipe:1",
	}
	ffmpeg := exec.Command("ffmpeg", ffmpegArgs...)
	ffmpeg.Stdout = w
	defer w.Close()

	err := ffmpeg.Start()
	if err != nil {
		return err
	}

	err = ffmpeg.Wait()
	if err != nil {
		return err
	}

	if ffmpeg.Err != nil {
		return ffmpeg.Err
	}

	return nil
}

// Remember to install dca
func dca(r io.Reader, w io.WriteCloser) error {
	dca := exec.Command("dca")
	dca.Stdin = r
	dca.Stdout = w
	dca.Stderr = os.Stderr
	defer w.Close()

	err := dca.Start()
	if err != nil {
		return err
	}

	err = dca.Wait()
	if err != nil {
		return err
	}

	return nil
}

func read(r io.Reader, output chan<- []byte) error {
	var opuslen int16
	defer close(output)
	for {
		err := binary.Read(r, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return fmt.Errorf("eof")
		} else if err != nil {
			return fmt.Errorf("error reading from dca file: %v", err)
		}

		// Read encoded pcm from dca file.
		buf := make([]byte, opuslen)
		binary.Read(r, binary.LittleEndian, &buf)

		output <- buf
	}
}

func Encode(path string, outputChan chan []byte) error {
	dr, fw := io.Pipe()
	or, dw := io.Pipe()

	errs, ctx := errgroup.WithContext(context.Background())
	defer ctx.Done()

	errs.Go(func() error { return ffmpeg(fw, path) })
	errs.Go(func() error { return dca(dr, dw) })
	errs.Go(func() error { return read(or, outputChan) })

	err := errs.Wait()

	if err != nil {
		return err
	}

	return nil
}
