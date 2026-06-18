package main

type RedirectConfig struct {
	Stdout bool
	Stderr bool
	Append bool
}

var redirectionMap = map[string]RedirectConfig{
	">":   {Stdout: false, Stderr: true, Append: false},
	"1>":  {Stdout: false, Stderr: true, Append: false},
	">>":  {Stdout: false, Stderr: true, Append: true},
	"1>>": {Stdout: false, Stderr: true, Append: true},
	"2>":  {Stdout: true, Stderr: false, Append: false},
	"2>>": {Stdout: true, Stderr: false, Append: true},
}
