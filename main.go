package main

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
	rpmdb "github.com/knqyf263/go-rpmdb/pkg"
)

func main() {
	watchedPackages := []string{"pingus"}

	db, err := rpmdb.Open("/var/lib/rpm/Packages")
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
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
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					for _, pkgname := range watchedPackages {
						pkginfo, _ := db.Package(pkgname)
						fmt.Printf("\t%+v\n", *pkginfo)
					}
					//log.Println("Some RPM package was installed, updated or removed.")
					//log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("/var/lib/rpm/Packages")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
