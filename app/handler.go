package main

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var builtinCommands = map[string]bool{
	"echo":     true,
	"exit":     true,
	"type":     true,
	"pwd":      true,
	"cd":       true,
	"complete": true,
	"jobs":     true,
	"history":  true,
	"declare":  true,
}

var completionScript = map[string]string{}

func Handler(command string, args []string, outStream *os.File, errStream *os.File, isBg bool) (string, error) {
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
	case "jobs":
		return handleJobs(args)
	case "history":
		return handleHistory(args)
	case "declare":
		return handleDeclare(args)
	default:
		return handleExternal(command, args, outStream, errStream, isBg)
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

func handleExternal(command string, args []string, outStream *os.File, errStream *os.File, isBg bool) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = outStream
	cmd.Stderr = errStream

	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("%v: command not found", command)
	}

	if isBg {
		jobID := AddJob(cmd.Process.Pid, command, args)
		fmt.Printf("[%v] %d\n", jobID, cmd.Process.Pid)

		go func(id int, c *exec.Cmd) {
			c.Wait()
			MarkJobDone(id)
		}(jobID, cmd)

		return "", nil
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
		return "", fmt.Errorf("Invalid command")
	}
}

func handleJobs(args []string) (string, error) {
	return PrintJobs(""), nil
}

func runPipeline(cmds [][]string, finalOut *os.File, finalErr *os.File, isBg bool) error {
	var waitFuncs []func() error
	var nextStdin *os.File = os.Stdin
	var lastPid int

	for i, c := range cmds {
		if len(c) == 0 {
			continue
		}
		cmdName := c[0]
		args := c[1:]
		isLast := (i == len(cmds)-1)

		var outStream *os.File
		var pipeReader *os.File
		if isLast {
			outStream = finalOut
		} else {
			var err error
			pipeReader, outStream, err = os.Pipe()
			if err != nil {
				return err
			}
		}

		currentStdin := nextStdin

		if builtinCommands[cmdName] {
			errCh := make(chan error, 1)
			go func(name string, args []string, in *os.File, out *os.File, isL bool) {
				if in != os.Stdin {
					defer in.Close()
				}
				if !isL {
					defer out.Close()
				}
				outStr, err := Handler(name, args, out, finalErr, false)
				if err != nil {
					fmt.Fprintf(finalErr, "%s\n", err.Error())
				} else if outStr != "" {
					fmt.Fprintf(out, "%s\n", outStr)
				}
				errCh <- err
			}(cmdName, args, currentStdin, outStream, isLast)

			waitFuncs = append(waitFuncs, func() error { return <-errCh })
		} else {
			cmd := exec.Command(cmdName, args...)
			cmd.Stdin = currentStdin
			cmd.Stdout = outStream
			cmd.Stderr = finalErr

			err := cmd.Start()
			if err != nil {
				fmt.Fprintf(finalErr, "%v: command not found\n", cmdName)
			} else {
				lastPid = cmd.Process.Pid
			}

			if !isLast {
				outStream.Close()
			}
			if currentStdin != os.Stdin {
				currentStdin.Close()
			}

			if err != nil {
				waitFuncs = append(waitFuncs, func() error { return err })
			} else {
				waitFuncs = append(waitFuncs, func() error { return cmd.Wait() })
			}
		}

		if !isLast {
			nextStdin = pipeReader
		}
	}

	if isBg {
		var fullCmdStr []string
		for _, c := range cmds {
			fullCmdStr = append(fullCmdStr, strings.Join(c, " "))
		}
		cmdStr := strings.Join(fullCmdStr, " | ")
		jobID := AddJob(lastPid, cmdStr, nil)
		fmt.Printf("[%v] %d\n", jobID, lastPid)

		go func() {
			for _, wf := range waitFuncs {
				wf()
			}
			MarkJobDone(jobID)
		}()
		return nil
	}

	for _, wf := range waitFuncs {
		wf()
	}

	return nil
}

func handleHistory(args []string) (string, error) {
	if len(args) > 0 {
		flag := args[0]
		switch flag {
		case "-r":
			if len(args) > 1 {
				err := AppendHistory(args[1])
				return "", err
			}
			return "", nil
		case "-w":
			if len(args) > 1 {
				err := WriteHistory(args[1], false)
				return "", err
			}
			return "", nil
		case "-a":
			if len(args) > 1 {
				err := WriteHistory(args[1], true)
				return "", err
			}
			return "", nil
		}
	}
	cnt := 1000
	if len(args) != 0 {
		var err error
		cnt, err = strconv.Atoi(args[0])
		if err != nil {
			return "", fmt.Errorf("invalid number: %v", args[0])
		}
	}
	var res strings.Builder
	for i := int(math.Max(0, float64(len(history)-cnt))); i < len(history) && cnt >= 0; i++ {
		fmt.Fprintf(&res, "\t %v %v\n", i+1, strings.Join(history[i], " "))
		cnt--
	}
	return strings.TrimSuffix(res.String(), "\n"), nil
}

func handleDeclare(args []string) (string, error) {
	x := args[0]
	if x == "-p" {
		k := args[1]
		if val, exists := variableStore[k]; exists {
			return fmt.Sprintf("declare -- %v=\"%v\"", k, val), nil
		}
		return "", fmt.Errorf("declare: %v: not found", k)
	} else {
		kvpair := strings.Split(x, "=")
		if len(kvpair) == 1 {
			return "", fmt.Errorf("declare: Invalid format")
		}
		k := kvpair[0]
		v := kvpair[1]
		_, err := strconv.Atoi(string(k[0]))
		if err != nil  {
			return "", fmt.Errorf("declare: `%v=%v': not a valid identifier", k, v)
		}
		variableStore[k] = v
	}
	return "", nil
}
