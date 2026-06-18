package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var builtinCommands = map[string]bool{
	"echo" : true,
	"exit" : true,
	"type" : true,
	"pwd"  : true,
}

func Handler(command string, args []string) string {
	switch command {
	case "echo":
		return handleEcho(args)
	case "type":
		return handleType(args)
	case "pwd":
		return handlePwd(args)
	default:
		return handleExternal(command, args)
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

func handleExternal(command string, args []string) string {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Sprintf("%v: command not found", command)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Sprintf("%s", err.Error())
	}

	return "";
}

func handlePwd(args []string) string {
	dir, err := os.Getwd()
	if err != nil {
		return err.Error()
	}
	return dir
}