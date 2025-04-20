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

func handlerAddFeed(s *state, cmd command, dbUser database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("addfeed command requires two arguments; provided %v", len(cmd.args))
	}

	name := cmd.args[0]
	url := cmd.args[1]
	dbFeed, err := s.db.CreateFeed(context.Background(),
		database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      name,
			Url:       url,
			UserID:    dbUser.ID,
		})
	if err != nil {
		os.Exit(1)
	}

	fmt.Printf("Feed has been added: %s\n", dbFeed.Name)

	dbFollow, err := s.db.CreateFeedFollow(context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    dbUser.ID,
			FeedID:    dbFeed.ID,
		})
	if err != nil {
		os.Exit(1)
	}

	fmt.Printf("User %s is now following feed '%s'\n", dbFollow.UserName, dbFollow.FeedName)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("feeds command doesn't require arguments; provided %v", len(cmd.args))
	}

	dbFeeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		os.Exit(1)
	}

	for _, feed := range dbFeeds {
		fmt.Printf("[%s] '%s' at %s\n", feed.UserName.String, feed.Name, feed.Url)
	}

	return nil
}

func handlerAddFollow(s *state, cmd command, dbUser database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("follow command requires one argument; provided %v", len(cmd.args))
	}

	url := cmd.args[0]
	dbFeed, err := s.db.GetFeed(context.Background(), url)
	if err != nil {
		os.Exit(1)
	}

	dbFollow, err := s.db.CreateFeedFollow(context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    dbUser.ID,
			FeedID:    dbFeed.ID,
		})
	if err != nil {
		os.Exit(1)
	}

	fmt.Printf("User %s is now following feed '%s'\n", dbFollow.UserName, dbFollow.FeedName)

	return nil
}

func handlerFollowing(s *state, cmd command, dbUser database.User) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("following command doesn't require arguments; provided %v", len(cmd.args))
	}

	dbFollows, err := s.db.GetFeedFollowsForUser(context.Background(), dbUser.ID)
	if err != nil {
		os.Exit(1)
	}

	for _, dbFollow := range dbFollows {
		fmt.Printf("Following '%s'\n", dbFollow.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, dbUser database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("unfollow command requires one argument; provided %v", len(cmd.args))
	}

	url := cmd.args[0]
	err := s.db.DeleteFollow(context.Background(),
		database.DeleteFollowParams{
			UserID: dbUser.ID,
			Url:    url,
		})
	if err != nil {
		os.Exit(1)
	}

	fmt.Printf("Feed %s was unfollowed\n", url)

	return nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		dbUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			os.Exit(1)
		}
		return handler(s, cmd, dbUser)
	}
}
