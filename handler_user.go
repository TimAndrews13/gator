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
