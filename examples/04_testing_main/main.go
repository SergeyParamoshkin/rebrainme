package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

var (
	stdout io.Writer = os.Stdout
	stderr io.Writer = os.Stderr
	args   []string  = os.Args[1:]
)

func main() {
	os.Exit(mainWithExitCode())
}

func mainWithExitCode() int {
	fs := flag.NewFlagSet("myapp", flag.ContinueOnError)
	fs.SetOutput(stderr)

	name := fs.String("name", "World", "name to greet")
	verbose := fs.Bool("v", false, "verbose output")
	repeat := fs.Int("repeat", 1, "number of times to repeat greeting")

	if err := fs.Parse(args); err != nil {
		return 2
	}

	if *verbose {
		fmt.Fprintln(stdout, "[VERBOSE] Starting application...")
	}

	for i := 0; i < *repeat; i++ {
		fmt.Fprintf(stdout, "Hello, %s!\n", *name)
	}

	if *verbose {
		fmt.Fprintln(stdout, "[VERBOSE] Application finished successfully")
	}

	return 0
}

func run() (int, error) {
	fs := flag.NewFlagSet("myapp", flag.ContinueOnError)

	name := fs.String("name", "World", "name to greet")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return 2, err
	}

	result := greet(*name)
	fmt.Println(result)

	return 0, nil
}

func greet(name string) string {
	if name == "" {
		return "Hello, World!"
	}
	return fmt.Sprintf("Hello, %s!", name)
}
