package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

func main() {
	err := newWatcher()
	if err != nil {
		fmt.Println("error: unable to create a file watcher:", err)
		os.Exit(1)
	}
}

func newWatcher() error {
	//filePath := path.Dir("/etc/traefik/dyn.toml")
	filePath := "/etc/traefik/dyn.toml"

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
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

					// TODO take care of the link cf certdumper
					//respawnFile(event.Name)

					// add the file back to watcher, since it is removed from it
					// when file is moved or deleted
					//log.Printf("add to watcher file:  %s\n", filePath)
					// add appears to be concurrently safe so calling from multiple go routines is ok
					//err = watcher.Add(filePath)
					//if err != nil {
					//	log.Fatal(err)
					//}

					// there is  not need to break the loop
					// we just continue waiting for events from the same watcher

				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(filePath)
	if err != nil {
		log.Fatal(err)
	}
	<-done

	return nil
}

//func respawnFile(filepath string) {
//	log.Printf("re creating file %s\n", filepath)
//
//	// you just need the os.Create()
//	respawned, err := os.Create(filepath)
//	if err != nil {
//		log.Fatalf("Err re-spawning file: %v", filepath)
//	}
//	defer respawned.Close()
//
//	// there is no need to call monitorFile again, it never returns
//	// the call to "go monitorFile(filepath)" was causing another go routine leak
//}
