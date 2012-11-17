package main

import (
	"fmt"
	"github.com/howeyc/fsnotify"
	"log"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: aroe DIRECTORY COMMAND [ARGS…]")
		os.Exit(1)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Watch(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	var cmd *exec.Cmd

	go func() {
		for {
			select {
			case _ = <-watcher.Event:
				log.Println("Changes in directory, restarting")
				cmd.Process.Signal(os.Interrupt)
			case err := <-watcher.Error:
				log.Fatal("error:", err)
			}
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
