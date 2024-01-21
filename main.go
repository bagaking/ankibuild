package main

import (
	"context"
	"log"
)

func main() {
	err := BuildAPKGsFromConfig(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
