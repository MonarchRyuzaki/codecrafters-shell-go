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

func WriteHistory(filename string, append bool) error {
	t := os.O_TRUNC
	if append {
		t = os.O_APPEND
	}
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|t, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, cmd := range history {
		line := strings.Join(cmd, " ")
		_, err := file.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}
