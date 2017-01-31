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

const (
	COLUMN_ORDERED_TABLES_NUMBER = iota
)

var orderedTablesTreeView *gtk.TreeView
var orderedTablesListStore *gtk.ListStore

var orderedTableNumber int

func createOrderedTablesTreeView() {
	var err error
	orderedTablesTreeView, err = gtk.TreeViewNew()

	if err != nil {
		log.Fatal("Unable to create tables tree view: ", err)
	}

	orderedTablesTreeView.AppendColumn(createColumn("Номер", COLUMN_ORDERED_TABLES_NUMBER))

	orderedTablesListStore, err = gtk.ListStoreNew(glib.TYPE_INT)
	if err != nil {
		log.Fatal("Unable to create tables list store: ", err)
	}

	orderedTablesTreeView.SetModel(orderedTablesListStore)
}

func orderedTableAddRow(number int) {
	iter := orderedTablesListStore.Append()

	err := orderedTablesListStore.Set(iter, []int{COLUMN_ORDERED_TABLES_NUMBER},
		[]interface{}{number})

	if err != nil {
		log.Fatal("Unable to add tables row: ", err)
	}
}

func getOrderedTables() {
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
		if responseMap["current_order"] == "-1" {
			continue
		}

		number, _ := strconv.Atoi(responseMap["number"])
		orderedTableAddRow(number)
	}
}

func orderEditButtonClicked(btn *gtk.Button) {
	selection, err := orderedTablesTreeView.GetSelection()
	if err != nil {
		log.Fatal("Error on getting new order selection")
	}

	rows := selection.GetSelectedRows(orderedTablesListStore)
	if rows == nil {
		return
	}
	path := rows.Data().(*gtk.TreePath)
	iter, err := orderedTablesListStore.GetIter(path)
	if err != nil {
		log.Fatal("Error on getting iter: ", err)
	}
	value, err := orderedTablesListStore.GetValue(iter, COLUMN_ORDERED_TABLES_NUMBER)
	if err != nil {
		log.Fatal("Error on getting value: ", err)
	}
	orderedTableNumber = value.GetInt()

	btn.SetSensitive(false)

	existingOrderWindow := existingOrderCreateWindow()
	getExistingOrderCategories()

	existingOrderListWindow := existingOrderListCreateWindow()
	getOrder()

	existingOrderWindow.Connect("destroy", func(window *gtk.Window) {
		existingOrderPrice = 0.0
		orderedTableNumber = 0
		freeTablesListStore.Clear()
		getFreeTables()
		orderedTablesListStore.Clear()
		getOrderedTables()
		existingOrderListWindow.Destroy()
		btn.SetSensitive(true)
	})
	existingOrderListWindow.Connect("destroy", func(window *gtk.Window) {
		existingOrderPrice = 0.0
		orderedTableNumber = 0
		freeTablesListStore.Clear()
		getFreeTables()
		orderedTablesListStore.Clear()
		getOrderedTables()
		btn.SetSensitive(true)
		existingOrderWindow.Destroy()
	})

	existingOrderWindow.ShowAll()
	existingOrderListWindow.ShowAll()
}

func orderedTablesCreatePage() *gtk.Box {
	//creates tables tabpage
	orderedTablesVbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	if err != nil {
		log.Fatal("Unable to create main vertical box: ", err)
	}

	createOrderedTablesTreeView()
	scrolledWindow, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		log.Fatalf("Error on creating workers scrolled window")
	}
	scrolledWindow.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	scrolledWindow.Add(orderedTablesTreeView)

	orderEditButton, err := gtk.ButtonNewWithLabel("Закрыть/Изменить заказ")
	if err != nil {
		log.Fatal("Unable to create edit order button: ", err)
	}
	orderEditButton.Connect("clicked", orderEditButtonClicked, nil)

	orderedTablesVbox.PackStart(scrolledWindow, true, true, 3)
	orderedTablesVbox.PackStart(orderEditButton, false, true, 3)

	return orderedTablesVbox
}
