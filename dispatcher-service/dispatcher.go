package main

import "launchpad.net/~jamesh/go-dbus/trunk"
import "fmt"
import "log"

func main () {
	var (
		err error
		conn *dbus.Connection
	)

	// Connect to Session or System buses.
	if conn, err = dbus.Connect(dbus.SessionBus); err != nil {
		log.Fatal("Connection error:", err)
	}
	if err = conn.Authenticate(); err != nil {
		log.Fatal("Authentication error:", err)
	}

	fmt.Println("Hello World")
}
