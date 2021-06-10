package main

import (
	"os"
	"strings"
)

// Args contains the command line arguments as key value pairs. In sodago this
// is handled quite simple: in the command line a key is identified via a hyphen
// prefix (`-`) followed by the corresponding value (e.g. `-data ./data -port
// 8080`).
type Args map[string]string

// ParseArgs parses the command line arguments
func ParseArgs() Args {
	args := make(map[string]string)
	if len(os.Args) < 2 {
		return args
	}
	flag := ""
	for i, val := range os.Args {
		if i == 0 {
			continue
		}
		arg := strings.TrimSpace(val)
		if flag != "" {
			args[flag] = arg
			flag = ""
			continue
		}
		if strings.HasPrefix(arg, "-") {
			flag = arg
		}
	}
	return args
}

func (args Args) GetOrDefault(key, defaultValue string) string {
	val, ok := args[key]
	if !ok {
		return defaultValue
	} else {
		return val
	}
}

func (args Args) DataDir() string {
	return args.GetOrDefault("-data", "data")
}

func (args Args) Port() string {
	return args.GetOrDefault("-port", "8080")
}
