package main

import (
	"fmt"
	"internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name    string
	args	[]string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.handlers[cmd.name]
	if !ok {
		return fmt.Errorf("handler for %s not registered", cmd.name)
	}
	return handler(s, cmd)
}
