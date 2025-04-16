package main

import (
	"fmt"
	"internal/config"
	"os"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(cfg)
}
