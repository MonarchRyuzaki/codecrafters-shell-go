package main

import (
	"strings"
	"fmt"
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
	_, ok := builtinCommands[args[0]]
	if ok {
		return fmt.Sprintf("%v is a shell builtin", args[0])
	}
	return fmt.Sprintf("%v: not found", args[0])
}