package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// handleAutocomplete decides whether to auto-complete a command or a file path
func handleAutocomplete(command []byte, lastTabPress bool) ([]byte, bool) {
	typedStr := string(command)

	spaceIndex := strings.LastIndex(typedStr, " ")
	if spaceIndex == -1 {
		return completeCommand(command, typedStr, lastTabPress)
	}

	baseStr := typedStr[:spaceIndex+1]       
	prefixToComplete := typedStr[spaceIndex+1:]
	return completeFile(command, baseStr, prefixToComplete, lastTabPress)
}

func completeCommand(command []byte, typedPrefix string, lastTabPress bool) ([]byte, bool) {
	var matches []string

	for k := range autocompleteSet {
		if strings.HasPrefix(k, typedPrefix) {
			matches = append(matches, k) 
		}
	}

	return performCompletion(command, typedPrefix, matches, lastTabPress, true)
}

func completeFile(command []byte, baseStr string, prefix string, lastTabPress bool) ([]byte, bool) {
	dir := "."
	filePrefix := prefix

	if lastSlash := strings.LastIndex(prefix, "/"); lastSlash != -1 {
		dir = prefix[:lastSlash+1]        
		filePrefix = prefix[lastSlash+1:]
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Print("\a")
		return command, false
	}

	var matches []string
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), filePrefix) {
			matchName := entry.Name()
			
			if dir != "." {
				matchName = dir + matchName
			}
			
			if entry.IsDir() {
				matchName += "/" 
			}
			matches = append(matches, matchName)
		}
	}

	return performCompletion(command, prefix, matches, lastTabPress, false)
}

// performCompletion contains the shared LCP and double-tab logic
func performCompletion(command []byte, typedPrefix string, matches []string, lastTabPress bool, isCommand bool) ([]byte, bool) {
	if len(matches) == 1 {
		completion := matches[0][len(typedPrefix):]
		
		if isCommand || !strings.HasSuffix(matches[0], "/") {
			completion += " "
		}

		command = append(command, []byte(completion)...)
		fmt.Print(completion)
		return command, false
		
	} else if len(matches) > 1 {
		lcp := matches[0]
		for _, match := range matches[1:] {
			i := 0
			for i < len(lcp) && i < len(match) && lcp[i] == match[i] {
				i++
			}
			lcp = lcp[:i]
		}

		if len(lcp) > len(typedPrefix) {
			completion := lcp[len(typedPrefix):]
			command = append(command, []byte(completion)...)
			fmt.Print(completion)
			fmt.Print("\a")
			return command, false
		} else {
			if !lastTabPress {
				fmt.Print("\a")
				return command, true
			} else {
				fmt.Print("\r\n")

				var displayMatches []string
				for _, m := range matches {
					displayName := m
					if lastSlash := strings.LastIndex(m[:len(m)-1], "/"); lastSlash != -1 {
						displayName = m[lastSlash+1:]
					}
					displayMatches = append(displayMatches, strings.TrimRight(displayName, "/"))
				}

				sort.Strings(displayMatches)

				fmt.Print(strings.Join(displayMatches, "  ") + "\r\n")
				fmt.Print("$ " + string(command))
				return command, false
			}
		}
	}

	fmt.Print("\a")
	return command, false
}
