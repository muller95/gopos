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
	COLUMN_EXISTING_ORDER_CATEGORIES_ID = iota
	COLUMN_EXISTING_ORDER_CATEGORIES_NAME
)

const (
	COLUMN_EXISTING_ORDER_DISHES_ID = iota
	COLUMN_EXISTING_ORDER_DISHES_NAME
	COLUMN_EXISTING_ORDER_DISHES_PRICE
)

var existingOrderWindow *gtk.Window

var existingOrderCategoriesTreeView *gtk.TreeView
var existingOrderCategoriesListStore *gtk.ListStore

var existingOrderDishesTreeView *gtk.TreeView
var existingOrderDishesListStore *gtk.ListStore

func createExistingOrderCategoriesTreeView() {
	var err error
	existingOrderCategoriesTreeView, err = gtk.TreeViewNew()

	if err != nil {
		log.Fatal("Unable to create cateogories tree view: ", err)
	}

	existingOrderCategoriesTreeView.AppendColumn(createColumn("ID",
		COLUMN_EXISTING_ORDER_CATEGORIES_ID))
	existingOrderCategoriesTreeView.AppendColumn(createColumn("Название",
		COLUMN_EXISTING_ORDER_CATEGORIES_NAME))

	existingOrderCategoriesListStore, err = gtk.ListStoreNew(glib.TYPE_INT, glib.TYPE_STRING)
	if err != nil {
		log.Fatal("Unable to create existingOrderCategories store: ", err)
	}

	existingOrderCategoriesTreeView.SetModel(existingOrderCategoriesListStore)
}

func existingOrderCategoryAddRow(id int, name string) {
	iter := existingOrderCategoriesListStore.Append()

	err := existingOrderCategoriesListStore.Set(iter, []int{COLUMN_EXISTING_ORDER_CATEGORIES_ID,
		COLUMN_NEW_ORDER_CATEGORIES_NAME}, []interface{}{id, name})

	if err != nil {
		log.Fatal("Unable to add newOrderCategories row: ", err)
	}
}

func createExistingOrderDishesTreeView() {
	var err error
	existingOrderDishesTreeView, err = gtk.TreeViewNew()

	if err != nil {
		log.Fatal("Unable to create existingOrderDishes tree view: ", err)
	}

	existingOrderDishesTreeView.AppendColumn(createColumn("ID", COLUMN_EXISTING_ORDER_DISHES_ID))
	existingOrderDishesTreeView.AppendColumn(createColumn("Название", COLUMN_EXISTING_ORDER_DISHES_NAME))
	existingOrderDishesTreeView.AppendColumn(createColumn("Цена", COLUMN_EXISTING_ORDER_DISHES_PRICE))

	existingOrderDishesListStore, err = gtk.ListStoreNew(glib.TYPE_INT, glib.TYPE_STRING,
		glib.TYPE_DOUBLE)
	if err != nil {
		log.Fatal("Unable to create existingOrderDishes store: ", err)
	}

	existingOrderDishesTreeView.SetModel(existingOrderDishesListStore)
}

func existingOrderDishAddRow(id int, name string, price float64) {
	iter := existingOrderDishesListStore.Append()

	err := existingOrderDishesListStore.Set(iter, []int{COLUMN_EXISTING_ORDER_DISHES_ID, COLUMN_EXISTING_ORDER_DISHES_NAME,
		COLUMN_EXISTING_ORDER_DISHES_PRICE}, []interface{}{id, name, price})

	if err != nil {
		log.Fatal("Unable to add existingOrderDish row: ", err)
	}
}

func getExistingOrderCategories() {
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
		existingOrderCategoryAddRow(id, responseMap["name"])
	}
	conn.Close()
}

func existingOrderCategoriesSelectionChanged(selection *gtk.TreeSelection) {
	existingOrderDishesListStore.Clear()
	selection, err := existingOrderCategoriesTreeView.GetSelection()
	if err != nil {
		log.Fatal("Error on getting existingOrderCategories selection")
	}
	rows := selection.GetSelectedRows(existingOrderCategoriesListStore)
	if rows == nil {
		return
	}

	path := rows.Data().(*gtk.TreePath)
	iter, err := existingOrderCategoriesListStore.GetIter(path)

	value, err := existingOrderCategoriesListStore.GetValue(iter, 0)
	if err != nil {
		log.Fatal("Error on getting value: ", err)
	}
	existingOrderCategoryId := value.GetInt()

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", goposServerIp,
		goposServerPort))

	if err != nil {
		log.Fatal("Unable to connect to server")
	}

	requestMap := make(map[string]string)
	requestMap["group"] = "DISH"
	requestMap["action"] = "GET"
	requestMap["password"] = goposServerPassword
	requestMap["category_id"] = fmt.Sprintf("%d", existingOrderCategoryId)
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
		existingOrderDishAddRow(id, responseMap["name"], price)
	}
}

/*func existingOrderDishAddButtonClicked(btn *gtk.Button) {
	errMessage := ""
	selection, err := existingOrderDishesTreeView.GetSelection()
	if err != nil {
		log.Fatal("Error on getting existingOrderCategories selection")
	}
	rows := selection.GetSelectedRows(existingOrderDishesListStore)
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
	iter, err := existingOrderDishesListStore.GetIter(path)
	if err != nil {
		log.Fatal("Error on getting iter: ", err)
	}

	value, err := existingOrderDishesListStore.GetValue(iter, 0)
	if err != nil {
		log.Fatal("Error on getting value: ", err)
	}
	existingOrderDishId := value.GetInt()

	value, err = existingOrderDishesListStore.GetValue(iter, 1)
	if err != nil {
		log.Fatal("Error on getting value: ", err)
	}
	existingOrderDishName, err := value.GetString()
	if err != nil {
		log.Fatal("Error on getting string: ", err)
	}

	value, err = existingOrderDishesListStore.GetValue(iter, 2)
	if err != nil {
		log.Fatal("Error on getting value: ", err)
	}
	existingOrderDishPrice := value.GetDouble()
	orderPrice += existingOrderDishPrice
	existingOrderPriceLabel.SetText(fmt.Sprintf("Цена: %.2f", orderPrice))
	existingOrderAddRow(existingOrderDishId, existingOrderDishName, existingOrderDishPrice)
}*/

func existingOrderCreateWindow() *gtk.Window {
	//creates menu tabpage
	var err error
	existingOrderWindow, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	existingOrderWindow.SetTitle("gopos-monitor-existingorder")
	existingOrderWindow.SetDefaultSize(800, 600)

	existingOrderVbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	if err != nil {
		log.Fatal("Unable to create main vertica; box: ", err)
	}

	menuHbox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	if err != nil {
		log.Fatal("Unable to create menu horizontal box: ", err)
	}

	existingOrderCategoriesVbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	if err != nil {
		log.Fatal("Unable to create existingOrderCategories vertical box: ", err)
	}

	existingOrderDishesVbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	if err != nil {
		log.Fatal("Unable to create existingOrderDishes vertical box: ", err)
	}

	existingOrderCategoriesFrame, err := gtk.FrameNew("Категории меню")
	if err != nil {
		log.Fatal("Error on creating existingOrderCategories frame")
	}

	existingOrderDishesFrame, err := gtk.FrameNew("Блюда")
	if err != nil {
		log.Fatal("Error on creating existingOrderDishes frame")
	}

	createExistingOrderCategoriesTreeView()
	existingOrderCategoriesScrolledWindow, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		log.Fatalf("Error on creating existingOrderCategories scrolled window")
	}
	existingOrderCategoriesScrolledWindow.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	existingOrderCategoriesScrolledWindow.Add(existingOrderCategoriesTreeView)

	existingOrderCategoriesSelection, err := existingOrderCategoriesTreeView.GetSelection()
	if err != nil {
		log.Fatal("Error on getting existingOrderCategories selection")
	}
	existingOrderCategoriesSelection.Connect("changed", existingOrderCategoriesSelectionChanged, nil)

	existingOrderCategoriesVbox.Add(existingOrderCategoriesFrame)
	existingOrderCategoriesVbox.PackStart(existingOrderCategoriesScrolledWindow, true, true, 3)

	createExistingOrderDishesTreeView()
	existingOrderDishesScrolledWindow, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		log.Fatalf("Error on creating existingOrderDishes scrolled window")
	}
	existingOrderDishesScrolledWindow.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	existingOrderDishesScrolledWindow.Add(existingOrderDishesTreeView)

	existingOrderDishesVbox.Add(existingOrderDishesFrame)
	existingOrderDishesVbox.PackStart(existingOrderDishesScrolledWindow, true, true, 3)

	menuHbox.PackStart(existingOrderCategoriesVbox, false, false, 3)
	menuHbox.PackStart(existingOrderDishesVbox, true, true, 3)

	existingOrderDishAddButton, err := gtk.ButtonNewWithLabel("Добавить к заказу")
	if err != nil {
		log.Fatal("Unable to create add button: ", err)
	}
	// existingOrderDishAddButton.Connect("clicked", existingOrderDishAddButtonClicked, nil)

	existingOrderVbox.PackStart(menuHbox, true, true, 3)
	existingOrderVbox.PackStart(existingOrderDishAddButton, false, true, 3)

	existingOrderWindow.Add(existingOrderVbox)
	return existingOrderWindow
}
