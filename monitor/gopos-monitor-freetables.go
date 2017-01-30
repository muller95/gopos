package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var freeTablesTreeView *gtk.TreeView
var freeTablesListStore *gtk.ListStore

var freeTableNumber int

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

	err := freeTablesListStore.Set(iter, []int{COLUMN_FREE_TABLES_NUMBER},
		[]interface{}{number})

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

func orderCreateButtonClicked(btn *gtk.Button) {
	selection, err := freeTablesTreeView.GetSelection()
	if err != nil {
		log.Fatal("Error on getting new order selection")
	}

	rows := selection.GetSelectedRows(freeTablesListStore)
	if rows == nil {
		return
	}
	path := rows.Data().(*gtk.TreePath)
	iter, err := freeTablesListStore.GetIter(path)
	if err != nil {
		log.Fatal("Error on getting iter: ", err)
	}
	value, err := freeTablesListStore.GetValue(iter, COLUMN_FREE_TABLES_NUMBER)
	if err != nil {
		log.Fatal("Error on getting value: ", err)
	}
	freeTableNumber = value.GetInt()

	btn.SetSensitive(false)

	newOrderWindow := newOrderCreateWindow()
	getNewOrderCategories()

	newOrderListWindow := newOrderListCreateWindow()

	newOrderWindow.Connect("destroy", func(window *gtk.Window) {
		newOrderPrice = 0.0
		freeTableNumber = 0.0
		freeTablesListStore.Clear()
		getFreeTables()
		orderedTablesListStore.Clear()
		getOrderedTables()
		newOrderListWindow.Destroy()
		btn.SetSensitive(true)
	})
	newOrderListWindow.Connect("destroy", func(window *gtk.Window) {
		newOrderPrice = 0.0
		freeTableNumber = 0.0
		freeTablesListStore.Clear()
		getFreeTables()
		orderedTablesListStore.Clear()
		getOrderedTables()
		btn.SetSensitive(true)
		newOrderWindow.Destroy()
	})

	newOrderWindow.ShowAll()
	newOrderListWindow.ShowAll()
}

func updateTablesButtonClicked() {
	freeTablesListStore.Clear()
	getFreeTables()
	orderedTablesListStore.Clear()
	getOrderedTables()
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

	orderCreateButton, err := gtk.ButtonNewWithLabel("Создать заказ")
	if err != nil {
		log.Fatal("Unable to create create order button: ", err)
	}
	orderCreateButton.Connect("clicked", orderCreateButtonClicked, nil)

	updateTablesButton, err := gtk.ButtonNewWithLabel("Обновить данные")
	if err != nil {
		log.Fatal("Unable to create create order button: ", err)
	}
	updateTablesButton.Connect("clicked", updateTablesButtonClicked, nil)

	freeTablesVbox.PackStart(scrolledWindow, true, true, 3)
	freeTablesVbox.PackStart(orderCreateButton, false, true, 3)
	freeTablesVbox.PackStart(updateTablesButton, false, true, 3)

	return freeTablesVbox
}
