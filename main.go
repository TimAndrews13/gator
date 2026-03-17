package main

import (
	"fmt"
	"os"

	"github.com/TimAndrews13/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("Error Reading Config File: %v", err)
		return
	}

	cfgState := state{
		cfg: &cfg,
	}

	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)

	args := os.Args
	if len(args) < 2 {
		fmt.Printf("no arguments supplied in cli\ntry again\n")
		os.Exit(1)
	}

	cmd := command{
		name:      args[1],
		arguments: args[2:],
	}

	err = cmds.run(&cfgState, cmd)
	if err != nil {
		fmt.Printf("error running %s: %v", cmd.name, err)
		os.Exit(1)
	}
	/*
	   cfg, err = config.Read()

	   	if err != nil {
	   		fmt.Printf("Error Reading Rewritten Config File: %v", err)
	   		return
	   	}

	   fmt.Print(cfg)
	*/
}
