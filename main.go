package main

import (
	"log"
	"net/http"

	"github.com/jaschaephraim/lrserver"

	"gopkg.in/fsnotify.v1"
)

func main() {
	// Create file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalln(err)
	}
	defer watcher.Close()

	// Add dir to watcher
	if err := watcher.Add("."); err != nil {
		log.Fatalln(err)
	}

	// Create and start LiveReload server
	lr, err := lrserver.New(lrserver.DefaultName, lrserver.DefaultPort)
	if err != nil {
		log.Fatal(err)
	}
	go lr.ListenAndServe()

	// Start goroutine that requests reload upon watcher event
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)
				lr.Reload(event.Name)
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("."))))
	http.ListenAndServe(":8080", nil)
}
