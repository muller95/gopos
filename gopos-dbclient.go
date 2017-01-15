package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"net"
	"os"
	"strconv"
	"time"
)

var workersTreeView *gtk.TreeView
var workersListStore *gtk.ListStore
var mainWindow *gtk.Window
var goposServerIp, goposServerPassword, goposServerPort string

const (
	COLUMN_ID = iota
	COLUMN_NAME
	COLUMN_DATE
)

func createColumn(title string, id int) *gtk.TreeViewColumn {
	cellRenderer, err := gtk.CellRendererTextNew()
	if err != nil {
		log.Fatal("Unable to create text cell renderer: ", err)
	}

	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "text", id)
	if err != nil {
		log.Fatal("Unable to create cell column")
	}

	return column
}

func createTreeView() (*gtk.TreeView, *gtk.ListStore) {
	treeView, err := gtk.TreeViewNew()

	if err != nil {
		log.Fatal("Unable to create tree view: ", err)
	}

	treeView.AppendColumn(createColumn("ID", COLUMN_ID))
	treeView.AppendColumn(createColumn("Имя", COLUMN_NAME))
	treeView.AppendColumn(createColumn("Дата", COLUMN_DATE))

	listStore, err := gtk.ListStoreNew(glib.TYPE_INT, glib.TYPE_STRING, glib.TYPE_STRING)
	if err != nil {
		log.Fatal("Unable to create list store: ", err)
	}

	treeView.SetModel(listStore)

	return treeView, listStore
}

func workerAddRow(listStore *gtk.ListStore, id int, name string, date string) {
	iter := listStore.Append()

	err := listStore.Set(iter, []int{ COLUMN_ID, COLUMN_NAME, COLUMN_DATE },
		[]interface{} { id, name, date })

	if err != nil {
		log.Fatal("Unable to add row: ", err)
	}
}

func workerAddButtonClicked(btn *gtk.Button, workerNameEntry *gtk.Entry) {
	workerName, err := workerNameEntry.GetText()

	if err == nil {
		if len(workerName) > 0 {
			currTime := time.Now().Local()
			timeString := fmt.Sprintf("%d-%d-%d", currTime.Year(), currTime.Month(),
				currTime.Day())


			conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", goposServerIp,
				goposServerPort))

			if err != nil {
				log.Fatal("Unable to connect to server")
			}

			requestMap := make(map[string]string)
			requestMap["group"] = "WORKER"
			requestMap["action"] = "ADD"
			requestMap["password"] = goposServerPassword
			requestMap["name"] = workerName
			requestMap["date"] = timeString
			encoder := json.NewEncoder(conn)
			err = encoder.Encode(requestMap)
			if err != nil {
				log.Fatal("Error on encode request map: ", requestMap)
			}

			decoder := json.NewDecoder(conn)
			responseMap := make(map[string]string)
			err = decoder.Decode(&responseMap)
			if err != nil {
				log.Fatal("Error on decoding response: ", err)
			}

			id, err := strconv.Atoi(responseMap["id"])
			if err != nil {
				log.Fatal("Error on converting id to int: ", err)
			}
			if id < 0 {
				messageDialog := gtk.MessageDialogNew(mainWindow,
					gtk.DIALOG_MODAL, gtk.MESSAGE_WARNING, gtk.BUTTONS_OK,
					responseMap["error"])
				messageDialog.Run()
				messageDialog.Destroy()
				return
			}
			workerAddRow(workersListStore, id, workerName, timeString)
		} else {
			messageDialog := gtk.MessageDialogNew(mainWindow, gtk.DIALOG_MODAL,
				gtk.MESSAGE_WARNING, gtk.BUTTONS_OK, "Введите имя работника")
			messageDialog.Run()
			messageDialog.Destroy()
		}
	} else {
		log.Fatal("Unable to get worker name entry text: ", err)
	}
}

func getWorkers() {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", goposServerIp,
		goposServerPort))

	if err != nil {
		log.Fatal("Unable to connect to server")
	}

	requestMap := make(map[string]string)
	requestMap["group"] = "WORKER"
	requestMap["action"] = "GET"
	requestMap["password"] = goposServerPassword
	encoder := json.NewEncoder(conn)
	err = encoder.Encode(requestMap)
	if err != nil {
		log.Fatal("Error on encode request map: ", requestMap)
	}

	decoder := json.NewDecoder(conn)
	for {
		responseMap := make(map[string]string)
		err = decoder.Decode(&responseMap)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Error on decoding response: ", err)
		}
		id, err := strconv.Atoi(responseMap["id"])
		if err != nil {
			log.Fatal("Error on converting id to int: ", err)
		}
		workerAddRow(workersListStore, id, responseMap["name"], responseMap["date"])
	}
}

func workerDeleteSelectedButtonClicked() {
	selection, err := workersTreeView.GetSelection()
	if err != nil {
		log.Fatal("Error on getting workers selection")
	}

	rows := selection.GetSelectedRows(workersListStore)
	path := rows.Data().(*gtk.TreePath)
	iter, err := workersListStore.GetIter(path)
	if err != nil {
		log.Fatal("Error on getting iter: ", err)
	}
	value, err := workersListStore.GetValue(iter, 0)
	if err != nil {
		log.Fatal("Error on getting value: ", err)
	}
	id := value.GetInt()

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", goposServerIp,
		goposServerPort))

	if err != nil {
		log.Fatal("Unable to connect to server")
	}

	requestMap := make(map[string]string)
	requestMap["group"] = "WORKER"
	requestMap["action"] = "DELETE"
	requestMap["password"] = goposServerPassword
	requestMap["id"] = fmt.Sprintf("%d", id)
	encoder := json.NewEncoder(conn)
	err = encoder.Encode(requestMap)
	if err != nil {
		log.Fatal("Error on encode request map: ", requestMap)
	}

	decoder := json.NewDecoder(conn)
	responseMap := make(map[string]string)
	err = decoder.Decode(&responseMap)
	if err != nil {
		log.Fatal("Error on decoding response: ", err)
	}
	if responseMap["result"] == "OK" {
		workersListStore.Remove(iter)
	}
}

func main() {
	var err error

	goposServerIp = os.Getenv("GOPOS_SERVER_IP")
	if goposServerIp == "" {
		log.Fatal("GOPOS_SERVER_IP is not set")
	}
	fmt.Println(goposServerIp)

	goposServerPassword = os.Getenv("GOPOS_SERVER_PASSWORD")
	if goposServerPassword == "" {
		log.Fatal("GOPOS_SERVER_PASSWORD is not set")
	}
	fmt.Println(goposServerPassword)

	goposServerPort = os.Getenv("GOPOS_SERVER_PORT")
	if goposServerPort == "" {
		log.Fatal("GOPOS_SERVER_PORT is not set")
	}
	fmt.Println(goposServerPort)

	// Initialize GTK without parsing any command line arguments.
	gtk.Init(nil)


	// Create a new toplevel window, set its title, and connect it to the
	// "destroy" signal to exit the GTK main loop when it is destroyed.
	mainWindow, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}

	mainWindow.SetTitle("gopos-dbclient")
	mainWindow.Connect("destroy", func() {
		gtk.MainQuit()
	})


	workersVbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	if err != nil {
		log.Fatal("Unable to create main vertical box: ", err)
	}


	workersTreeView, workersListStore = createTreeView()

	workersFormHbox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	if err != nil {
		log.Fatal("Unable to create workers form horizontal box: ", err)
	}

	workerNameLabel, err := gtk.LabelNew("Имя работника:")
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}

	workerNameEntry, err := gtk.EntryNew()
	if err != nil {
		log.Fatal("Unable to create entry: ", err)
	}

	workerAddButton, err := gtk.ButtonNewWithLabel("Добавить")
	if err != nil {
		log.Fatal("Unable to create add button: ", err)
	}
	workerAddButton.Connect("clicked", workerAddButtonClicked, workerNameEntry)

	workerDeleteSelectedButton, err := gtk.ButtonNewWithLabel("Удалить выбранного")
	if err != nil {
		log.Fatal("Unable to create add button: ", err)
	}
	workerDeleteSelectedButton.Connect("clicked", workerDeleteSelectedButtonClicked, nil)

	workersFormHbox.PackStart(workerNameLabel, false, false, 3)
	workersFormHbox.PackStart(workerNameEntry, true, true, 3)
	workersFormHbox.PackStart(workerAddButton, true, true, 3)
	workersFormHbox.PackStart(workerDeleteSelectedButton, true, true, 3)


	workersVbox.PackStart(workersTreeView, true, true, 3)
	workersVbox.PackStart(workersFormHbox, false, false, 3)
	mainWindow.Add(workersVbox)

	getWorkers()

	// Set the default window size.
	mainWindow.SetDefaultSize(800, 600)

	// Recursively show all widgets contained in this window.
	mainWindow.ShowAll()

	// Begin executing the GTK main loop.  This blocks until
	// gtk.MainQuit() is run.
	gtk.Main()
}

