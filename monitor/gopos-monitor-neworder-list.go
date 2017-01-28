package main

import (
	//	"encoding/json"
	//	"fmt"
	"fmt"
	"log"
	//	"net"
	//	"strconv"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var newOrderTreeView *gtk.TreeView
var newOrderListStore *gtk.ListStore

const (
	COLUMN_NEW_ORDER_DISH_ID = iota
	COLUMN_NEW_ORDER_DISH_NAME
	COLUMN_NEW_ORDER_DISH_PRICE
)

func createNewOrderTreeView() {
	var err error
	newOrderTreeView, err = gtk.TreeViewNew()

	if err != nil {
		log.Fatal("Unable to create tables tree view: ", err)
	}

	newOrderTreeView.AppendColumn(createColumn("ID", COLUMN_NEW_ORDER_DISH_ID))
	newOrderTreeView.AppendColumn(createColumn("Название блюда", COLUMN_NEW_ORDER_DISH_NAME))
	newOrderTreeView.AppendColumn(createColumn("Цена", COLUMN_NEW_ORDER_DISH_PRICE))

	newOrderListStore, err = gtk.ListStoreNew(glib.TYPE_INT, glib.TYPE_STRING, glib.TYPE_DOUBLE)
	if err != nil {
		log.Fatal("Unable to create tables list store: ", err)
	}

	newOrderTreeView.SetModel(newOrderListStore)
}

func newOrderAddRow(id int, name string, price float64) {
	iter := newOrderListStore.Append()
	fmt.Println(id, name, price)
	err := newOrderListStore.Set(iter, []int{COLUMN_NEW_ORDER_DISH_ID, COLUMN_NEW_ORDER_DISH_NAME,
		COLUMN_NEW_ORDER_DISH_PRICE}, []interface{}{id, name, price})

	if err != nil {
		log.Fatal("Unable to add tables row: ", err)
	}
}

func newOrderListCreateWindow() *gtk.Window {
	newOrderListWindow, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	newOrderListWindow.SetTitle("gopos-monitor-neworder-lsit")
	newOrderListWindow.SetDefaultSize(800, 600)

	newOrderVbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	if err != nil {
		log.Fatal("Unable to create new order vertical box: ", err)
	}

	createNewOrderTreeView()
	scrolledWindow, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		log.Fatalf("Error on creating workers scrolled window")
	}
	scrolledWindow.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	scrolledWindow.Add(newOrderTreeView)

	newOrderVbox.PackStart(scrolledWindow, true, true, 3)

	newOrderListWindow.Add(newOrderVbox)

	return newOrderListWindow
}
