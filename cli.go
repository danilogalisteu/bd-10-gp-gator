package main

import (
	"context"
	"fmt"
	"internal/config"
	"internal/database"
	"os"
	"time"

	"github.com/google/uuid"
)

type command struct {
	name string
	args []string
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

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("login command requires one argument; provided %v", len(cmd.args))
	}

	user := cmd.args[0]
	dbUser, err := s.db.GetUser(context.Background(), user)
	if err != nil {
		os.Exit(1)
	}

	s.cfg.CurrentUserName = dbUser.Name
	err := config.SetUser(user)
	if err != nil {
		return err
	}
	s.cfg.CurrentUserName = user
	fmt.Printf("User has been set: %s\n", user)

	return nil
}
