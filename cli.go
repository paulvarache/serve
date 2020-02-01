package main

import "flag"

import "os"

type CliOptions struct {
	SpaIndex string
	Dir      string
	Open     bool
	Port     int
	Hostname string
}

func FlagString(target *string, name string, shorthand string, defaultVal string, usage string) {
	flag.StringVar(target, name, defaultVal, usage)
	flag.StringVar(target, shorthand, defaultVal, usage+" (shorthand)")
}

func FlagBool(target *bool, name string, shorthand string, defaultVal bool, usage string) {
	flag.BoolVar(target, name, defaultVal, usage)
	flag.BoolVar(target, shorthand, defaultVal, usage+" (shorthand)")
}

func FlagInt(target *int, name string, shorthand string, defaultVal int, usage string) {
	flag.IntVar(target, name, defaultVal, usage)
	flag.IntVar(target, shorthand, defaultVal, usage+" (shorthand)")
}

func NewCliOptions() *CliOptions {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return &CliOptions{
		SpaIndex: "",
		Dir:      cwd,
		Open:     false,
		Port:     8000,
		Hostname: "localhost",
	}
}
