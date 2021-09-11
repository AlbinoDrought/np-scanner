package main

import (
	"log"

	"go.albinodrought.com/neptunes-pride/internal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
