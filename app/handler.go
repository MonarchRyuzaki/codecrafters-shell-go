package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var builtinCommands = map[string]bool{
	"echo":     true,
	"exit":     true,
	"type":     true,
	"pwd":      true,
	"cd":       true,
	"complete": true,
}

var completionScript = map[string]string{}

func Handler(command string, args []string, outStream *os.File, errStream *os.File) (string, error) {
	switch command {
	case "echo":
		return handleEcho(args)
	case "type":
		return handleType(args)
	case "pwd":
		return handlePwd(args)
	case "cd":
		return handleCd(args)
	case "complete":
		return handleComplete(args)
	default:
		return handleExternal(command, args, outStream, errStream)
	}
}

func handleEcho(args []string) (string, error) {
	return strings.Join(args, " "), nil
}

func handleType(args []string) (string, error) {
	cmd := args[0]
	_, ok := builtinCommands[cmd]
	if ok {
		return fmt.Sprintf("%v is a shell builtin", cmd), nil
	}
	path, err := exec.LookPath(cmd)
	if err == nil {
		return fmt.Sprintf("%v is %v", cmd, path), nil
	}
	return "", fmt.Errorf("%v: not found", cmd)
}

func handleExternal(command string, args []string, outStream *os.File, errStream *os.File) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = outStream
	cmd.Stderr = errStream

	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("%v: command not found", command)
	}

	if err := cmd.Wait(); err != nil {
		return "", nil
	}

	return "", nil
}

func handlePwd(args []string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return dir, nil
}

func handleCd(args []string) (string, error) {
	newDir := args[0]
	if newDir == "~" {
		d, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("home directory not defined")
		}
		newDir = d
	}
	_, err := os.Stat(newDir)
	if err != nil {
		return "", fmt.Errorf("cd: %v: No such file or directory", newDir)
	}
	err = os.Chdir(newDir)
	if err != nil {
		return "", fmt.Errorf("cd: %v: No such file or directory", newDir)
	}
	return "", nil
}

func handleComplete(args []string) (string, error) {
	switch args[0] {
	case "-p", "-P":
		prog := args[1]
		if path, exists := completionScript[prog]; exists {
			return fmt.Sprintf("complete -C '%v' %v", path, prog), nil
		}
		return "", fmt.Errorf("complete: %v: no completion specification", prog)
	case "-C", "-c":
		path := args[1]
		prog := args[2]
		completionScript[prog] = path
		return "", nil
	case "-r", "-R":
		prog := args[1]
		delete(completionScript, prog)
		return "", nil
	default:
		return "", fmt.Errorf("Invalid command");
	}
}