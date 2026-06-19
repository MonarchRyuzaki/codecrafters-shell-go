package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/term"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

var autocompleteSet map[string]bool

func init() {
	autocompleteSet = make(map[string]bool)
	for k, v := range builtinCommands {
		autocompleteSet[k] = v
	}
	pathEnv := os.Getenv("PATH")
	paths := filepath.SplitList(pathEnv)

	for _, dir := range paths {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			name := entry.Name()
			if !entry.IsDir() {
				info, err := entry.Info()
				if err == nil && info.Mode()&0111 != 0 {
					autocompleteSet[name] = true
				}
			}
		}
	}
}

// ReadCommand puts the terminal in raw mode and reads keystrokes individually
func readCommand() (string, error) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	var command []byte
	buf := make([]byte, 1)
	lastTabPress := false

	for {
		_, err := os.Stdin.Read(buf)
		if err != nil {
			return string(command), err
		}

		switch buf[0] {
		case '\r', '\n': // Enter
			fmt.Print("\r\n")
			return string(command), nil
		case '\t': // Tab
			typedStr := string(command)
			var matches []string

			for k := range autocompleteSet {
				if strings.HasPrefix(k, typedStr) {
					matches = append(matches, k)
				}
			}

			if len(matches) == 1 {
				matchedBuiltin := matches[0]

				completion := matchedBuiltin[len(typedStr):] + " "

				command = append(command, []byte(completion)...)

				fmt.Print(completion)
				lastTabPress = false 
			} else if len(matches) > 1 {
				lcp := matches[0]
				for _, match := range matches[1:] {
					i := 0
					for i < len(lcp) && i < len(match) && lcp[i] == match[i] {
						i++
					}
					lcp = lcp[:i]
				}

				if len(lcp) > len(typedStr) {
					completion := lcp[len(typedStr):]
					command = append(command, []byte(completion)...)
					fmt.Print(completion)
					
					fmt.Print("\a") 
					lastTabPress = false
				} else {
					if lastTabPress == false {
						fmt.Print("\a")
						lastTabPress = true
					} else {
						fmt.Print("\r\n")
						sort.Strings(matches)
						
						fmt.Print(strings.Join(matches, "  ") + "\r\n")
						fmt.Print("$ " + string(command))
						lastTabPress = false
					}
				}
			} else {
				fmt.Print("\a")
				lastTabPress = false
			}

		case '\x03': // Ctrl + C
			fmt.Print("^C\r\n")
			term.Restore(int(os.Stdin.Fd()), oldState)
			os.Exit(0)
		case '\x7f', '\b': // Backspace
			lastTabPress = false 
			if len(command) > 0 {
				command = command[:len(command)-1]
				fmt.Print("\b \b")
			}
		default: // Normal character
			lastTabPress = false 
			command = append(command, buf[0])
			fmt.Print(string(buf[0]))
		}

	}
}

func main() {
	for {
		fmt.Print("$ ")
		command, err := readCommand()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}
		result := parseInput(command)
		if len(result) == 0 {
			continue
		}
		outStream := os.Stdout
		errStream := os.Stderr
		index := len(result)
		for i := 0; i < len(result); i++ {
			if config, exists := redirectionMap[result[i]]; exists {
				flags := os.O_CREATE | os.O_WRONLY
				if config.Append {
					flags |= os.O_APPEND
				} else {
					flags |= os.O_TRUNC
				}

				file, _ := os.OpenFile(result[i+1], flags, 0644)

				if !config.Stdout {
					outStream = file
				}
				if !config.Stderr {
					errStream = file
				}
				index = i
				break
			}
		}
		command = result[0]
		args := result[1:index]
		// fmt.Println(command, args)
		if command == "exit" {
			break
		}
		out, err := Handler(command, args, outStream, errStream)
		if err != nil {
			fmt.Fprintf(errStream, "%s\n", err.Error())
		} else if out != "" {
			fmt.Fprintf(outStream, "%s\n", out)
		}
	}
}
