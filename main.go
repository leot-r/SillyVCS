package main

import (
	"fmt"
	"os"

	"SillyVCS/commands"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("run <help> for commands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "help":
		fmt.Print("Commands:\n  init <>\n  commit <>\n")
	case "init":
		if len(os.Args) < 3 {
			commands.Init("")
		} else {
			commands.Init(os.Args[2])
		}
	case "commit":
		if len(os.Args) < 3 {
			fmt.Println("Please supply file to commit")
			return
		}
		commands.CommitFile(os.Args[2])
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
