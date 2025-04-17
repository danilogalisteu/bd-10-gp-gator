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

	st := state{
		cfg: &cfg,
	}

	handlers := commands{
		handlers: make(map[string]func(*state, command) error),
	}

	handlers.register("login", handlerLogin)

	if len(os.Args) < 2 {
		fmt.Println("missing arguments, exiting...")
		os.Exit(1)
	}

	cmd := command{
		name: os.Args[1],
	}

	if len(os.Args) > 2 {
		cmd.args = os.Args[2:]
	} else {
		cmd.args = make([]string, 0)
	}

	err = handlers.run(&st, cmd)
	if err != nil {
		fmt.Printf("error running command %s: %v\n", cmd.name, err)
		os.Exit(1)
	}
}
