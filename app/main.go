package main

import (
	"bufio"
	"fmt"
	"os"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	for {
		fmt.Print("$ ")
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
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
				index = i;
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
			fmt.Fprintf(outStream, "%s\n", out);
		}
	}
}
