package main

import (
	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 2 {
		color.New(color.BgBlue, color.FgWhite).Printf("Usage: ./CompileDaemonWatcher <path-to-watch>\n")
		os.Exit(1)
	}

	pathToWatch := os.Args[1]

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		color.New(color.FgRed).Printf("Failed to create watcher : %v\n", err)
		log.Fatalf("Failed to create watcher: %v", err)
	}
	defer watcher.Close()

	err = watcher.Add(pathToWatch)
	if err != nil {
		color.New(color.FgRed).Printf("Failed to add path to watcher: %v\n", err)
		log.Fatalf("Failed to add path to watcher: %v", err)
	}

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					color.New(color.FgGreen).Printf("Modified file: %s\n", event.Name)
					buildProject(pathToWatch)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				color.New(color.FgRed).Printf("Error: %v\n", err)
				log.Println("Error:", err)
			}
		}
	}()

	<-done
}

func buildProject(path string) {
	color.New(color.FgYellow).Printf("Building project...\n")
	cmd := exec.Command("go", "build", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		color.New(color.FgRed).Printf("Error building project: %v\n", err)
	} else {
		color.New(color.FgGreen).Printf("Build successful\n")
	}
}
