package main

import (
	"flag"
	"log"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/deckarep/gosx-notifier"
	"github.com/tubbebubbe/transmission"
)

func main() {

	var server, port, user, password string
	flag.StringVar(&server, "server", "127.0.0.1", "Hostname of transmission server")
	flag.StringVar(&port, "port", "9091", "Port of transmission server")
	flag.StringVar(&user, "user", "admin", "Username")
	flag.StringVar(&password, "password", "****", "Password")

	flag.Parse()

	// Connect to transmission client
	client := transmission.New("http://"+server+":"+port, user, password)

	// Test the connection
	_, err := client.GetTorrents()
	if err != nil {
		log.Panic(err)
	}

	log.Println("Listening for magnet links...")

	// Loop until sigterm (Ctrl-C)
	for {
		text, err := clipboard.ReadAll()
		if err != nil {
			panic(err)
		}

		// Check for a magnet link in the clipboard
		if strings.HasPrefix(text, "magnet:") {

			// Create a new Add Command from the magnet link
			cmd, err := transmission.NewAddCmdByMagnet(text)
			if err != nil {
				panic(err)
			}

			// Exectute the Add Command on the remote transmission server
			tAdded, err := client.ExecuteAddCommand(cmd)
			if err != nil {
				panic(err)
			}

			// If we added a new torrent, notify
			name := tAdded.Name
			if len(name) >= 1 {
				notify(name)
				log.Println("Added magnet link for:", name)
			}

			// Clear the clipboard
			err = clipboard.WriteAll("")
			if err != nil {
				panic(err)
			}
		}

		// Sleep for 1 second, because everyone needs a break?
		time.Sleep(time.Second)
	}

}

// Notify OSX
func notify(name string) {
	note := gosxnotifier.NewNotification(name)
	note.Title = "Torrent Added!"
	note.Sound = gosxnotifier.Tink

	err := note.Push()
	if err != nil {
		panic(err)
	}
}
