package main

import (
	"os"
	"strings"
)

var history [][]string
var lastAppendIdx int

func init() {
	history = make([][]string, 0)
	file := os.Getenv("HISTFILE")
	if len(file) != 0 {
		AppendHistory(file)
	}
}

func AppendHistory(filename string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	strContent := strings.Split(string(content), "\n")
	for _, line := range strContent {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		result := parseInput(line)
		if len(result) > 0 {
			history = append(history, result)
		}
	}
	lastAppendIdx = len(history)
	return nil
}

func WriteHistory(filename string, appendFlag bool) error {
	t := os.O_TRUNC
	startIndex := 0
	if appendFlag {
		t = os.O_APPEND
		startIndex = lastAppendIdx
	}

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|t, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for i := startIndex; i < len(history); i++ {
		line := strings.Join(history[i], " ")
		_, err := file.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	lastAppendIdx = len(history)
	return nil
}
