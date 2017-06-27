package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/mmatur/httppollerbeat/beater"
)

func main() {
	err := beat.Run("httppollerbeat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
