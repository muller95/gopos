//main package for gopos database client
//this file contains functions for tables tree view
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
)

var freeTablesTreeView *gtk.TreeView
var freeTablesListStore *gtk.ListStore

const (
	COLUMN_FREE_TABLES_NUMBER = iota
)

func createFreeTablesTreeView() {
	var err error
	freeTablesTreeView, err = gtk.TreeViewNew()

	if err != nil {
		log.Fatal("Unable to create tables tree view: ", err)
	}

	freeTablesTreeView.AppendColumn(createColumn("Номер", COLUMN_FREE_TABLES_NUMBER))

	freeTablesListStore, err = gtk.ListStoreNew(glib.TYPE_INT)
	if err != nil {
		log.Fatal("Unable to create tables list store: ", err)
	}

	freeTablesTreeView.SetModel(freeTablesListStore)
}

func freeTableAddRow(number int) {
	iter := freeTablesListStore.Append()

	err := freeTablesListStore.Set(iter, []int{ COLUMN_FREE_TABLES_NUMBER },
		[]interface{} { number })

	if err != nil {
		log.Fatal("Unable to add tables row: ", err)
	}
}

func getFreeTables() {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", goposServerIp,
		goposServerPort))

	if err != nil {
		log.Fatal("Unable to connect to server")
	}

	requestMap := make(map[string]string)
	requestMap["group"] = "TABLE"
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
		if responseMap["current_order"] != "-1" {
			continue
		}

		number, _ := strconv.Atoi(responseMap["number"])
		freeTableAddRow(number)
	}
}


func freeTablesCreatePage() *gtk.Box {
	//creates tables tabpage
	freeTablesVbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	if err != nil {
		log.Fatal("Unable to create main vertical box: ", err)
	}


	createFreeTablesTreeView()
	scrolledWindow, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		log.Fatalf("Error on creating workers scrolled window")
	}
	scrolledWindow.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	scrolledWindow.Add(freeTablesTreeView)

/*	freeTablesFormHbox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	if err != nil {
		log.Fatal("Unable to create tables form horizontal box: ", err)
	}*/

	orderCreateButton, err := gtk.ButtonNewWithLabel("Создать заказ")
	if err != nil {
		log.Fatal("Unable to create create order button: ", err)
	}
	//tableAddButton.Connect("clicked", tableAddButtonClicked, tableNumberEntry)

//	freeTablesFormHbox.PackStart(orderCreateButton, false, true, 3)

	freeTablesVbox.PackStart(scrolledWindow, true, true, 3)
	freeTablesVbox.PackStart(orderCreateButton, false, true, 3)

	return freeTablesVbox
}
