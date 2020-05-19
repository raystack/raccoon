package main

import (
	"clickstream-service/app"
	"log"
)

func main() {
	err := app.Run()
	if err != nil {
		log.Fatal("init failure", err)
	}
}
