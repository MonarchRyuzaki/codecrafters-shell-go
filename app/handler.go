package main

import (
	"strings"
	"fmt"
)

func Handler(command string, args []string) string {
	switch command {
	case "echo":
		return handleEcho(args)
	default:
		return fmt.Sprintf("%v: command not found", command)
	}
}

func handleEcho(args []string) string {
	return strings.Join(args, " ")
}