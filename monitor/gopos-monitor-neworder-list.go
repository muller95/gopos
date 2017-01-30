package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var newOrderTreeView *gtk.TreeView
var newOrderListStore *gtk.ListStore
var newOrderPriceLabel *gtk.Label
var newOrderListWindow *gtk.Window

var newOrderPrice float64

const (
	COLUMN_NEW_ORDER_LIST_DISH_ID = iota
	COLUMN_NEW_ORDER_LIST_DISH_NAME
	COLUMN_NEW_ORDER_LIST_DISH_PRICE
)

func createNewOrderTreeView() {
	var err error
	newOrderTreeView, err = gtk.TreeViewNew()

	if err != nil {
		log.Fatal("Unable to create tables tree view: ", err)
	}

	newOrderTreeView.AppendColumn(createColumn("ID", COLUMN_NEW_ORDER_LIST_DISH_ID))
	newOrderTreeView.AppendColumn(createColumn("Название блюда", COLUMN_NEW_ORDER_LIST_DISH_NAME))
	newOrderTreeView.AppendColumn(createColumn("Цена", COLUMN_NEW_ORDER_LIST_DISH_PRICE))

	newOrderListStore, err = gtk.ListStoreNew(glib.TYPE_INT, glib.TYPE_STRING, glib.TYPE_DOUBLE)
	if err != nil {
		log.Fatal("Unable to create tables list store: ", err)
	}

	newOrderTreeView.SetModel(newOrderListStore)
}

func newOrderAddRow(id int, name string, price float64) {
	iter := newOrderListStore.Append()

	err := newOrderListStore.Set(iter, []int{COLUMN_NEW_ORDER_LIST_DISH_ID, COLUMN_NEW_ORDER_LIST_DISH_NAME,
		COLUMN_NEW_ORDER_LIST_DISH_PRICE}, []interface{}{id, name, price})

	if err != nil {
		log.Fatal("Unable to add tables row: ", err)
	}
}

func dishNewOrderDeleteSelectedButtonClicked() {
	selection, err := newOrderTreeView.GetSelection()
	if err != nil {
		log.Fatal("Error on getting new order selection")
	}

	rows := selection.GetSelectedRows(newOrderListStore)
	if rows == nil {
		return
	}
	path := rows.Data().(*gtk.TreePath)
	iter, err := newOrderListStore.GetIter(path)
	if err != nil {
		log.Fatal("Error on getting iter: ", err)
	}
	value, err := newOrderListStore.GetValue(iter, 2)
	if err != nil {
		log.Fatal("Error on getting value: ", err)
	}
	newOrderPrice -= value.GetDouble()
	newOrderPriceLabel.SetText(fmt.Sprintf("Цена: %.2f", newOrderPrice))
	newOrderListStore.Remove(iter)
}

func newOrderConfirmButtonClicked() {
	newOrderMap := make(map[int]int)

	iter, isNotEmpty := newOrderListStore.GetIterFirst()
	if !isNotEmpty {
		return
	}

	for {
		value, err := newOrderListStore.GetValue(iter, COLUMN_NEW_ORDER_LIST_DISH_ID)
		if err != nil {
			log.Fatal("Error on getting value: ", err)
		}
		dishId := value.GetInt()
		newOrderMap[dishId]++
		if !newOrderListStore.IterNext(iter) {
			break
		}
	}

	orderString := ""
	for id, count := range newOrderMap {
		orderString += fmt.Sprintf(" %d:%d", id, count)
	}
	orderString = strings.Trim(orderString, " ")
	fmt.Println(newOrderMap)
	fmt.Println(orderString)

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", goposServerIp,
		goposServerPort))

	if err != nil {
		log.Fatal("Unable to connect to server")
	}

	requestMap := make(map[string]string)
	requestMap["group"] = "ORDER"
	requestMap["action"] = "CREATE"
	requestMap["password"] = goposServerPassword
	requestMap["order_string"] = orderString
	requestMap["table_number"] = fmt.Sprintf("%d", freeTableNumber)
	requestMap["price"] = fmt.Sprintf("%f", newOrderPrice)
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
		conn.Close()
		return
	} else {
		newOrderWindow.Destroy()
		newOrderListWindow.Destroy()
	}
}

func newOrderListCreateWindow() *gtk.Window {
	var err error
	newOrderListWindow, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	newOrderListWindow.SetTitle("gopos-monitor-neworder-lsit")
	newOrderListWindow.SetDefaultSize(800, 600)

	newOrderPriceLabel, err = gtk.LabelNew("Цена: 0.0")
	if err != nil {
		log.Fatal("Error on creating price label")
	}

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

	dishDeleteButton, err := gtk.ButtonNewWithLabel("Удалить из заказа")
	if err != nil {
		log.Fatal("Unable to create delete button: ", err)
	}
	dishDeleteButton.Connect("clicked", dishNewOrderDeleteSelectedButtonClicked, nil)

	newOrderConfirm, err := gtk.ButtonNewWithLabel("Подтвердить заказ")
	if err != nil {
		log.Fatal("Unable to create comfirm button: ", err)
	}
	newOrderConfirm.Connect("clicked", newOrderConfirmButtonClicked)

	newOrderVbox.PackStart(scrolledWindow, true, true, 3)
	newOrderVbox.PackStart(newOrderPriceLabel, false, false, 3)
	newOrderVbox.PackStart(dishDeleteButton, false, true, 3)
	newOrderVbox.PackStart(newOrderConfirm, false, true, 3)

	newOrderListWindow.Add(newOrderVbox)

	return newOrderListWindow
}
