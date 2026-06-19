package main

import (
	"fmt"
	"strings"
)

type Job struct {
	ID      int
	PID     int
	Command string
	Status  string
}

var jobs []Job
var bgCounter int = 1

func AddJob(pid int, command string, args []string) int {
	fullCmd := command
	if len(args) > 0 {
		fullCmd += " " + strings.Join(args, " ")
	}
	fullCmd += " &"

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

func PrintJobs() string {
	var sb strings.Builder
	for i, job := range jobs {
		marker := " "
		if i == len(jobs)-1 {
			marker = "+"
		}
		if i == len(jobs) - 2 {
			marker = "-"
		}

		fmt.Fprintf(&sb, "[%d]%s  %-24s%s\n", job.ID, marker, job.Status, job.Command)
	}

	return strings.TrimRight(sb.String(), "\n")
}
