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

type DishInfo struct {
	DishNameEntry  *gtk.Entry
	DishPriceEntry *gtk.Entry
}

const (
	COLUMN_CATEGORIES_ID = iota
	COLUMN_CATEGORIES_NAME
)

const (
	COLUMN_DISHES_ID = iota
	COLUMN_DISHES_NAME
	COLUMN_DISHES_PRICE
)

var newOrderWindow *gtk.Window

var categoriesTreeView *gtk.TreeView
var categoriesListStore *gtk.ListStore

var dishesTreeView *gtk.TreeView
var dishesListStore *gtk.ListStore

func createCategoriesTreeView() {
	var err error
	categoriesTreeView, err = gtk.TreeViewNew()

	if err != nil {
		log.Fatal("Unable to create cateogories tree view: ", err)
	}

	categoriesTreeView.AppendColumn(createColumn("ID", COLUMN_CATEGORIES_ID))
	categoriesTreeView.AppendColumn(createColumn("Название", COLUMN_CATEGORIES_NAME))

	categoriesListStore, err = gtk.ListStoreNew(glib.TYPE_INT, glib.TYPE_STRING)
	if err != nil {
		log.Fatal("Unable to create categories store: ", err)
	}

	categoriesTreeView.SetModel(categoriesListStore)
}

func categoryAddRow(id int, name string) {
	iter := categoriesListStore.Append()

	err := categoriesListStore.Set(iter, []int{COLUMN_CATEGORIES_ID, COLUMN_CATEGORIES_NAME},
		[]interface{}{id, name})

	if err != nil {
		log.Fatal("Unable to add categories row: ", err)
	}
}

func createDishesTreeView() {
	var err error
	dishesTreeView, err = gtk.TreeViewNew()

	if err != nil {
		log.Fatal("Unable to create dishes tree view: ", err)
	}

	dishesTreeView.AppendColumn(createColumn("ID", COLUMN_DISHES_ID))
	dishesTreeView.AppendColumn(createColumn("Название", COLUMN_DISHES_NAME))
	dishesTreeView.AppendColumn(createColumn("Цена", COLUMN_DISHES_PRICE))

	dishesListStore, err = gtk.ListStoreNew(glib.TYPE_INT, glib.TYPE_STRING,
		glib.TYPE_DOUBLE)
	if err != nil {
		log.Fatal("Unable to create dishes store: ", err)
	}

	dishesTreeView.SetModel(dishesListStore)
}

func dishAddRow(id int, name string, price float64) {
	iter := dishesListStore.Append()

	err := dishesListStore.Set(iter, []int{COLUMN_DISHES_ID, COLUMN_DISHES_NAME,
		COLUMN_DISHES_PRICE}, []interface{}{id, name, price})

	if err != nil {
		log.Fatal("Unable to add dish row: ", err)
	}
}

func getCategories() {
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
		categoryAddRow(id, responseMap["name"])
	}
	conn.Close()
}

func categoriesSelectionChanged(selection *gtk.TreeSelection) {
	dishesListStore.Clear()
	selection, err := categoriesTreeView.GetSelection()
	if err != nil {
		log.Fatal("Error on getting categories selection")
	}
	rows := selection.GetSelectedRows(categoriesListStore)
	if rows == nil {
		return
	}

	path := rows.Data().(*gtk.TreePath)
	iter, err := categoriesListStore.GetIter(path)

	value, err := categoriesListStore.GetValue(iter, 0)
	if err != nil {
		log.Fatal("Error on getting value: ", err)
	}
	categoryId := value.GetInt()

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", goposServerIp,
		goposServerPort))

	if err != nil {
		log.Fatal("Unable to connect to server")
	}

	requestMap := make(map[string]string)
	requestMap["group"] = "DISH"
	requestMap["action"] = "GET"
	requestMap["password"] = goposServerPassword
	requestMap["category_id"] = fmt.Sprintf("%d", categoryId)
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
		dishAddRow(id, responseMap["name"], price)
	}
}

func dishAddButtonClicked(btn *gtk.Button) {
	errMessage := ""
	selection, err := dishesTreeView.GetSelection()
	if err != nil {
		log.Fatal("Error on getting categories selection")
	}
	rows := selection.GetSelectedRows(dishesListStore)
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
	iter, err := dishesListStore.GetIter(path)
	if err != nil {
		log.Fatal("Error on getting iter: ", err)
	}

	value, err := dishesListStore.GetValue(iter, 0)
	if err != nil {
		log.Fatal("Error on getting value: ", err)
	}
	dishId := value.GetInt()

	value, err = dishesListStore.GetValue(iter, 1)
	if err != nil {
		log.Fatal("Error on getting value: ", err)
	}
	dishName, err := value.GetString()
	if err != nil {
		log.Fatal("Error on getting string: ", err)
	}

	value, err = dishesListStore.GetValue(iter, 2)
	if err != nil {
		log.Fatal("Error on getting value: ", err)
	}
	dishPrice := value.GetDouble()
	orderPrice += dishPrice
	newOrderPriceLabel.SetText(fmt.Sprintf("Цена: %.2f", orderPrice))
	newOrderAddRow(dishId, dishName, dishPrice)
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

	categoriesVbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	if err != nil {
		log.Fatal("Unable to create categories vertical box: ", err)
	}

	dishesVbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	if err != nil {
		log.Fatal("Unable to create dishes vertical box: ", err)
	}

	categoriesFrame, err := gtk.FrameNew("Категории меню")
	if err != nil {
		log.Fatal("Error on creating categories frame")
	}

	dishesFrame, err := gtk.FrameNew("Блюда")
	if err != nil {
		log.Fatal("Error on creating dishes frame")
	}

	createCategoriesTreeView()
	categoriesScrolledWindow, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		log.Fatalf("Error on creating categories scrolled window")
	}
	categoriesScrolledWindow.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	categoriesScrolledWindow.Add(categoriesTreeView)

	categoriesSelection, err := categoriesTreeView.GetSelection()
	if err != nil {
		log.Fatal("Error on getting categories selection")
	}
	categoriesSelection.Connect("changed", categoriesSelectionChanged, nil)

	categoriesVbox.Add(categoriesFrame)
	categoriesVbox.PackStart(categoriesScrolledWindow, true, true, 3)

	createDishesTreeView()
	dishesScrolledWindow, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		log.Fatalf("Error on creating dishes scrolled window")
	}
	dishesScrolledWindow.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	dishesScrolledWindow.Add(dishesTreeView)

	dishesVbox.Add(dishesFrame)
	dishesVbox.PackStart(dishesScrolledWindow, true, true, 3)

	menuHbox.PackStart(categoriesVbox, false, false, 3)
	menuHbox.PackStart(dishesVbox, true, true, 3)

	dishAddButton, err := gtk.ButtonNewWithLabel("Добавить к заказу")
	if err != nil {
		log.Fatal("Unable to create add button: ", err)
	}
	dishAddButton.Connect("clicked", dishAddButtonClicked, nil)

	newOrderVbox.PackStart(menuHbox, true, true, 3)
	newOrderVbox.PackStart(dishAddButton, false, true, 3)

	newOrderWindow.Add(newOrderVbox)
	return newOrderWindow
}
