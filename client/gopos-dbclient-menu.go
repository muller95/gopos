//main package for gopos database client
//this file contains functions for menu tree view
package main

import (
//	"encoding/json"
//	"fmt"
//	"io"
	"log"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
//	"net"
//	"strconv"
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

	categoriesListStore, err = gtk.ListStoreNew(glib.TYPE_INT)
	if err != nil {
		log.Fatal("Unable to create tables categories store: ", err)
	}

	categoriesTreeView.SetModel(categoriesListStore)
}

/*func categoryAddRow(number int) {
	iter := tablesListStore.Append()

	err := tablesListStore.Set(iter, []int{ COLUMN_TABLES_NUMBER },
		[]interface{} { number })

	if err != nil {
		log.Fatal("Unable to add tables row: ", err)
	}
}*/

func menuCreatePage() *gtk.Box {
	//creates tables tabpage
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

/*	tablesFormHbox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	if err != nil {
		log.Fatal("Unable to create tables form horizontal box: ", err)
	}

	tableNumberLabel, err := gtk.LabelNew("Номер столика:")
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}

	tableNumberEntry, err := gtk.EntryNew()
	if err != nil {
		log.Fatal("Unable to create entry: ", err)
	}

	tableAddButton, err := gtk.ButtonNewWithLabel("Добавить")
	if err != nil {
		log.Fatal("Unable to create add button: ", err)
	}
	tableAddButton.Connect("clicked", tableAddButtonClicked, tableNumberEntry)

	tableDeleteSelectedButton, err := gtk.ButtonNewWithLabel("Удалить выбранного")
	if err != nil {
		log.Fatal("Unable to create add button: ", err)
	}
	tableDeleteSelectedButton.Connect("clicked", tableDeleteSelectedButtonClicked, nil)

	tablesFormHbox.PackStart(tableNumberLabel, false, false, 3)
	tablesFormHbox.PackStart(tableNumberEntry, true, true, 3)
	tablesFormHbox.PackStart(tableAddButton, true, true, 3)
	tablesFormHbox.PackStart(tableDeleteSelectedButton, true, true, 3)*/


//	tablesVbox.PackStart(tablesFormHbox, false, false, 3)
	categoriesVbox.Add(categoriesFrame)
	categoriesVbox.PackStart(scrolledWindow, true, true, 3)
	dishesVbox.Add(dishesFrame)
	menuHbox.PackStart(categoriesVbox, false, false, 3)
	menuHbox.PackStart(dishesVbox, true, true, 3)
	return menuHbox
}
