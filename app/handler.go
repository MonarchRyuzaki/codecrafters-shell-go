package main

import (
	"strings"
	"fmt"
	"os/exec"
)

var builtinCommands = map[string]bool{
	"echo" : true,
	"exit" : true,
	"type" : true,
}

func Handler(command string, args []string) string {
	switch command {
	case "echo":
		return handleEcho(args)
	case "type":
		return handleType(args)
	default:
		return fmt.Sprintf("%v: command not found", command)
	}
}

func handleEcho(args []string) string {
	return strings.Join(args, " ")
}

func handleType(args []string) string {
	cmd := args[0]
	_, ok := builtinCommands[cmd]
	if ok {
		return fmt.Sprintf("%v is a shell builtin", cmd)
	}
	path, err := exec.LookPath(cmd)
	if err == nil {
		return fmt.Sprintf("%v is %v", cmd, path)
	}
	return fmt.Sprintf("%v: not found", cmd)
}