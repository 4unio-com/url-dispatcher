package main

import "launchpad.net/~jamesh/go-dbus/trunk"
import "fmt"
import "log"
import "regexp"

type urlHandler struct {
	application string
	regex string
}

func (handler urlHandler) handleURL (url string, returnchan chan bool, conn *dbus.Connection) {
	if matched, _ := regexp.MatchString(handler.regex, url); !matched {
		fmt.Println("Unable to match URL", url, "with regular expression", handler.regex)
		returnchan <- false
		return
	}

	message := dbus.NewMethodCallMessage("com.ubuntu.Upstart", "/com/ubuntu/Upstart/jobs/application", "com.ubuntu.Upstart0_6.Job", "Start")
	
	var env []string;
	env = append(env, fmt.Sprintf("APP_ID=%s", handler.application))
	message.AppendArgs(env, true);

	if _, error := conn.SendWithReply(message); error != nil {
		fmt.Println("Unable to start application")
		returnchan <- false
	}

	returnchan <- true
}

var urlHandlers []urlHandler

func initHandlers () {
	urlHandlers = append(urlHandlers,
		urlHandler{application: "webbrowser-app", regex: "^http://"},
		urlHandler{application: "webbrowser-app", regex: "^https://"})
}

func handleURLMessage (message *dbus.Message, conn *dbus.Connection) {
	var url string
	handlersHandled := make(chan bool, len(urlHandlers))

	message.Args(&url)

	for _, handler := range urlHandlers {
		go handler.handleURL(url, handlersHandled, conn)
	}

	handlerCount := 0
	for i := 0; i < len(urlHandlers); i++ {
		if <-handlersHandled {
			handlerCount = handlerCount + 1
		}
	}

	var returnmessage *dbus.Message
	if handlerCount > 0 {
		fmt.Println("Finished handling", url)
		returnmessage = dbus.NewMethodReturnMessage(message)
	} else {
		fmt.Println("Invalid URL", url)
		returnmessage = dbus.NewErrorMessage(message, "InvalidURL", "No handlers for URL")
	}

	conn.Send(returnmessage)
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

	fmt.Println("DBus Unique Name", conn.UniqueName)

	conn.RegisterObjectPath("/com/canonical/URLDispatcher", messages);

	for {
		message := <- messages
		fmt.Println("Got Message:", message)
		switch message.Member {
		case "DispatchURL":
			go handleURLMessage(message, conn)
		default:
			errormsg := dbus.NewErrorMessage(message, "InvalidInterface", "Unsupported function")
			conn.Send(errormsg)
		}
	}
}
