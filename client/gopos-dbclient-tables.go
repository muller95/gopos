//main package for gopos database client
//this file contains functions for tables tree view
package main

import (
	"encoding/json"
	"fmt"
//	"io"
	"log"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"net"
	"strconv"
)

var tablesTreeView *gtk.TreeView
var tablesListStore *gtk.ListStore

const (
	COLUMN_TABLES_NUMBER = iota
)

func createTablesTreeView() {
	var err error
	tablesTreeView, err = gtk.TreeViewNew()

	if err != nil {
		log.Fatal("Unable to create tables tree view: ", err)
	}

	tablesTreeView.AppendColumn(createColumn("Номер", COLUMN_TABLES_NUMBER))

	tablesListStore, err = gtk.ListStoreNew(glib.TYPE_INT)
	if err != nil {
		log.Fatal("Unable to create tables list store: ", err)
	}

	tablesTreeView.SetModel(tablesListStore)
}

func tableAddRow(number int) {
	iter := tablesListStore.Append()

	err := tablesListStore.Set(iter, []int{ COLUMN_TABLES_NUMBER },
		[]interface{} { number })

	if err != nil {
		log.Fatal("Unable to add tables row: ", err)
	}
}


func tableAddButtonClicked(btn *gtk.Button, tableNumberEntry *gtk.Entry) {
	tableNumber, err := tableNumberEntry.GetText()

	if err == nil {
		if len(tableNumber) > 0 {
			number, err := strconv.Atoi(tableNumber)
			if err != nil || number < 0{
				messageDialog := gtk.MessageDialogNew(mainWindow,
					gtk.DIALOG_MODAL, gtk.MESSAGE_WARNING, gtk.BUTTONS_OK,
					"Номер столика должен быть целым неотрицательным числом")
				messageDialog.Run()
				messageDialog.Destroy()
				return
			}

			conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", goposServerIp,
				goposServerPort))

			if err != nil {
				log.Fatal("Unable to connect to server")
			}

			requestMap := make(map[string]string)
			requestMap["group"] = "TABLES"
			requestMap["action"] = "ADD"
			requestMap["password"] = goposServerPassword
			requestMap["number"] = tableNumber
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

			if responseMap["result"] != "OK" {
				messageDialog := gtk.MessageDialogNew(mainWindow,
					gtk.DIALOG_MODAL, gtk.MESSAGE_WARNING, gtk.BUTTONS_OK,
					responseMap["error"])
				messageDialog.Run()
				messageDialog.Destroy()
				return
			}

			tableAddRow(number)
		} else {
			messageDialog := gtk.MessageDialogNew(mainWindow, gtk.DIALOG_MODAL,
				gtk.MESSAGE_WARNING, gtk.BUTTONS_OK, "Введите номер столика")
			messageDialog.Run()
			messageDialog.Destroy()
		}
	} else {
		log.Fatal("Unable to get worker name entry text: ", err)
	}
}


func tablesCreatePage() *gtk.Box {
	//creates tables tabpage
	tablesVbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	if err != nil {
		log.Fatal("Unable to create main vertical box: ", err)
	}


	createTablesTreeView()
	scrolledWindow, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		log.Fatalf("Error on creating workers scrolled window")
	}
	scrolledWindow.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	scrolledWindow.Add(tablesTreeView)

	tablesFormHbox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	if err != nil {
		log.Fatal("Unable to create tables form horizontal box: ", err)
	}

	tableNumberLabel, err := gtk.LabelNew("Номер столика:")
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}

	tableNumberEntry, err := gtk.EntryNew()
	if err != nil {
		log.Fatal("Unable to create entry: ", err)
	}

	tableAddButton, err := gtk.ButtonNewWithLabel("Добавить")
	if err != nil {
		log.Fatal("Unable to create add button: ", err)
	}
	tableAddButton.Connect("clicked", tableAddButtonClicked, tableNumberEntry)

	tableDeleteSelectedButton, err := gtk.ButtonNewWithLabel("Удалить выбранного")
	if err != nil {
		log.Fatal("Unable to create add button: ", err)
	}
//	tableDeleteSelectedButton.Connect("clicked", tableDeleteSelectedButtonClicked, nil)

	tablesFormHbox.PackStart(tableNumberLabel, false, false, 3)
	tablesFormHbox.PackStart(tableNumberEntry, true, true, 3)
	tablesFormHbox.PackStart(tableAddButton, true, true, 3)
	tablesFormHbox.PackStart(tableDeleteSelectedButton, true, true, 3)


	tablesVbox.PackStart(scrolledWindow, true, true, 3)
	tablesVbox.PackStart(tablesFormHbox, false, false, 3)

	return tablesVbox
}
