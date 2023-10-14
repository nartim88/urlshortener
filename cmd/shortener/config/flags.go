package config

import "flag"

// parseFlags парсит флаги командной строки
func parseFlags() {
	flag.StringVar(&CFG.RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&CFG.BaseURL, "b", "http://localhost:8080", "server address before shorten URL")

	flag.Parse()
}
