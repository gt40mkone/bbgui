package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"xt995"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

type ewApp struct {
	displayText  *widget.Label
	appStatus    *widget.Label
	serverStatus *widget.Label

	bServerUp bool

	iMsgCount int

	//sClientAddr string
	//conClient   *websocket.Conn

	client *xt995.WsClient
	server *xt995.WsServer

	//	chToClient   chan string
	//	chFromClient chan string

	lApp *log.Logger //	client log
	lCli *log.Logger //	client log
	lSrv *log.Logger //	server log
}

var a ewApp

func loggerSetup(sPath string, sPrefix string) *log.Logger {

	file, err := os.OpenFile(sPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	//	l := log.New(file, sPrefix, log.Ldate|log.Ltime|log.Lshortfile)
	l := log.New(file, sPrefix, log.Ldate|log.Ltime)

	l.Print("============================================")

	return l
}

func clientConnect(a *ewApp) error {

	a.lCli.Print("client init phase")

	a.client.Connect("127.0.0.1:5757")

	err := a.client.Connect("localhost:5757")
	if err != nil {
		return err
	}

	//	defer c.Close()

	a.iMsgCount = 0

	a.lCli.Printf("client connection established")

	return nil
}

func clientSend(a *ewApp, sMsg string, iCount int) error {

	if !a.bServerUp {
		return errors.New(updateAppStatus(a, "no server available"))
	}
	a.iMsgCount++
	myString := fmt.Sprintf("%s : %d", sMsg, a.iMsgCount)

	a.lCli.Printf("client about to send message to the server")

	err := a.client.Send(myString)

	if err != nil {
		a.lCli.Printf("clientSend : %s\n", err.Error())
		return errors.New(updateAppStatus(a, err.Error()))
	}

	a.lCli.Println("client assumes success of send :", myString)

	return nil
}

/*
	//	possible client read example - receiving messages back from the server
	mType, sbMsg, err := a.conClient.ReadMessage()
	if err != nil {
		return err
	}
	if mType == 234 {
		panic(err)
	}

	a.chToClient <- string(sbMsg)
*/

/*
	// client write examples
	go func() {
		for {
			_ = c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("first ws msg : %d", count)))
			time.Sleep(2000 * time.Millisecond)
		}
	}()
		//	=========================================
*/

/*
	// client reads off of the replies by the server
	go func() {
		for sMsg := range chanMessages {
			l.Print("reply from Server :", sMsg)
		}
	}()
	//	=========================================
*/

func main() {

	aFyne := app.New()
	w := aFyne.NewWindow("WebSocket SwissArmy")

	//	a = new(ewApp)

	a.client = new(xt995.WsClient)
	a.server = new(xt995.WsServer)

	a.lApp = loggerSetup("./a.log", "app:")
	a.lCli = loggerSetup("./c.log", "cli:")
	a.lSrv = loggerSetup("./s.log", "srv:")

	a.lApp.Print("loggers loaded, moving on to the core of the app logic")

	//	a.sServerAddr = "127.0.0.1:5757"
	//	a.sClientAddr = "127.0.0.1:5757"

	a.bServerUp = false

	// Buffered channels to account for bursts or spikes in data
	//	a.chToClient = make(chan string)

	a.lCli.Printf("msg channels established")

	//	UI setup

	a.displayText = &widget.Label{Text: ""}
	a.displayText.Alignment = fyne.TextAlignTrailing
	a.displayText.TextStyle.Monospace = true

	a.appStatus = &widget.Label{Text: ""}
	a.appStatus.Alignment = fyne.TextAlignTrailing
	a.appStatus.TextStyle.Monospace = true
	updateAppStatus(&a, "Normal")

	a.serverStatus = &widget.Label{Text: ""}
	a.serverStatus.TextStyle.Monospace = true

	updateServerStatus(&a)

	//	myApp.displayText.SetText(err.Error())

	//execGrid := fyne.NewContainerWithLayout(layout.NewGridLayout(1), newExecGrid(&myApp))

	//uiContainer := fyne.NewContainerWithLayout(layout.NewGridLayout(1), myApp.displayText, execGrid)

	clearGrid := fyne.NewContainerWithLayout(layout.NewGridLayout(1), newClearGrid(&a))
	//	buttons := fyne.NewContainerWithLayout(layout.NewGridLayout(2), clearGrid, newOperatorGrid(&myApp))
	uiContainer :=
		fyne.NewContainerWithLayout(
			layout.NewGridLayout(1),
			a.displayText,
			a.appStatus,
			a.serverStatus,
			clearGrid)

	//myApp.displayText.SetText("Error : " + err.Error())
	//	myApp.displayText.SetText("Syd")

	//	activate the UI

	w.SetContent(uiContainer)
	w.ShowAndRun()
}

func serverStopStart(a *ewApp, sCmd string) error {

	if sCmd != "stop" && sCmd != "start" {
		return errors.New(updateAppStatus(a, "only stop and start are supported : "+sCmd))
	}

	if sCmd == "stop" {
		a.lApp.Printf("ServerStopStart : Server stopping")

		err := a.server.Stop()

		if err != nil {
			return err
		}

		a.bServerUp = false
		updateServerStatus(a)

		a.lApp.Printf("ServerStopStart : Server stopped")

		return nil
	}

	err := a.server.Start()

	if err != nil {
		return err
	}

	a.bServerUp = true

	a.lApp.Printf("ServerStopStart : Server started")

	return nil
}

func newClearGrid(a *ewApp) *fyne.Container {
	var clears []fyne.CanvasObject

	clears = append(clears, &widget.Button{
		Text: "Start", OnTapped: func() { serverStopStart(a, "start") },
	})
	clears = append(clears, &widget.Button{
		Text: "Stop", OnTapped: func() { serverStopStart(a, "stop") },
	})
	clears = append(clears, &widget.Button{
		Text: "Connect", OnTapped: func() { clientConnect(a) },
	})
	clears = append(clears, &widget.Button{
		Text: "Send", OnTapped: func() { clientSend(a, "message from Syd!", 7) }, //add count here
	})
	clears = append(clears, &widget.Button{
		Text: "Quit", OnTapped: func() {
			a.lApp.Printf(" ")
			a.lApp.Printf("fyneControl : app shutdown")
			os.Exit(0)
		},
	})

	return fyne.NewContainerWithLayout(layout.NewGridLayout(2), clears...)
}

func checkError(a *ewApp, err error) {
	if err != nil {
		a.appStatus.SetText("Error: " + err.Error())
		fmt.Println("Error: ", err)
		a.lApp.Print("Error: ", err)
		os.Exit(35)
	}
}

func updateAppStatus(a *ewApp, sInStatus string) string {
	a.lApp.Print("AppStatus : %v \n", sInStatus)
	a.appStatus.SetText(sInStatus)
	return sInStatus
}

func updateServerStatus(a *ewApp) {

	if a.bServerUp {
		a.serverStatus.SetText("Up")
	} else {
		a.serverStatus.SetText("Down")
	}
}

/*
func newExecGrid(myApp *ewApp) *fyne.Container {
	var clears []fyne.CanvasObject

	clears = append(clears, &widget.Button{
		Text: "Quit", OnTapped: func() { os.Exit(0) },
	})

	//iconFile, _ := os.Open("images/usb.png")
	//r := bufio.NewReader(iconFile)
	//b, _ := ioutil.ReadAll(r)
	//clears = append(clears, &widget.Button{
	//	Icon: fyne.NewStaticResource("icon", b), OnTapped: func() { os.Exit(0) },
	//})

	return fyne.NewContainerWithLayout(layout.NewGridLayout(1), clears...)
}
*/
