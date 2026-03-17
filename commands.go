package main

import (
	"fmt"

	"github.com/TimAndrews13/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name      string
	arguments []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("username required\n")
	}

	userName := cmd.arguments[0]

	err := s.cfg.SetUser(userName)
	if err != nil {
		return fmt.Errorf("error when setting user %s: %w", userName, err)
	}

	fmt.Printf("User has been set to %s\n", cmd.arguments[0])
	return nil
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.handlers[cmd.name]
	if !ok {
		return fmt.Errorf("command %s does not exist", cmd.name)
	}

	err := handler(s, cmd)
	if err != nil {
		return fmt.Errorf("error running command %s: %w", cmd.name, err)
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}
