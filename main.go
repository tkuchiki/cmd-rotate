package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"os"
	"strings"
	"sync"
)

var (
	stdoutFile = kingpin.Flag("stdout-log", "stdout log name").Default("stdout.log").String()
	stderrFile = kingpin.Flag("stderr-log", "stderr log name").Default("stderr.log").String()
	mergeLog   = kingpin.Flag("merge-log", "stdout and stderr write to same file").Bool()
	logDir     = kingpin.Flag("logdir", "log directory location").PlaceHolder("$TMPDIR").Default(os.TempDir()).String()
	fileMode   = kingpin.Flag("file-mode", "file permission").Default("0644").String()
	fileSize   = kingpin.Flag("file-size", "rotate file size").Default("10485760").Int()
	fileNum    = kingpin.Flag("file-num", "number of files").Default("20").Int()
	args       = kingpin.Arg("args", "command").Strings()
)

func main() {
	kingpin.Version("0.1.0")
	kingpin.Parse()

	if *fileNum < 1 {
		log.Fatal("Invalid --file-num")
	}

	mu := new(sync.Mutex)
	fout := file{
		dir:  *logDir,
		name: *stdoutFile,
		mode: *fileMode,
		num:  int64(*fileNum),
		size: int64(*fileSize),
		mu:   mu,
	}

	fpOut, fpOutErr := fout.openFile()
	if fpOutErr != nil {
		log.Fatal(fpOutErr)
	}
	defer fout.Close()
	fout.fp = fpOut

	c := command{
		args: strings.Join(*args, " "),
	}

	if !*mergeLog {
		ferr := file{
			dir:  *logDir,
			name: *stderrFile,
			mode: *fileMode,
			num:  int64(*fileNum),
			size: int64(*fileSize),
			mu:   mu,
		}

		fpErr, fpErrErr := ferr.openFile()
		if fpErrErr != nil {
			log.Fatal(fpErrErr)
		}
		defer ferr.Close()
		ferr.fp = fpErr

		c.stderrLog = ferr
	}

	c.stdoutLog = fout

	_, cerr := c.runCommand()

	if cerr != nil {
		log.Fatal(cerr)
	}
}
