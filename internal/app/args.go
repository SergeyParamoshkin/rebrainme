package app

import "flag"

type Args struct {
	ConfigFileName string
}

func NewArgs() *Args {
	args := &Args{}

	flag.StringVar(&args.ConfigFileName, "c", "config.yaml", "config")
	flag.Parse()

	return args
}
