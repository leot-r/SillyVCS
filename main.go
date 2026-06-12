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
		fmt.Print("Commands:\n  init\n  commit <file_name> <msg>\n")
	case "init":
		if len(os.Args) < 3 {
			commands.Init("")
		} else {
			commands.Init(os.Args[2])
		}
	case "commit":
		if len(os.Args) < 3 {
			fmt.Println("Please supply commit message")
			return
		}
		commands.Commit(os.Args[2])
	case "stage":
		if len(os.Args) < 3 {
			fmt.Println("Please supply file/s to stage")
			return
		}
		commands.Stage(os.Args)
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}


// TODO:
// Add staging (stage command) [IN PROGRESS]
// 	 Add check to see when staging if the file has not changed
// 	 make staging of just files work
//   Make staging append
//   Make command that clears staging
// Add tracking parent files
// Add HEAD file
// Add resore command
// Add status command
//   Will include StageFile.Print()
// Add security and comppresion to commiting
