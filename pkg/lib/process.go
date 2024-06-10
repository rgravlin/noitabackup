package lib

import (
	"golang.org/x/sys/windows"
	"os"
	"syscall"
)

const (
	processEntrySize = 568
	noitaProcessName = "noita.exe"
)

// processID takes a name parameter and returns the process ID of the process with that name. If no process is found,
// it returns an error.
//
// The function uses the Windows API function CreateToolhelp32Snapshot to take a snapshot of the current system information,
// including the list of running processes.
//
// It then iterates over each process entry in the snapshot and checks if the name of the executable matches the given name.
// If a match is found, it returns the process ID.
//
// If no process is found with the given name, it returns an error.
//
// Example usage:
//
//	pid, err := processID(noitaProcessName)
//	if err != nil {
//	    // handle error
//	}
//	// use pid
//
// Note: The constant processEntrySize needs to be declared before using this function. It specifies the size of the
// windows.ProcessEntry32 struct used by the Windows API.
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

// isNoitaRunning checks if the Noita process is currently running.
// It first retrieves the process ID of the Noita process using the processID function.
// Next, it uses os.FindProcess to find the process with the given process ID.
// If the process is found, it tries to send a Signal(0) to the process, which is a safe way to check if the process is running.
// If an error occurs during the Signal(0) call, we can assume that the process is running,
// and thus, the function returns true. Otherwise, it returns false.
// Note that the noitaProcessName constant needs to be declared and defined before using this function.
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
