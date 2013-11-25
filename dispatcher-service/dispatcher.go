package main

import "launchpad.net/~jamesh/go-dbus/trunk"
import "fmt"
import "log"
//import "regexp"

type urlHandler struct {
	application string
	regex string
}

var urlHandlers []urlHandler

func initHandlers () {
	urlHandlers = append(urlHandlers,
		urlHandler{application: "webbrowser-app", regex: "^http://"},
		urlHandler{application: "webbrowser-app", regex: "^https://"})
}

func handleURLMessage (message *dbus.Message) {

}

func main () {
	var (
		err error
		conn *dbus.Connection
		messages = make(chan *dbus.Message)
	)

	initHandlers()

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
