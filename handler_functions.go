package main

import (
	"context"
	"fmt"
	"time"

	"github.com/TimAndrews13/gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("username required\n")
	}

	userName := cmd.arguments[0]

	_, err := s.db.GetUser(context.Background(), userName)
	if err != nil {
		return fmt.Errorf("user does not exist: %s\n", userName)
	}

	err = s.cfg.SetUser(userName)
	if err != nil {
		return fmt.Errorf("error when setting user %s: %w\n", userName, err)
	}

	fmt.Printf("User has been set to %s\n", cmd.arguments[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("name required\n")
	}

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
	})
	if err != nil {
		return fmt.Errorf("error creating new user record: %v\n", err)
	}

	s.cfg.SetUser(cmd.arguments[0])

	fmt.Printf("user created: %v\n", user)

	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error resetting users table: %v\n", err)
	}

	fmt.Printf("user table reset\n")

	return nil
}

func handlerUsersList(s *state, cmd command) error {
	var users []database.User
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error retrieving users: %v\n", err)
	}

	for i := 0; i < len(users); i++ {
		if users[i].Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", users[i].Name)
		} else {
			fmt.Printf("* %s\n", users[i].Name)
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	/*  ADD BACK IN WHEN WE DO NOT WANT A HARDCODED URL
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("url required\n")
	}
	*/

	rssFeed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("error retrieving rss feed from url: %v", err)
	}

	fmt.Printf("Feed: %+v\n", rssFeed)

	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 2 {
		return fmt.Errorf("not enough arguments\nname and url required\n")
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
		Url:       cmd.arguments[1],
		UserID:    uuid.NullUUID{UUID: user.ID, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("error creating new feed record: %v\n", err)
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating new feed_follow record: %v", err)
	}

	fmt.Printf("feed created: %v\n", feed)

	return nil
}

func handlerGetFeeds(s *state, cmd command) error {
	var feeds []database.GetFeedsRow
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error retrieving feeds: %v\n", err)
	}

	for i := 0; i < len(feeds); i++ {
		fmt.Printf("* Feed: %s | URL: %s | User: %s\n", feeds[i].FeedName, feeds[i].Url, feeds[i].UserName)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("url required\n")
	}

	feed, err := s.db.GetFeed(context.Background(), cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("error pulling feed by url: %v\n", err)
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating new feed_follow record: %v\n", err)
	}

	fmt.Printf("Feed Followed: %s\n", feed.Name)
	fmt.Printf("Followed By: %s\n", user.Name)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return fmt.Errorf("error retrieving feeds for current user: %v\n", err)
	}

	for i := 0; i < len(feeds); i++ {
		fmt.Printf("* %s\n", feeds[i].FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("url required\n")
	}

	feed, err := s.db.GetFeed(context.Background(), cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("error retrieving feed id %v\n", err)
	}

	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error unfollowing feed for user: %v\n", err)
	}

	fmt.Printf("Feed %s Unfollowed by %s\n", feed.Name, user.Name)
	return nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s1 *state, cmd1 command) error {
		currentUser, err := s1.db.GetUser(context.Background(), s1.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("error pulling current user: %v\n", err)
		}

		resultErr := handler(s1, cmd1, currentUser)

		return resultErr
	}
}
