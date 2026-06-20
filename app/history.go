package main

import (
	"os"
	"strings"
)

var history [][]string

func init() {
	history = make([][]string, 0)
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
	return nil
}
