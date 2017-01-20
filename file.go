package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type file struct {
	fp   *os.File
	name string
	dir  string
	mode string
	num  int64
	size int64
	mu   *sync.Mutex
}

type fileInfos []os.FileInfo

func (fi fileInfos) Len() int {
	return len(fi)
}

func (fi fileInfos) Less(i, j int) bool {
	f1 := strings.Split(fi[i].Name(), "-")
	f2 := strings.Split(fi[j].Name(), "-")

	f1nano, _ := strconv.ParseInt(f1[len(f1)-1], 10, 64)
	f2nano, _ := strconv.ParseInt(f2[len(f2)-1], 10, 64)
	return f1nano < f2nano
}

func (fi fileInfos) Swap(i, j int) {
	fi[i], fi[j] = fi[j], fi[i]
}

func stringToFileMode(mode string) (os.FileMode, error) {
	m, err := strconv.ParseInt(mode, 8, 0)
	if err != nil {
		var fmode os.FileMode
		return fmode, err
	}

	return os.FileMode(m), nil
}

func sortFileInfos(fi fileInfos, reverse bool) fileInfos {
	if reverse {
		sort.Sort(sort.Reverse(fileInfos(fi)))
	} else {
		sort.Sort(fileInfos(fi))
	}

	return fi
}

func (f *file) absPath() string {
	return filepath.Join(f.dir, f.name)
}

func (f *file) openFile() (*os.File, error) {
	var fp *os.File
	m, err := stringToFileMode(f.mode)
	if err != nil {
		return fp, err
	}

	fp, err = os.OpenFile(f.absPath(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.FileMode(m))
	return fp, err
}

func (f *file) writeLine(str string) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.fp.WriteString(fmt.Sprintln(str))
}

func (f *file) rename(old, new string) error {
	return os.Rename(old, new)
}

func (f *file) rotate() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if err := f.close(); err != nil {
		return err
	}

	old := f.absPath()
	new := fmt.Sprintf("%s-%d", old, time.Now().UnixNano())
	if err := f.rename(old, new); err != nil {
		return err
	}

	fp, err := f.openFile()
	if err != nil {
		return err
	}

	f.fp = fp
	err = f.cleanup()

	return err
}

func (f *file) checkFileSize() (bool, error) {
	stat, err := f.fp.Stat()
	if err != nil {
		return false, err
	}

	return stat.Size() > f.size, nil
}

func (f *file) listFiles() ([]os.FileInfo, error) {
	infos := make([]os.FileInfo, 0)
	allInfos, err := ioutil.ReadDir(f.dir)

	for _, fi := range allInfos {
		if fi.IsDir() {
			continue
		}

		filename := fi.Name()
		if strings.Index(filename, fmt.Sprintf("%s-", f.name)) == 0 {
			infos = append(infos, fi)
		}
	}

	sortFileInfos(infos, true)

	return infos, err
}

func (f *file) remove(rmFile string) error {
	return os.Remove(rmFile)
}

func (f *file) cleanup() error {
	infos, err := f.listFiles()
	if err != nil {
		return err
	}

	if int64(len(infos)) < f.num {
		return nil
	}

	for _, fi := range infos[f.num-1:] {
		filename := filepath.Join(f.dir, fi.Name())
		if f.remove(filename) != nil {
			return err
		} else {
			fmt.Println("removed", filename)
		}
	}

	return nil
}

func (f *file) close() error {
	err := f.fp.Close()
	f.fp = nil
	return err
}

func (f *file) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.fp.Close()
}
