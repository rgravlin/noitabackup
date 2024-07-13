package internal

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

type Job struct {
	src string
	dst string
}

func worker(jobs chan Job, errs chan error, workersGroup *sync.WaitGroup) {
	defer workersGroup.Done()

	for job := range jobs {
		srcInfo, err := os.Stat(job.src)
		if err != nil {
			errs <- fmt.Errorf("%s: %w", ErrStatFile, err)
			return
		}
		switch srcInfo.Mode() & os.ModeType {
		case os.ModeDir:
			continue
		default:
			if err := copyFile(job.src, job.dst); err != nil {
				errs <- fmt.Errorf("%s: %v", ErrCopyFile, err)
				return
			}
		}
	}
}

func concurrentCopy(src, dst string, dirCounter, fileCounter *int, numOfWorkers int) error {
	jobs := make(chan Job)
	errs := make(chan error)

	go hydrateChannel(jobs, errs, src, dst, dirCounter, fileCounter)

	var workersGroup sync.WaitGroup
	for i := 0; i < numOfWorkers; i++ {
		workersGroup.Add(1)
		go worker(jobs, errs, &workersGroup)
	}
	workersGroup.Wait()

	close(errs)
	errors := make([]string, 0)
	for err := range errs {
		errors = append(errors, err.Error())
	}
	if len(errors) > 0 {
		return fmt.Errorf("%d worker errors occurred: %s", len(errors), strings.Join(errors, ","))
	}

	return nil
}

func hydrateChannel(jobs chan Job, errs chan error, src, dst string, dirCounter, fileCounter *int) {
	defer close(jobs)
	if err := buildDirectory(jobs, src, dst, dirCounter, fileCounter); err != nil {
		errs <- err
		return
	}
}
