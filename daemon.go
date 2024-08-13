package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
)

func printHelp() {
	fmt.Println(`Usage: compile-daemon-watcher [OPTIONS]

Options:
  --path <path>   Path to watch for file changes
  --help          Show this help message`)
}

func main() {
	// Определение флагов командной строки
	pathToWatch := flag.String("path", "", "Path to watch")
	showHelp := flag.Bool("help", false, "Show help message")
	flag.Parse()

	// Проверка флага --help
	if *showHelp {
		printHelp()
		os.Exit(0)
	}

	// Проверка обязательного параметра
	if *pathToWatch == "" {
		color.New(color.BgBlue, color.FgWhite).Printf("Usage: compile-daemon-watcher --path=<path-to-watch>\n")
		os.Exit(1)
	}

	// Создание нового наблюдателя файлов
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		color.New(color.FgRed).Printf("Failed to create watcher: %v\n", err)
		log.Fatalf("Failed to create watcher: %v", err)
	}
	defer watcher.Close()

	// Добавление пути к наблюдателю
	err = watcher.Add(*pathToWatch)
	if err != nil {
		color.New(color.FgRed).Printf("Failed to add path to watcher: %v\n", err)
		log.Fatalf("Failed to add path to watcher: %v", err)
	}

	// Основной цикл для обработки событий
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
					buildProject(*pathToWatch)
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

// buildProject выполняет сборку проекта
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
