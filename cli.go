package main

import (
	"context"
	"fmt"
	"internal/config"
	"internal/database"
	"internal/rss"
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
	err = config.SetUser(dbUser.Name)
	if err != nil {
		return err
	}

	fmt.Printf("User has been set: %s\n", dbUser.Name)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("register command requires one argument; provided %v", len(cmd.args))
	}

	user := cmd.args[0]
	dbUser, err := s.db.CreateUser(context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      user,
		})
	if err != nil {
		os.Exit(1)
	}

	fmt.Printf("User has been created: %s\n", dbUser.Name)

	s.cfg.CurrentUserName = dbUser.Name
	err = config.SetUser(dbUser.Name)
	if err != nil {
		return err
	}

	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("reset command doesn't require arguments; provided %v", len(cmd.args))
	}

	err := s.db.ResetUsers(context.Background())
	if err != nil {
		os.Exit(1)
	}

	fmt.Println("User table has been reset")

	return nil
}

func handlerUsers(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("users command doesn't require arguments; provided %v", len(cmd.args))
	}

	dbUsers, err := s.db.GetUsers(context.Background())
	if err != nil {
		os.Exit(1)
	}

	for _, user := range dbUsers {
		name := user.Name
		if name == s.cfg.CurrentUserName {
			name += " (current)"
		}
		fmt.Printf("* %s\n", name)
	}

	return nil
}

func handlerAggregator(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("users command doesn't require arguments; provided %v", len(cmd.args))
	}

	feedURL := "https://www.wagslane.dev/index.xml"
	feed, err := rss.FetchFeed(context.Background(), feedURL)
	if err != nil {
		return err
	}
	fmt.Println(feed)

	return nil
}
