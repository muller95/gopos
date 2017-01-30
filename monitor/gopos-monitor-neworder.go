//tsi package is gopos monitor module
//this file contains functions for menu tree view
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
	COLUMN_NEW_ORDER_CATEGORIES_ID = iota
	COLUMN_NEW_ORDER_CATEGORIES_NAME
)

const (
	COLUMN_NEW_ORDER_DISH_ID = iota
	COLUMN_NEW_ORDER_DISH_NAME
	COLUMN_NEW_ORDER_DISH_PRICE
)

var newOrderWindow *gtk.Window

var newOrderCategoriesTreeView *gtk.TreeView
var newOrderCategoriesListStore *gtk.ListStore

var newOrderDishesTreeView *gtk.TreeView
var newOrderDishesListStore *gtk.ListStore

func createNewOrderCategoriesTreeView() {
	var err error
	newOrderCategoriesTreeView, err = gtk.TreeViewNew()

	if err != nil {
		log.Fatal("Unable to create cateogories tree view: ", err)
	}

	newOrderCategoriesTreeView.AppendColumn(createColumn("ID", COLUMN_NEW_ORDER_CATEGORIES_ID))
	newOrderCategoriesTreeView.AppendColumn(createColumn("Название",
		COLUMN_NEW_ORDER_CATEGORIES_NAME))

	newOrderCategoriesListStore, err = gtk.ListStoreNew(glib.TYPE_INT, glib.TYPE_STRING)
	if err != nil {
		log.Fatal("Unable to create newOrderCategories store: ", err)
	}

	newOrderCategoriesTreeView.SetModel(newOrderCategoriesListStore)
}

func newOrderCategoryAddRow(id int, name string) {
	iter := newOrderCategoriesListStore.Append()

	err := newOrderCategoriesListStore.Set(iter, []int{COLUMN_NEW_ORDER_CATEGORIES_ID,
		COLUMN_NEW_ORDER_CATEGORIES_NAME}, []interface{}{id, name})

	if err != nil {
		log.Fatal("Unable to add newOrderCategories row: ", err)
	}
}

func createNewOrderDishesTreeView() {
	var err error
	newOrderDishesTreeView, err = gtk.TreeViewNew()

	if err != nil {
		log.Fatal("Unable to create newOrderDishes tree view: ", err)
	}

	newOrderDishesTreeView.AppendColumn(createColumn("ID", COLUMN_NEW_ORDER_DISH_ID))
	newOrderDishesTreeView.AppendColumn(createColumn("Название", COLUMN_NEW_ORDER_DISH_NAME))
	newOrderDishesTreeView.AppendColumn(createColumn("Цена", COLUMN_NEW_ORDER_DISH_PRICE))

	newOrderDishesListStore, err = gtk.ListStoreNew(glib.TYPE_INT, glib.TYPE_STRING,
		glib.TYPE_DOUBLE)
	if err != nil {
		log.Fatal("Unable to create newOrderDishes store: ", err)
	}

	newOrderDishesTreeView.SetModel(newOrderDishesListStore)
}

func newOrderDishAddRow(id int, name string, price float64) {
	iter := newOrderDishesListStore.Append()

	err := newOrderDishesListStore.Set(iter, []int{COLUMN_NEW_ORDER_DISH_ID, COLUMN_NEW_ORDER_DISH_NAME,
		COLUMN_NEW_ORDER_DISH_PRICE}, []interface{}{id, name, price})

	if err != nil {
		log.Fatal("Unable to add newOrderDish row: ", err)
	}
}

func getNewOrderCategories() {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", goposServerIp,
		goposServerPort))

	if err != nil {
		log.Fatal("Unable to connect to server")
	}

	requestMap := make(map[string]string)
	requestMap["group"] = "CATEGORY"
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

		id, _ := strconv.Atoi(responseMap["id"])
		newOrderCategoryAddRow(id, responseMap["name"])
	}
	conn.Close()
}

func newOrderCategoriesSelectionChanged(selection *gtk.TreeSelection) {
	newOrderDishesListStore.Clear()
	selection, err := newOrderCategoriesTreeView.GetSelection()
	if err != nil {
		log.Fatal("Error on getting newOrderCategories selection")
	}
	rows := selection.GetSelectedRows(newOrderCategoriesListStore)
	if rows == nil {
		return
	}

	path := rows.Data().(*gtk.TreePath)
	iter, err := newOrderCategoriesListStore.GetIter(path)

	value, err := newOrderCategoriesListStore.GetValue(iter, 0)
	if err != nil {
		log.Fatal("Error on getting value: ", err)
	}
	newOrderCategoryId := value.GetInt()

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", goposServerIp,
		goposServerPort))

	if err != nil {
		log.Fatal("Unable to connect to server")
	}

	requestMap := make(map[string]string)
	requestMap["group"] = "DISH"
	requestMap["action"] = "GET"
	requestMap["password"] = goposServerPassword
	requestMap["category_id"] = fmt.Sprintf("%d", newOrderCategoryId)
	encoder := json.NewEncoder(conn)
	err = encoder.Encode(requestMap)
	if err != nil {
		log.Fatal("Error on encode request map: ", requestMap)
	}

	decoder := json.NewDecoder(conn)
	for {
		responseMap := make(map[string]string)
		err = decoder.Decode(&responseMap)
		fmt.Printf("%v %v\n", responseMap, err)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Error on decoding response: ", err)
		}

		id, _ := strconv.Atoi(responseMap["id"])
		price, _ := strconv.ParseFloat(responseMap["price"], 64)
		newOrderDishAddRow(id, responseMap["name"], price)
	}
}

func newOrderDishAddButtonClicked(btn *gtk.Button) {
	errMessage := ""
	selection, err := newOrderDishesTreeView.GetSelection()
	if err != nil {
		log.Fatal("Error on getting newOrderCategories selection")
	}
	rows := selection.GetSelectedRows(newOrderDishesListStore)
	if rows == nil {
		errMessage += "Выберите блюда для добавления в заказ."
	}

	if errMessage != "" {
		messageDialog := gtk.MessageDialogNew(mainWindow, gtk.DIALOG_MODAL,
			gtk.MESSAGE_WARNING, gtk.BUTTONS_OK, errMessage)
		messageDialog.Run()
		messageDialog.Destroy()
		return
	}

	path := rows.Data().(*gtk.TreePath)
	iter, err := newOrderDishesListStore.GetIter(path)
	if err != nil {
		log.Fatal("Error on getting iter: ", err)
	}

	value, err := newOrderDishesListStore.GetValue(iter, 0)
	if err != nil {
		log.Fatal("Error on getting value: ", err)
	}
	newOrderDishId := value.GetInt()

	value, err = newOrderDishesListStore.GetValue(iter, 1)
	if err != nil {
		log.Fatal("Error on getting value: ", err)
	}
	newOrderDishName, err := value.GetString()
	if err != nil {
		log.Fatal("Error on getting string: ", err)
	}

	value, err = newOrderDishesListStore.GetValue(iter, 2)
	if err != nil {
		log.Fatal("Error on getting value: ", err)
	}
	newOrderDishPrice := value.GetDouble()
	newOrderPrice += newOrderDishPrice
	newOrderPriceLabel.SetText(fmt.Sprintf("Цена: %.2f", newOrderPrice))
	newOrderAddRow(newOrderDishId, newOrderDishName, newOrderDishPrice)
}

func newOrderCreateWindow() *gtk.Window {
	//creates menu tabpage
	var err error
	newOrderWindow, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	newOrderWindow.SetTitle("gopos-monitor-neworder")
	newOrderWindow.SetDefaultSize(800, 600)

	newOrderVbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	if err != nil {
		log.Fatal("Unable to create main vertica; box: ", err)
	}

	menuHbox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	if err != nil {
		log.Fatal("Unable to create menu horizontal box: ", err)
	}

	newOrderCategoriesVbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	if err != nil {
		log.Fatal("Unable to create newOrderCategories vertical box: ", err)
	}

	newOrderDishesVbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	if err != nil {
		log.Fatal("Unable to create newOrderDishes vertical box: ", err)
	}

	newOrderCategoriesFrame, err := gtk.FrameNew("Категории меню")
	if err != nil {
		log.Fatal("Error on creating newOrderCategories frame")
	}

	newOrderDishesFrame, err := gtk.FrameNew("Блюда")
	if err != nil {
		log.Fatal("Error on creating newOrderDishes frame")
	}

	createNewOrderCategoriesTreeView()
	newOrderCategoriesScrolledWindow, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		log.Fatalf("Error on creating newOrderCategories scrolled window")
	}
	newOrderCategoriesScrolledWindow.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	newOrderCategoriesScrolledWindow.Add(newOrderCategoriesTreeView)

	newOrderCategoriesSelection, err := newOrderCategoriesTreeView.GetSelection()
	if err != nil {
		log.Fatal("Error on getting newOrderCategories selection")
	}
	newOrderCategoriesSelection.Connect("changed", newOrderCategoriesSelectionChanged, nil)

	newOrderCategoriesVbox.Add(newOrderCategoriesFrame)
	newOrderCategoriesVbox.PackStart(newOrderCategoriesScrolledWindow, true, true, 3)

	createNewOrderDishesTreeView()
	newOrderDishesScrolledWindow, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		log.Fatalf("Error on creating newOrderDishes scrolled window")
	}
	newOrderDishesScrolledWindow.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	newOrderDishesScrolledWindow.Add(newOrderDishesTreeView)

	newOrderDishesVbox.Add(newOrderDishesFrame)
	newOrderDishesVbox.PackStart(newOrderDishesScrolledWindow, true, true, 3)

	menuHbox.PackStart(newOrderCategoriesVbox, false, false, 3)
	menuHbox.PackStart(newOrderDishesVbox, true, true, 3)

	newOrderDishAddButton, err := gtk.ButtonNewWithLabel("Добавить к заказу")
	if err != nil {
		log.Fatal("Unable to create add button: ", err)
	}
	newOrderDishAddButton.Connect("clicked", newOrderDishAddButtonClicked, nil)

	newOrderVbox.PackStart(menuHbox, true, true, 3)
	newOrderVbox.PackStart(newOrderDishAddButton, false, true, 3)

	newOrderWindow.Add(newOrderVbox)
	return newOrderWindow
}
