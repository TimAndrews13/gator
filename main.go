package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/TimAndrews13/gator/internal/config"
	"github.com/TimAndrews13/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/gator")
	if err != nil {
		fmt.Printf("error connecting to postgres database: %v", err)
		return
	}

	dbQueries := database.New(db)

	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("error reading config file: %v", err)
		return
	}

	cfgState := state{
		db:  dbQueries,
		cfg: &cfg,
	}

	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsersList)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)

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
		fmt.Printf("error running %s: %v\n", cmd.name, err)
		os.Exit(1)
	}
}
