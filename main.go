package main

import (
	"math/rand"
	"time"

	"github.com/comp500/filebisect/cmd"
)

func main() {
	// Ensure rand is seeded
	rand.Seed(time.Now().Unix())

	cmd.Execute()
}
