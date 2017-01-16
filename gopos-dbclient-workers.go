//main package for gopos database client
//this file contains functions for workers tree view
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"net"
	"strconv"
	"time"
)

var workersTreeView *gtk.TreeView
var workersListStore *gtk.ListStore

const (
	COLUMN_WORKERS_ID = iota
	COLUMN_WORKERS_NAME
	COLUMN_WORKERS_DATE
)

func createWorkersTreeView() {
	var err error
	workersTreeView, err = gtk.TreeViewNew()

	if err != nil {
		log.Fatal("Unable to create tree view: ", err)
	}

	workersTreeView.AppendColumn(createColumn("ID", COLUMN_WORKERS_ID))
	workersTreeView.AppendColumn(createColumn("Имя", COLUMN_WORKERS_NAME))
	workersTreeView.AppendColumn(createColumn("Дата", COLUMN_WORKERS_DATE))

	workersListStore, err = gtk.ListStoreNew(glib.TYPE_INT, glib.TYPE_STRING,
		glib.TYPE_STRING)
	if err != nil {
		log.Fatal("Unable to create list store: ", err)
	}

	workersTreeView.SetModel(workersListStore)
}

func workerAddRow(id int, name string, date string) {
	iter := workersListStore.Append()

	err := workersListStore.Set(iter, []int{ COLUMN_WORKERS_ID, COLUMN_WORKERS_NAME,
		COLUMN_WORKERS_DATE },
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
			workerAddRow(id, workerName, timeString)
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
		workerAddRow(id, responseMap["name"], responseMap["date"])
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

func workersCreatePage() *gtk.Box {
	//creates workers tabpage
	workersVbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	if err != nil {
		log.Fatal("Unable to create main vertical box: ", err)
	}

	createWorkersTreeView()
	scrolledWindow, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		log.Fatalf("Error on creating workers scrolled window")
	}
	scrolledWindow.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	scrolledWindow.Add(workersTreeView)

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


	workersVbox.PackStart(scrolledWindow, true, true, 3)
	workersVbox.PackStart(workersFormHbox, false, false, 3)

	return workersVbox
}
