package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var (
	lock sync.RWMutex
)

func main() {
	db := database{"shoes": 50, "socks": 5}
	mux := http.NewServeMux()
	mux.HandleFunc("/list", db.list)
	mux.HandleFunc("/price", db.price)
	mux.HandleFunc("/create", db.create)
	mux.HandleFunc("/delete", db.delete)
	mux.HandleFunc("/update", db.update)

	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type database map[string]dollars

func (db database) list(w http.ResponseWriter, req *http.Request) {
	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func (db database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	lock.RLock()
	if price, ok := db[item]; ok {
		fmt.Fprintf(w, "%s\n", price)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
	lock.RUnlock()
}

func (db database) create(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	newPrice := req.URL.Query().Get("price")
	value, _ := strconv.ParseFloat(newPrice, 32)
	floatPrice := dollars(float32(value))
	lock.Lock()
	if _, ok := db[item]; ok {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "already such item: %q\n", item)
	} else {
		db[item] = floatPrice
		fmt.Fprintf(w, "insert item %q price %s\n", item, floatPrice)
	}
	lock.Unlock()
}

func (db database) delete(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	lock.Lock()
	if _, ok := db[item]; ok {
		delete(db, item)
		fmt.Fprintf(w, "deltete item %q\n", item)

	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
	lock.Unlock()
}

func (db database) update(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	newPrice := req.URL.Query().Get("price")
	value, _ := strconv.ParseFloat(newPrice, 32)
	floatPrice := dollars(float32(value))
	lock.Lock()
	if _, ok := db[item]; ok {
		db[item] = floatPrice
		fmt.Fprintf(w, "update item %q price %s\n", item, floatPrice)

	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "already such item: %q\n", item)
	}
	lock.Unlock()
}
