package internal

import (
	"golang.org/x/sys/windows"
	"os"
	"syscall"
)

const (
	processEntrySize = 568
	noitaProcessName = "noita.exe"
)

func processID(name string) (uint32, error) {
	h, e := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if e != nil {
		return 0, e
	}
	p := windows.ProcessEntry32{Size: processEntrySize}
	for {
		e := windows.Process32Next(h, &p)
		if e != nil {
			return 0, e
		}
		if windows.UTF16ToString(p.ExeFile[:]) == name {
			return p.ProcessID, nil
		}
	}
}

func isNoitaRunning() bool {
	pid, err := processID(noitaProcessName)
	if err != nil {
		return false
	}

	process, err := os.FindProcess(int(pid))
	if err != nil {
		return false
	} else {
		err := process.Signal(syscall.Signal(0))
		if err != nil {
			return true
		} else {
			return false
		}
	}
}
