package config

import "flag"

var (
	FlagRunAddr  string
	FlagBaseAddr string
)

func parseFlags() {
	flag.StringVar(&FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&FlagBaseAddr, "b", "http://localhost:8080", "server address before shorten URL")

	flag.Parse()
}
