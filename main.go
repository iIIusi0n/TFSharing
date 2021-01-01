package main

import (
	"flag"
	"log"
)

func main() {
	isServerPtr := flag.Bool("server", false, "Is it server?")
	flag.Parse()
	isServer := *isServerPtr
	if isServer {
		InitServer()
	} else {
		InitClient()
	}
}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}