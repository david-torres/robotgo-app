package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	mw "github.com/labstack/echo/middleware"
	"github.com/zserge/webview"
	"golang.org/x/net/websocket"
)

// Message to receive/send
type Message struct {
	Message string `json:"msg"`
}

func main() {
	// channel to post window handles to
	messages := make(chan string, 1)

	// start the webserver
	e := getServer(messages)
	go e.Run(standard.New(":3000"))

	// start the window handle listener
	go startWindowListener(messages)

	// start the mouse listener
	go startMouseListener(messages)

	// start the clipboard listener
	go startClipboardListener(messages)

	// start webview UI
	webview.Open("RobotGo App", "http://localhost:3000", 400, 300, true)
}

func getServer(messages chan string) *echo.Echo {
	// init the web server
	e := echo.New()

	// init app-wide middleware
	e.Use(mw.Logger())
	e.Use(mw.Recover())
	e.Use(mw.Gzip())

	// routes
	e.Static("/", "public")
	e.File("/", "public/index.html")
	e.GET("/ws", standard.WrapHandler(websocket.Handler(func(ws *websocket.Conn) {
		// loop forever receiving from channel
		for {
			msg := <-messages
			log.Println("Received message: " + msg)

			outgoing := new(Message)
			outgoing.Message = msg
			err := websocket.JSON.Send(ws, outgoing)

			if err != nil {
				log.Println("websocket send error")
				log.Println(err.Error())
				continue
			}
		}
	})))

	return e
}

func startWindowListener(messages chan string) {
	// loop forever polling for a window title every second and post to channel
	for {
		time.Sleep(1 * time.Second)
		title := robotgo.GetTitle()

		log.Println("Sending window title: " + title)
		messages <- "Window title: " + title
	}
}

func startMouseListener(messages chan string) {
	// loop forever polling for mouse position every second and post to channel
	for {
		time.Sleep(1 * time.Second)
		x, y := robotgo.GetMousePos()
		mpos := fmt.Sprintf("(%d, %d)", x, y)
		log.Println("Sending mouse position: " + mpos)

		messages <- "Mouse position: " + mpos
	}
}

func startClipboardListener(messages chan string) {
	// loop forever polling for clipboard every second and post to channel
	for {
		time.Sleep(1 * time.Second)
		text, _ := robotgo.ReadAll()

		log.Println("Sending clipboard: " + text)
		messages <- "Clipboard: " + text
	}
}
