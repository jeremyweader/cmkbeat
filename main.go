package main

import (
	"github.com/elastic/beats/libbeat/beat"
	"github.com/jeremyweader/cmkbeat/beater"
)

func main() {
	beat.Run("cmkbeat", "", beater.New)
}
