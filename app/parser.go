package main

import "strings"

func parseInput(input string) []string {
	var args []string
	var current strings.Builder
	inSingleQuote := false
	inDoubleQuote := false

	input = strings.TrimSpace(input)

	for i := 0; i < len(input); i++ {
		c := input[i]

		if c == '\'' && !inDoubleQuote {
			inSingleQuote = !inSingleQuote
			continue
		}

		if c == '"' && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
			continue
		}

		if c == '\\' && !inSingleQuote && !inDoubleQuote && i < len(input) - 1 {
			current.WriteByte(input[i+1])
			i++; 
			continue
		}

		if (c == ' ' || c == '\t') && !inSingleQuote && !inDoubleQuote {
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
			continue
		}

		current.WriteByte(c)
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}

	return args
}
