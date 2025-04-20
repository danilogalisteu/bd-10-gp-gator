package main

import (
	"database/sql"
	"fmt"
	"internal/config"
	"internal/database"
	"os"

	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

const dbURL = "postgres://postgres:abc123@localhost:5432/gator?sslmode=disable"

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	st := state{
		db:  database.New(db),
		cfg: &cfg,
	}

	handlers := commands{
		handlers: make(map[string]func(*state, command) error),
	}

	handlers.register("login", handlerLogin)
	handlers.register("register", handlerRegister)
	handlers.register("reset", handlerReset)
	handlers.register("users", handlerUsers)
	handlers.register("agg", handlerAggregator)
	handlers.register("addfeed", handlerAddFeed)
	handlers.register("feeds", handlerFeeds)

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
