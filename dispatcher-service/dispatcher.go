package main

import "launchpad.net/~jamesh/go-dbus/trunk"
import "fmt"
import "log"

func handleURLMessage (message *dbus.Message) {

}

func main () {
	var (
		err error
		conn *dbus.Connection
		messages = make(chan *dbus.Message)
	)

	// Connect to Session bus.
	if conn, err = dbus.Connect(dbus.SessionBus); err != nil {
		log.Fatal("Connection error:", err)
	}

	conn.RegisterObjectPath("/com/canonical/URLDispatcher", messages);

	fmt.Println("Hello World")

	for {
		message := <- messages
		fmt.Println("Got Message:", message)
		switch message.Member {
		case "DispatchURL":
			go handleURLMessage(message)
		default:
			errormsg := dbus.NewErrorMessage(message, "InvalidInterface", "Unsupported function")
			conn.Send(errormsg)
		}
	}
}
