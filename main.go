package main

import (
	"context"
	"log"
)

func main() {
	err := BuildAPKGsFromToml(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
