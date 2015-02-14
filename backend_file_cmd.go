package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// FileBackendCmd represents a buffer file,
// whose content is the output from an extenal command.
type FileBackendCmd struct {
	FileBackend
	dir    string
	runner *exec.Cmd
	title  *string
}

// if title == nil then will show the command name
func (e *Editor) NewFileBackendCmd(args []string, dir string, viewId int, title *string) (*FileBackendCmd, error) {
	b, err := e.NewFileBackend(e.BufferFile(viewId), viewId)
	if err != nil {
		return nil, err
	}
	fb := &FileBackendCmd{
		FileBackend: *b,
		dir:         dir,
		runner:      exec.Command(args[0], args[1:]...),
		title:       title,
	}
	go fb.start()
	return fb, nil
}

func (f *FileBackendCmd) Refresh() {
	f.stop()
	go f.start()
}

func (f *FileBackendCmd) Close() error {
	f.stop()
	return nil
}

func (f *FileBackendCmd) start() {
	workDir, _ := filepath.Abs(f.dir)
	v := Ed.ViewById(f.viewId)
	v.WorkDir = workDir
	f.runner.Stdout = f.file
	f.runner.Stderr = f.file
	f.runner.Dir = workDir
	if f.title == nil {
		title := strings.Join(f.runner.Args, " ")
		f.title = &title
	}
	v.title = fmt.Sprintf("[RUNNING] %s", *f.title)
	Ed.Render()
	err := f.runner.Run()
	// TODO: autorefresh every n seconds or new output available ?
	Ed.Open(f.srcLoc, v, "")
	v.WorkDir = workDir // open() would have modified this
	if err != nil {
		v.title = fmt.Sprintf("[FAILED] %s", *f.title)
		Ed.SetStatusErr(err.Error())
	} else {
		v.title = *f.title
	}
	Ed.Render()
}

func (f *FileBackendCmd) stop() {
	if f.runner != nil && f.runner.Process != nil {
		f.runner.Process.Release()
		f.runner.Process.Kill()
	}
	f.runner = nil
}
