package main

import (
	"bufio"
	"io"
	"os/exec"
	"syscall"
)

type command struct {
	stdoutLog file
	stderrLog file
	args      string
}

var (
	stdoutCh = make(chan string)
	stderrCh = make(chan string)
	quitCh   = make(chan struct{})
)

func (c *command) runCommand() (int, error) {
	cmd := exec.Command("sh", "-c", c.args)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return 1, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return 1, err
	}

	if err = cmd.Start(); err != nil {
		return 1, err
	}

	if (file{}) == c.stderrLog {
		go c.readIo(stdout, stdoutCh)
		go c.readIo(stderr, stdoutCh)
	} else {
		go c.readIo(stdout, stdoutCh)
		go c.readIo(stderr, stderrCh)
	}

	go c.writeFile()

	err = cmd.Wait()
	quitCh <- struct{}{}

	var exitCode int
	if err != nil {
		if err2, ok := err.(*exec.ExitError); ok {
			if s, ok := err2.Sys().(syscall.WaitStatus); ok {
				err = nil
				exitCode = s.ExitStatus()
			}
		}
	}

	return exitCode, err
}

func (c *command) readIo(r io.Reader, q chan string) {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		q <- scanner.Text()
	}

}

func (c *command) writeFile() {
	var f *file
	var txt string
	for {
		select {
		case outTxt := <-stdoutCh:
			f = &c.stdoutLog
			txt = outTxt
		case errTxt := <-stderrCh:
			f = &c.stderrLog
			txt = errTxt
		case <-quitCh:
			break
		}

		f.writeLine(txt)

		isExceeded, err := f.checkFileSize()
		if err != nil {
			panic(err)
		}

		if isExceeded {
			if ferr := f.rotate(); ferr != nil {
				panic(ferr)
			}
		}
	}
}
