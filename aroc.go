package main

import (
	"fmt"
	"github.com/sdegutis/go.fsevents"
	"log"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: aroc DIRECTORY COMMAND [ARGS…]")
		os.Exit(1)
	}

	ch := fsevents.WatchPaths([]string{os.Args[1]})

	var cmd *exec.Cmd

	go func() {
		for _ = range ch {
			log.Println("Changes in directory, restarting")
			cmd.Process.Signal(os.Interrupt)
		}
	}()

	for {
		cmd = exec.Command(os.Args[2])
		cmd.Args = os.Args[2:]
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		err := cmd.Run()
		if err != nil {
			if err, ok := err.(*exec.ExitError); !ok {
				log.Fatal(err)
			}
		}
	}
}
