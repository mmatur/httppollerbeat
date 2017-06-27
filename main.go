package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/mmatur/httppollerbeat/beater"
)

var name = "httppollerbeat"

func main() {
	err := beat.Run(name, "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
