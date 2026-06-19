package main

import (
	"fmt"
	"strings"
	"sync"
)

type Job struct {
	ID      int
	PID     int
	Command string
	Status  string
}

var jobs []Job
var bgCounter int = 1
var jobsMutex sync.Mutex 

func AddJob(pid int, command string, args []string) int {
	jobsMutex.Lock()
	defer jobsMutex.Unlock()

	fullCmd := command
	if len(args) > 0 {
		fullCmd += " " + strings.Join(args, " ")
	}

	job := Job{
		ID:      bgCounter,
		PID:     pid,
		Command: fullCmd,
		Status:  "Running",
	}
	jobs = append(jobs, job)

	bgCounter++
	return job.ID
}

func MarkJobDone(id int) {
	jobsMutex.Lock()
	defer jobsMutex.Unlock()

	for i := range jobs {
		if jobs[i].ID == id {
			jobs[i].Status = "Done"
			break
		}
	}
}

func PrintJobs() string {
	jobsMutex.Lock()
	defer jobsMutex.Unlock()

	var sb strings.Builder
	var remainingJobs []Job

	for i, job := range jobs {
		marker := " "
		if i == len(jobs)-1 {
			marker = "+"
		} else if i == len(jobs)-2 {
			marker = "-"
		}

		sb.WriteString(fmt.Sprintf("[%d]%s  %-24s%s\n", job.ID, marker, job.Status, job.Command))

		if job.Status != "Done" {
			remainingJobs = append(remainingJobs, job)
		}
	}

	jobs = remainingJobs

	return strings.TrimRight(sb.String(), "\n")
}
