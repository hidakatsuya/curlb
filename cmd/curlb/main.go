package main

import (
	"flag"
	"fmt"
	"os"

	"curlb/internal/app"
)

var version = "0.1.0"

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "curlb %s - Simple curl command builder\n", version)
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: curlb [--version|--help] [-- additional-curl-args]\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Examples:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  curlb\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  curlb -- --connect-timeout 2  # pass extra curl flags\n")
	}

	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		return
	}

	extraArgs := flag.Args()

	if err := app.Run(extraArgs); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
