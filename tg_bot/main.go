package main

import (
	"github.com/gulldan/cp2024omsk-pmsk/bot"
	"github.com/gulldan/cp2024omsk-pmsk/config"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	if err := bot.New(&cfg); err != nil {
		panic(err)
	}
}
