package main

import (
	"fmt"

	"github.com/TimAndrews13/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("Error Reading Config File: %v", err)
		return
	}

	err = cfg.SetUser("tim")
	if err != nil {
		fmt.Printf("Error Setting User in Config Struct: %v", err)
		return
	}

	cfg, err = config.Read()
	if err != nil {
		fmt.Printf("Error Reading Rewritten Config File: %v", err)
		return
	}

	fmt.Print(cfg)
}
