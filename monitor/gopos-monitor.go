//this package is gopos monitor module
package main

import (
	"log"
	"os"
	"github.com/gotk3/gotk3/gtk"
)

var mainWindow *gtk.Window
var goposServerIp, goposServerPassword, goposServerPort string

func main() {
	var err error

	goposServerIp = os.Getenv("GOPOS_SERVER_IP")
	if goposServerIp == "" {
		log.Fatal("GOPOS_SERVER_IP is not set")
	}

	goposServerPassword = os.Getenv("GOPOS_SERVER_PASSWORD")
	if goposServerPassword == "" {
		log.Fatal("GOPOS_SERVER_PASSWORD is not set")
	}

	goposServerPort = os.Getenv("GOPOS_SERVER_PORT")
	if goposServerPort == "" {
		log.Fatal("GOPOS_SERVER_PORT is not set")
	}

	// Initialize GTK without parsing any command line arguments.
	gtk.Init(nil)

	// Create a new toplevel window, set its title, and connect it to the
	// "destroy" signal to exit the GTK main loop when it is destroyed.
	mainWindow, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}

	mainWindow.SetTitle("gopos-monitor")
	mainWindow.Connect("destroy", func() {
		gtk.MainQuit()
	})

	//creates notebook
	notebook, err := gtk.NotebookNew()
	if err != nil {
		log.Fatal("Error on creating notebook: ", err)
	}

	freeTablesPageLabel, err := gtk.LabelNew("Свободные столики")
	if err != nil {
		log.Fatal("Error on creating new order page label: ", err)
	}

	notebook.AppendPage(freeTablesCreatePage(), freeTablesPageLabel)
	mainWindow.Add(notebook)

	getFreeTables()

	// Set the default window size.
	mainWindow.SetDefaultSize(800, 600)

	// Recursively show all widgets contained in this window.
	mainWindow.ShowAll()

	// Begin executing the GTK main loop.  This blocks until
	// gtk.MainQuit() is run.
	gtk.Main()
}
