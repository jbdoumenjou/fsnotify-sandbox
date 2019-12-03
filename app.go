package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/fsnotify/fsnotify"
)

func main() {
	var target = flag.String("target", "file", "use file, dir or both as watcher target")
	var reload = flag.Bool("reload", false, "reload the watcher if true")
	flag.Parse()

	fmt.Println("Started with target: '", *target, "' and reload: '", *reload, "'")

	var filePaths []string
	fileToWatch := "/etc/traefik/dyn.toml"

	switch *target {
	case "dir":
		filePaths = []string{path.Dir(fileToWatch)}
	case "both":
		filePaths = []string{fileToWatch, path.Dir(fileToWatch)}
	case "file":
		filePaths = []string{fileToWatch}
	default:
		fmt.Println("error: unknown argument: '", *target, "'. Choose in the following: 'file', 'dir' or 'both'")
		os.Exit(1)
	}

	err := newWatcher(filePaths, *reload)
	if err != nil {
		fmt.Println("error: unable to create a file watcher:", err)
		os.Exit(1)
	}
}

func newWatcher(paths []string, reload bool) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	done := make(chan bool)
	go func(reload bool) {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				switch event.Op {
				case fsnotify.Create:
					log.Println("Created")

				case fsnotify.Write:
					log.Println("Write")

				case fsnotify.Chmod:
					log.Println("Chmod")

				case fsnotify.Rename:
					log.Println("Moved")

					if reload {
						if err := watcher.Remove("/etc/traefik/dyn.toml"); err != nil {
							log.Fatal(err)
							return
						}

						watcher.Add("/etc/traefik/dyn.toml")
					}

				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}(reload)

	for _, path := range paths {
		err = watcher.Add(path)
		if err != nil {
			log.Fatal(err)
		}
	}
	<-done

	return nil
}
