package main

import (
	"fmt"
	"go-armory/examples/future_monad/future"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 1 {
		fmt.Println("choose program eg. 0,1,2,3")
		os.Exit(1)
	}
	switch args[1] {
	case "0":
		future.Start()
	}

}
