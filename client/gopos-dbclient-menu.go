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
)

var categoriesTreeView *gtk.TreeView
var categoriesListStore *gtk.ListStore

const (
	COLUMN_CATEGORIES_ID = iota
	COLUMN_CATEGORIES_NAME
)

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
		fmt.Printf("%v %v\n", responseMap, err)
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

			id, err := strconv.Atoi(responseMap["id"])
			if err != nil {
				log.Fatal("Error on converting id to int: ", err)
			}
			if responseMap["result"] == "ERR" {
				messageDialog := gtk.MessageDialogNew(mainWindow,
					gtk.DIALOG_MODAL, gtk.MESSAGE_WARNING, gtk.BUTTONS_OK,
					responseMap["error"])
				messageDialog.Run()
				messageDialog.Destroy()
				return
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
		log.Fatal("Unable to get worker name entry text: ", err)
	}
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
	scrolledWindow, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		log.Fatalf("Error on creating workers scrolled window")
	}
	scrolledWindow.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	scrolledWindow.Add(categoriesTreeView)

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

	categoryDeleteSelectedButton, err := gtk.ButtonNewWithLabel("Удалить выбранного")
	if err != nil {
		log.Fatal("Unable to create add button: ", err)
	}
//	cateogoryDeleteSelectedButton.Connect("clicked", tableDeleteSelectedButtonClicked, nil)

	categoriesFormHbox.PackStart(categoryNameLabel, false, false, 3)
	categoriesFormHbox.PackStart(categoryNameEntry, true, true, 3)
	categoriesFormHbox.PackStart(categoryAddButton, true, true, 3)
	categoriesFormHbox.PackStart(categoryDeleteSelectedButton, true, true, 3)


	categoriesVbox.Add(categoriesFrame)
	categoriesVbox.PackStart(scrolledWindow, true, true, 3)
	categoriesVbox.PackStart(categoriesFormHbox, false, false, 3)
	dishesVbox.Add(dishesFrame)
	menuHbox.PackStart(categoriesVbox, false, false, 3)
	menuHbox.PackStart(dishesVbox, true, true, 3)
	return menuHbox
}
