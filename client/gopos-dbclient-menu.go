//main package for gopos database client
//this file contains functions for menu tree view
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
	"strings"
)

type DishInfo struct{
	DishNameEntry *gtk.Entry
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

	err := categoriesListStore.Set(iter, []int{ COLUMN_CATEGORIES_ID, COLUMN_CATEGORIES_NAME },
		[]interface{} { id, name })

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

	err := dishesListStore.Set(iter, []int{ COLUMN_DISHES_ID, COLUMN_DISHES_NAME,
		COLUMN_DISHES_PRICE }, []interface{} { id, name, price })

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

func categoryAddButtonClicked(btn *gtk.Button, categoryNameEntry *gtk.Entry) {
	categoryName, err := categoryNameEntry.GetText()

	if err == nil {
		categoryName = strings.Trim(categoryName, " ")
		if len(categoryName) > 0 {
			conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", goposServerIp,
				goposServerPort))

			if err != nil {
				log.Fatal("Unable to connect to server")
			}

			requestMap := make(map[string]string)
			requestMap["group"] = "CATEGORY"
			requestMap["action"] = "ADD"
			requestMap["password"] = goposServerPassword
			requestMap["name"] = categoryName
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

			if responseMap["result"] == "ERR" {
				messageDialog := gtk.MessageDialogNew(mainWindow,
					gtk.DIALOG_MODAL, gtk.MESSAGE_WARNING, gtk.BUTTONS_OK,
					responseMap["error"])
				messageDialog.Run()
				messageDialog.Destroy()
				return
			}

			id, err := strconv.Atoi(responseMap["id"])
			if err != nil {
				log.Fatal("Error on converting id to int: ", err)
			}

			categoryAddRow(id, categoryName)

			conn.Close()
		} else {
			messageDialog := gtk.MessageDialogNew(mainWindow, gtk.DIALOG_MODAL,
				gtk.MESSAGE_WARNING, gtk.BUTTONS_OK, "Введите название категории")
			messageDialog.Run()
			messageDialog.Destroy()
		}
	} else {
		log.Fatal("Unable to get category name entry text: ", err)
	}
}

func categoryDeleteSelectedButtonClicked() {
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
	if err != nil {
		log.Fatal("Error on getting iter: ", err)
	}
	value, err := categoriesListStore.GetValue(iter, 0)
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
	requestMap["group"] = "CATEGORY"
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
		categoriesListStore.Remove(iter)
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

func dishAddButtonClicked(btn *gtk.Button, dishInfo *DishInfo) {
	dishName, err := dishInfo.DishNameEntry.GetText()
	if err != nil {
		log.Fatal("Unable to get dish name entry text: ", err)
	}
	dishName = strings.Trim(dishName, " ")

	dishPrice, err := dishInfo.DishPriceEntry.GetText()
	if err != nil {
		log.Fatal("Unable to get dish price entry text: ", err)
	}
	dishPrice = strings.Trim(dishPrice, " ")

	errMessage := ""
	if len(dishName) == 0 {
		errMessage += "Введите название блюда. "
	}

	price := 0.0
	if len(dishPrice) == 0 {
		errMessage += "Введите цену блюда."
	} else {
		price, err = strconv.ParseFloat(dishPrice, 64)
		if err != nil  || price < 0.0{
			errMessage += "Цена должна быть положительным вещественным числом, " +
			"например, 123.45. "
		}
	}

	selection, err := categoriesTreeView.GetSelection()
	if err != nil {
		log.Fatal("Error on getting categories selection")
	}
	rows := selection.GetSelectedRows(dishesListStore)
	if rows == nil {
		errMessage += "Выберите категорию для добавления блюда. "
	}

	if errMessage != "" {
		messageDialog := gtk.MessageDialogNew(mainWindow, gtk.DIALOG_MODAL,
			gtk.MESSAGE_WARNING, gtk.BUTTONS_OK, errMessage)
		messageDialog.Run()
		messageDialog.Destroy()
		return
	}

	path := rows.Data().(*gtk.TreePath)
	iter, err := categoriesListStore.GetIter(path)
	if err != nil {
		log.Fatal("Error on getting iter: ", err)
	}
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
	requestMap["action"] = "ADD"
	requestMap["password"] = goposServerPassword
	requestMap["name"] = dishName
	requestMap["price"] = fmt.Sprintf("%f", price)
	requestMap["category_id"] = fmt.Sprintf("%d", categoryId)
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

	if responseMap["result"] == "ERR" {
		messageDialog := gtk.MessageDialogNew(mainWindow,
			gtk.DIALOG_MODAL, gtk.MESSAGE_WARNING, gtk.BUTTONS_OK,
			responseMap["error"])
			messageDialog.Run()
		messageDialog.Destroy()
		return
	}

	id, err := strconv.Atoi(responseMap["id"])
	if err != nil {
		log.Fatal("Error on converting id to int: ", err)
	}
	dishAddRow(id, dishName, price)
	conn.Close()
}

func dishDeleteSelectedButtonClicked() {
	selection, err := dishesTreeView.GetSelection()
	if err != nil {
		log.Fatal("Error on getting categories selection")
	}
	rows := selection.GetSelectedRows(dishesListStore)
	if rows == nil {
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
	id := value.GetInt()

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", goposServerIp,
		goposServerPort))

	if err != nil {
		log.Fatal("Unable to connect to server")
	}

	requestMap := make(map[string]string)
	requestMap["group"] = "DISH"
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
		dishesListStore.Remove(iter)
	}
	conn.Close()
}

func menuCreatePage() *gtk.Box {
	//creates menu tabpage
	menuHbox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	if err != nil {
		log.Fatal("Unable to create main horizontal box: ", err)
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

	categoriesFormHbox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	if err != nil {
		log.Fatal("Unable to create categories form horizontal box: ", err)
	}

	categoryNameLabel, err := gtk.LabelNew("Название категории:")
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}

	categoryNameEntry, err := gtk.EntryNew()
	if err != nil {
		log.Fatal("Unable to create entry: ", err)
	}

	categoryAddButton, err := gtk.ButtonNewWithLabel("Добавить")
	if err != nil {
		log.Fatal("Unable to create add button: ", err)
	}
	categoryAddButton.Connect("clicked", categoryAddButtonClicked, categoryNameEntry)

	categoryDeleteSelectedButton, err := gtk.ButtonNewWithLabel("Удалить категорию")
	if err != nil {
		log.Fatal("Unable to create add button: ", err)
	}
	categoryDeleteSelectedButton.Connect("clicked", categoryDeleteSelectedButtonClicked, nil)

	categoriesFormHbox.PackStart(categoryNameLabel, false, false, 3)
	categoriesFormHbox.PackStart(categoryNameEntry, true, true, 3)
	categoriesFormHbox.PackStart(categoryAddButton, true, true, 3)
	categoriesFormHbox.PackStart(categoryDeleteSelectedButton, true, true, 3)


	categoriesVbox.Add(categoriesFrame)
	categoriesVbox.PackStart(categoriesScrolledWindow, true, true, 3)
	categoriesVbox.PackStart(categoriesFormHbox, false, false, 3)


	createDishesTreeView()
	dishesScrolledWindow, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		log.Fatalf("Error on creating dishes scrolled window")
	}
	dishesScrolledWindow.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	dishesScrolledWindow.Add(dishesTreeView)

	dishesFormHbox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	if err != nil {
		log.Fatal("Unable to create dishes form horizontal box: ", err)
	}

	dishNameLabel, err := gtk.LabelNew("Название блюда:")
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}

	dishNameEntry, err := gtk.EntryNew()
	if err != nil {
		log.Fatal("Unable to create entry: ", err)
	}
	dishPriceEntry, err := gtk.EntryNew()
	if err != nil {
		log.Fatal("Unable to create entry: ", err)
	}

	dishAddButton, err := gtk.ButtonNewWithLabel("Добавить")
	if err != nil {
		log.Fatal("Unable to create add button: ", err)
	}
	dishInfo := new(DishInfo)
	dishInfo.DishNameEntry = dishNameEntry
	dishInfo.DishPriceEntry = dishPriceEntry
	dishAddButton.Connect("clicked", dishAddButtonClicked, dishInfo)

	dishDeleteSelectedButton, err := gtk.ButtonNewWithLabel("Удалить блюдо")
	if err != nil {
		log.Fatal("Unable to create add button: ", err)
	}
	dishDeleteSelectedButton.Connect("clicked", dishDeleteSelectedButtonClicked, nil)

	dishesFormHbox.PackStart(dishNameLabel, false, false, 3)
	dishesFormHbox.PackStart(dishNameEntry, true, true, 3)
	dishesFormHbox.PackStart(dishPriceEntry, true, true, 3)
	dishesFormHbox.PackStart(dishAddButton, true, true, 3)
	dishesFormHbox.PackStart(dishDeleteSelectedButton, true, true, 3)

	dishesVbox.Add(dishesFrame)
	dishesVbox.PackStart(dishesScrolledWindow, true, true, 3)
	dishesVbox.PackStart(dishesFormHbox, false, false, 3)

	menuHbox.PackStart(categoriesVbox, false, false, 3)
	menuHbox.PackStart(dishesVbox, true, true, 3)
	return menuHbox
}
